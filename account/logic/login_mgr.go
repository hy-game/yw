package logic

import (
	"com/log"
	"math"
	"pb"
	"time"
)

/*
todo 玩家所在game这些信息，重连后需要重建
*/
type LoginInfo struct {
	lastLoginTime int64
	lastAckTime   int64
	gameId        uint16
	loginCnt      uint8 //只用于 清除gameId
}

type GameInfo struct {
	Info       *pb.MsgGameHeartBeat
	LastHBTime int64
	closed     bool
}

func (this *GameInfo) isAlive() bool {
	return time.Now().Unix()-this.LastHBTime < 20
}

var loginMgr = newLoginMgr()

type LoginMgr struct {
	ack chan *LoginAck
	evt chan evtParam

	datas  map[string]*LoginInfo
	gmInfo map[uint16]GameInfo
}

func newLoginMgr() *LoginMgr {
	m := &LoginMgr{
		ack:    make(chan *LoginAck, 1000),
		evt:    make(chan evtParam, 1000),
		datas:  make(map[string]*LoginInfo),
		gmInfo: make(map[uint16]GameInfo),
	}

	go m.run()

	return m
}

type Op int

const (
	OpGameInfo Op = iota
	OpLogin
	OpRoleClear
	OpGameRoles
)

type evtParam struct {
	op       Op
	login    *pb.MsgLogin
	gmSrv    *pb.MsgGameHeartBeat
	roleInfo *pb.MsgROnOffLine
	serId    uint16
	gmRoles *pb.MsgGameRoles
}

type LoginAck struct {
	loginReq *pb.MsgLogin
	success  bool
}

func (m *LoginMgr) Ack(data *LoginAck) {
	m.ack <- data
}

func (m *LoginMgr) PostEvt(e evtParam) {
	m.evt <- e
}

func (m *LoginMgr) Login(data *pb.MsgLogin, gateId uint16) {
	data.GtID = uint32(gateId)
	m.PostEvt(evtParam{
		op:    OpLogin,
		login: data,
		serId: gateId,
	})
}

func (m *LoginMgr) run() {
	log.Debug("start login mgr run")
	defer func() {
		log.Debug("stop login mgr run")
	}()
	for {
		select {
		case ack := <-m.ack:
			m.onLoginCheckAck(ack.loginReq, ack.success)
		case e := <-m.evt:
			m.onEvent(e)
		}
	}
}

func (m *LoginMgr) loginCheck(req *pb.MsgLogin) {
	now := time.Now().Unix()
	info := m.datas[req.Acc]
	if info == nil {
		m.datas[req.Acc] = &LoginInfo{}
		info = m.datas[req.Acc]
	}
	if now < info.lastLoginTime+5 { //CD ing
		return
	}

	info.lastLoginTime = now

	//todo 投递sdk验证消息，开协程或者做个task，然后把返回消息投递回来
	//这里暂时先直接成功
	m.Ack(&LoginAck{req, true})
}

func (m *LoginMgr) onLoginCheckAck(loginReq *pb.MsgLogin, success bool) {
	info := m.datas[loginReq.Acc]
	if info == nil {
		return
	}
	msg := &pb.MsgLoginAck{Ret: pb.LoginCode_LCSuccess, Data: loginReq}

	if success {
		now := time.Now().Unix()
		if now < info.lastAckTime+5 {
			return //CD ing 过滤抖动
		}
		info.lastAckTime = now

		//获取 game id
		if info.gameId != 0 { //已登录过
			serInfo, ok := m.getGameInfo(info.gameId)
			if !ok { //宕机或者掉线时，需要等重启或者重连
				msg.Ret = pb.LoginCode_LCNoGame
			} else if serInfo.closed {
				gameId := m.getMinRoleGameId()
				if gameId == 0 {
					msg.Ret = pb.LoginCode_LCNoGame
				} else {
					info.gameId = gameId
					info.loginCnt++
					msg.GameID = uint32(gameId)
				}
			} else { //找到
				msg.GameID = uint32(info.gameId)
			}
		} else { //第一次登录
			gameId := m.getMinRoleGameId()
			if gameId == 0 {
				msg.Ret = pb.LoginCode_LCNoGame
			} else {
				info.gameId = gameId
				info.loginCnt++
				msg.GameID = uint32(gameId)
			}
		}
	}
	SendToGt(uint16(loginReq.GtID), pb.MsgIDS2S_Acc2GtLoginAck, msg)
}

func (m *LoginMgr) onEvent(e evtParam) {
	switch e.op {
	case OpLogin:
		m.loginCheck(e.login)
	case OpGameInfo:
		m.updateGameInfo(e.gmSrv, e.serId)
	case OpGameRoles:
		m.onGameRoles(e.gmRoles, e.serId)
	case OpRoleClear:
		m.onRoleClear(e.roleInfo.Acc)
	}
}

func (m *LoginMgr) getGameInfo(serId uint16) (GameInfo, bool) {
	if data, ok := m.gmInfo[serId]; ok {
		if !data.isAlive() {
			return data, false
		} else {
			return data, true
		}
	} else {
		return GameInfo{
			Info:       nil,
			LastHBTime: 0,
			closed:     false,
		}, false
	}
}

func (m *LoginMgr) getMinRoleGameId() (serId uint16) {
	min := int32(math.MaxInt32)
	for k, v := range m.gmInfo {
		if !v.isAlive() {
			continue
		}
		if v.Info.RoleCnt < min {
			min = v.Info.RoleCnt
			serId = k
		}
	}
	return
}

func (m *LoginMgr) updateGameInfo(info *pb.MsgGameHeartBeat, serID uint16) {
	m.gmInfo[serID] = GameInfo{
		Info:       info,
		LastHBTime: time.Now().Unix(),
	}
}

func (m *LoginMgr)onRoleClear(acc string){
	if data, ok := m.datas[acc]; ok {
		data.loginCnt--
		if data.loginCnt == 0 {
			data.gameId = 0
			log.Tracef("reset %s gameId=0", acc)
		}
	} else {
		log.Warnf("recv role online but not find login info")
	}
}

func (m *LoginMgr)onGameRoles(msg *pb.MsgGameRoles, serID uint16){
	for _, v := range msg.Roles{
		info := &LoginInfo{
			lastLoginTime: 0,
			lastAckTime:   0,
			gameId:        serID,
			loginCnt:      1,
		}
		m.datas[v] = info
	}
	log.Infof("rebuild role in game:%d", serID)
	WaitGameInfo.Done()
}