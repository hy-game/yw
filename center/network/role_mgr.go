package network

import (
	"com/log"
	"github.com/golang/protobuf/proto"
	"pb"
	"sync"
)

var	RoleInGame sync.Map

func OnRoleOnline(roleGuid uint32, gameID uint16){
	RoleInGame.Store(roleGuid, gameID)
}

func OnRoleOffline(roleGuid uint32){
	RoleInGame.Delete(roleGuid)
}

//GetGameIDByRole 获取角色所在的gameID
func GetGameIDByRole(roleGuid uint32)uint16{
	if gmID, ok := RoleInGame.Load(roleGuid); ok {
		return gmID.(uint16)
	}else{
		return 0
	}
}

//SendToRole	经过game转发消息给角色， 线程安全
func SendToRole(roleGuid uint32, msgID pb.MsgIDS2C, msg proto.Message){
	msgSend := &pb.MsgForwardToRole{MsgID:uint32(msgID)}
	if msg != nil{
		b, err := proto.Marshal(msg)
		if err != nil{
			log.Warnf("marshal err:%v when SendToRole", err)
			return
		}
		msgSend.Data = b
	}
	SendToGm(GetGameIDByRole(roleGuid), pb.MsgIDS2S_CtForwardToRole, msgSend)
}

