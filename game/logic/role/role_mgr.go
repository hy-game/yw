package role

import (
	"com/log"
	"game/network"
	"game/types"
	"pb"
	"time"
)

type RoleCacheState int

const (
	RoleInit RoleCacheState = iota
	RoleLoading
	RoleOnline
	RoleOffline
)

type RoleCache struct {
	Role      *types.Role //不能访问role里面的内容
	State     RoleCacheState
	StateTime int64
	SesID     uint32
}

func (this *RoleCache) setState(state RoleCacheState) {
	this.State = state
	this.StateTime = time.Now().Unix()
}

type cacheOp int

const (
	OpOnline cacheOp = iota
	OpReqData
	OpReConn
	OpOffline
	OpReLogin
	OpPackGameRole
)

type paramIn struct {
	acc       string
	op        cacheOp
	roleSesId uint32
}

type paramOut struct {
	data      *pb.MsgPlayerData
	roleSesId uint32
}

type RoleMgr struct {
	data map[string]*RoleCache
	ops  chan paramIn
	ret  chan paramOut
}

var Mgr4Role = newRoleMgr()

func newRoleMgr() *RoleMgr {
	m := &RoleMgr{
		data: make(map[string]*RoleCache),
		ops:  make(chan paramIn, 500),
		ret:  make(chan paramOut, 500),
	}
	go m.run()

	return m
}

//ReqRoleData	请求角色的数据
func (m *RoleMgr) ReqRoleData(msg *pb.MsgLogin, sesID uint32) {
	m.ops <- paramIn{acc: msg.Acc, op: OpReqData, roleSesId: sesID}
}

//ReqReConn	请求角重连
func (m *RoleMgr) ReqReConn(msg *pb.MsgReConn, sesID uint32) {
	m.ops <- paramIn{acc: msg.Acc, op: OpReConn, roleSesId: sesID}
}

//OnLoadRoleData	角色的数据加载完成
func (m *RoleMgr) OnLoadRoleData(data *pb.MsgPlayerData, roleSesId uint32) {
	m.ret <- paramOut{
		data:      data,
		roleSesId: roleSesId,
	}
}

//OnAccReqGameRoles	账号服启动，请求game上的角色信息
func (m *RoleMgr) AccReqGameRoles() {
	m.ops <- paramIn{op:OpPackGameRole}
}

//RoleOffline	角色下线
func (m *RoleMgr) RoleOffline(acc string) {
	m.ops <- paramIn{acc: acc, op: OpOffline}
}

//RoleOffline	角色异地登录
func (m *RoleMgr) RoleReLogin(acc string, newSesId uint32) {
	m.ops <- paramIn{acc: acc, op: OpReLogin, roleSesId: newSesId}
}

func (m *RoleMgr) run() {
	t := time.NewTicker(time.Minute)
	for { //不能访问cache里面的Role
		select {
		case p := <-m.ops:
			m.onOps(p)
		case ret := <-m.ret:
			m.onLoadRoleData(ret)
		case <-t.C:
			m.checkClear()
		}
	}
}

func (m *RoleMgr) onOps(p paramIn) {
	switch p.op {
	case OpReqData:
		m.reqData(p)
	case OpReConn:
		m.onReConn(p)
	case OpOffline:
		m.roleOffline(p.acc)
	case OpReLogin:
		m.roleReLogin(p)
	case OpPackGameRole:
		m.onAccReqGameRoles()
	}
}

func (m *RoleMgr) reqData(p paramIn) {
	v := m.data[p.acc]
	if v == nil {
		v = &RoleCache{}
		m.data[p.acc] = v
	}
	switch v.State {
	case RoleOnline:
		//异地登录
		ses := types.GetRoleSes(v.SesID)
		if ses != nil {
			//kick old when old is online
			types.PostToSes(v.SesID, types.Evt{
				Type: types.RoleReLogin,
				Data: &pb.MsgReLogin{SesID: p.roleSesId, Acc: p.acc},
			})
		} else {
			m._setOnline(v, p.roleSesId, p.acc)
		}
	case RoleOffline:
		//send data
		m._setOnline(v, p.roleSesId, p.acc)
	case RoleLoading:
		if time.Now().Unix()-v.StateTime < 5 {
			log.Warnf("%s login cd...", p.acc)
			return
		} else {
			v.setState(RoleLoading)
			//从数据库加载数据
			dbLoadRole(p.acc, p.roleSesId)
		}
	case RoleInit:
		v.setState(RoleLoading)

		//从数据库加载数据
		dbLoadRole(p.acc, p.roleSesId)
	}
}

func (m *RoleMgr) onReConn(p paramIn) {
	roleSes := types.GetRoleSes(p.roleSesId)
	if roleSes == nil {
		log.Warnf("%s reConn ses is nil", p.acc)
		return
	}

	if v, ok := m.data[p.acc]; ok && v != nil && v.State == RoleOffline {
		types.PostToSes(p.roleSesId, types.NewReConnEvt(v.Role))
	} else {
		types.PostToSes(p.roleSesId, types.Evt{
			Type: types.ReConn,
		})
	}
}

func (m *RoleMgr) onLoadRoleData(p paramOut) {
	v := m.data[p.data.Account]
	if v == nil {
		return
	}
	if v.State != RoleLoading {
		return
	}

	v.Role = NewRole(p.data)
	m._setOnline(v, p.roleSesId, p.data.Account)
}

func (m *RoleMgr) roleOffline(acc string) {
	if r, ok := m.data[acc]; ok {
		r.setState(RoleOffline)
		r.SesID = 0
	}
}

func (m *RoleMgr) _setOnline(v *RoleCache, sesId uint32, acc string) {
	roleSes := types.GetRoleSes(sesId)
	if roleSes == nil {
		log.Warnf("%s login ses is nil", acc)
		return
	}
	v.setState(RoleOnline)
	v.SesID = sesId

	types.PostToSes(sesId, types.NewLoadRoleSuccessEvt(v.Role))
}

func (m *RoleMgr) roleReLogin(p paramIn) {
	v := m.data[p.acc]
	if v == nil {
		return
	}

	m._setOnline(v, p.roleSesId, p.acc)
}

func (m *RoleMgr) checkClear() {
	now := time.Now().Unix()
	for k, v := range m.data {
		if v.State == RoleOffline && now-v.StateTime > int64(30*time.Minute.Seconds()) {
			network.SendToAcc(pb.MsgIDS2S_Gm2AccClearRole, &pb.MsgROnOffLine{Acc: k, Guid: v.Role.Guid})
			delete(m.data, k)
		}
	}
}

func (m *RoleMgr)onAccReqGameRoles(){
	msg := &pb.MsgGameRoles{Roles:make([]string, len(m.data))}

	i := 0
	for k, _ := range m.data {
		msg.Roles[i] = k
		i++
	}
	network.SendToAcc(pb.MsgIDS2S_Gm2AccGameRolesAck, msg)
}
