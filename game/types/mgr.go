package types

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var shutDown = make(chan struct{})
var roleSes sync.Map
var roleCnt int32
var guid2SesID sync.Map

//---------------------------------------
func addRoleSes(sesId uint32, ses *Session) {
	roleSes.Store(sesId, ses)
	atomic.AddInt32(&roleCnt, 1)
}

func delRoleSes(sesId uint32) {
	roleSes.Delete(sesId)
	atomic.AddInt32(&roleCnt, -1)
}

//GetRoleSes 获取指定id的网络会话
func GetRoleSes(sesId uint32) *Session {
	if ses, ok := roleSes.Load(sesId); ok {
		return ses.(*Session)
	}
	return nil
}

//GetRoleSesCnt 获取当前网络会话数
func GetRoleSesCnt() int32 {
	return atomic.LoadInt32(&roleCnt)
}

//---------------------------------------
//BindRoleSess	绑定角色guid和sesID
func BindRoleSess(role *Role, ses *Session) {
	role.Ses = ses
	ses.Role = role
	ses.desc = fmt.Sprintf("ses[%d]_%s[%d]", ses.ID, role.Data.Name ,role.Guid)
	guid2SesID.Store(role.Guid, ses.ID)
}

//UnBindRoleSess	删除已绑定的角色guid
func UnBindRoleSess(role *Role, ses *Session) {
	role.Ses = nil
	ses.Role = nil
	guid2SesID.Delete(role.Guid)
}
