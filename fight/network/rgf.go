/*
场景框架，驱动战斗场景逻辑，每个场景为一个协程
*/
package network

import (
	"com/log"
	"com/util"
	"fight/configs"
	"fight/types"
	"github.com/golang/protobuf/proto"
	"pb"
	"time"
)

type RoleMsg struct {
	RoleID uint32
	Msg    *pb.SrvMsg
}

type Evt struct {
	roleId uint32
	ses    *types.Session
}

type RgnForm struct {
	Id         uint32
	createTime int64
	nowTime    int64
	Roles      map[uint32]*types.Role
	MsgQ       chan RoleMsg
	StartData  *pb.MsgBattleStartData
	Battle     IRegion
	evts       chan Evt
	close      bool
}

func NewRgnForm(rgId uint32, msg *pb.MsgBattleStartData) *RgnForm {
	gk := configs.GetRegionLogic(msg.BattleID)
	//todo: you must init pos by use region logic
	if msg.BattleDifficulty == 0 {
		msg.BattleDifficulty = 1
	}
	for _, v := range msg.Roles {
		for _, v2 := range v.Fighters {
			if v2.Visible {
				v2.PosX = gk.BornBos.ServerCastPos[0]
				v2.PosY = gk.BornBos.ServerCastPos[1]
				v2.PosZ = gk.BornBos.ServerCastPos[2]
				v2.Yaw = gk.BornBos.ServerCastYaw
			}
		}
	}
	r := &RgnForm{
		Roles:     make(map[uint32]*types.Role, 10),
		MsgQ:      make(chan RoleMsg, 128),
		evts:      make(chan Evt, 128),
		Id:        rgId,
		StartData: msg,
	}

	for _, v := range msg.Roles {
		var vft *types.Fighter
		fts := make([]*types.Fighter, 0)
		for _, v2 := range v.Fighters {
			oneFt := CreateFighter(v2, false)
			if v2.Visible {
				vft = oneFt
			}
			fts = append(fts, oneFt)
		}
		r.Roles[v.Guid] = types.NewRole(v, fts, vft)
	}

	go r.run()

	return r
}

//Close	关闭，会清理掉
func (m *RgnForm) Close() {
	m.close = true
}

func (m *RgnForm) addMsg(roleID uint32, msg *pb.SrvMsg) {
	m.MsgQ <- RoleMsg{
		RoleID: roleID,
		Msg:    msg,
	}
}

func (m *RgnForm) bindRoleSession(roleID uint32, session *types.Session) {
	m.evts <- Evt{ses: session, roleId: roleID}
}

func (m *RgnForm) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			util.PrintStack()
		}
	}()

	m.createTime = time.Now().UnixNano() / (int64)(time.Millisecond)
	m.Battle.OnCreate(m)

	tRun := time.NewTicker(time.Millisecond * 50) //帧率暂定20
	tSec := time.NewTicker(time.Second)

	defer func() {
		tRun.Stop()
		tSec.Stop()
	}()
	for {
		select {
		case msg := <-m.MsgQ:
			role := m.Roles[msg.RoleID]
			if role == nil {
				//todo
				continue
			}
			err := HandleRoleMsg(msg.Msg.ID, msg.Msg.Msg, role, m)
			if err != nil {
				//todo
			}
		case e := <-m.evts:
			m.onEvt(e)
		case <-tRun.C:
			m.nowTime = time.Now().UnixNano() / (int64)(time.Millisecond)
			m.Battle.Run()
		case <-tSec.C:
			m.nowTime = time.Now().UnixNano() / (int64)(time.Millisecond)
			m.Battle.SecLoop()
		}

		if m.close {
			m.onClose()
			return
		}
	}
}

func (m *RgnForm) onClose() {
	MgrForRegion.Del(m.Id)
	for _, r := range m.Roles {
		if r != nil && r.Ses != nil {
			r.Ses.Close()
		}
	}
}

func (m *RgnForm) onEvt(e Evt) {
	//现在只有绑定
	role := m.Roles[e.roleId]
	if role == nil {
		e.ses.Close() //没找到，这属于异常情况（有场景玩家应该就存在），直接关闭连接
		return
	}
	if role.Ses != nil {
		role.Ses.Close() //挤掉其他连接
	}
	role.Ses = e.ses
	if e.ses == nil {
		m.Battle.OnRoleDisConnect(role)
		return //掉线了
	}
	m.Battle.OnRoleConnect(role)
}

func (m *RgnForm) SendTo(roleId uint32, msgID pb.MsgIDS2C, msgData proto.Message) {
	role := m.Roles[roleId]
	if role != nil {
		role.Send(msgID, msgData)
	}
}

func (m *RgnForm) SendToAll(msgID pb.MsgIDS2C, msgData proto.Message) {
	for _, v := range m.Roles {
		if v.HadEnter {
			v.Send(msgID, msgData)
		}
	}
}

func (m *RgnForm) SendToAllAndNoEntered(msgID pb.MsgIDS2C, msgData proto.Message) {
	for _, v := range m.Roles {
		v.Send(msgID, msgData)
	}
}
