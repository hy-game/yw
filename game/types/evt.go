package types

import (
	"github.com/golang/protobuf/proto"
)

type GameEvent int

const (
	GEInit        = iota
	LoadRoleData  //加载角色数据成功
	ReConn        //重连
	KickOut       //踢角色
	RoleReLogin   //异地登录
	SecLoop       //每秒循环调用
	ForwardToRole //收到转发的消息

	//逻辑
	GMOrder         //GM命令
	BattleCreateAck //战斗创建返回
	BattleFinish    //战斗结算
	RecvMail		//收到邮件
)

type Evt struct {
	Type GameEvent
	Data proto.Message
	role *Role //只在重连和登录的时候使用
}

//PostEvt post事件给指定guid的角色
func PostEvt(roleGuid uint32, e Evt) bool{
	sessID, ok := guid2SesID.Load(roleGuid)
	if ok {
		return PostToSes(sessID.(uint32), e)
	}else{
		return false
	}
}

//PostEvtToAll	投递事件给所有在线连接
func PostEvtToAllOnline(e Evt) {
	roleSes.Range(func(k interface{}, v interface{}) bool {
		ses := v.(*Session)
		if ses != nil {
			ses.postEvt(e)
		}
		return true
	})
}

//PostToSes	投递事件给指定的连接
func PostToSes(sessID uint32, evt Evt) bool {
	ses := GetRoleSes(sessID)
	if ses != nil {
		ses.postEvt(evt)
		return true
	} else {
		return false
	}
}

//NewLoadRoleSuccessEvt	创建一个加载角色成功事件
func NewLoadRoleSuccessEvt(r *Role) Evt {
	return Evt{Type: LoadRoleData, role: r}
}

//NewReConnEvt 重连事件
func NewReConnEvt(r *Role) Evt {
	return Evt{Type: ReConn, role: r}
}

//GetRoleForConnEvt	重连和登录时需要role
func GetRoleForConnEvt(e *Evt) *Role {
	return e.role
}
