package handler

import (
	"center/network"
	"github.com/golang/protobuf/proto"
	"pb"
)

//RegisteMsgHandle 注册消息处理函数
func RegisteGameMsgHandle() {
	network.RegisterHandle(pb.MsgIDS2S_Gm2CtLogin, func() proto.Message { return &pb.MsgKeyValueU{} }, onRoleLogin)
	network.RegisterHandle(pb.MsgIDS2S_Gm2CtOffline, func() proto.Message { return &pb.MsgKeyValueU{} }, onRoleOffline)
}

func onRoleLogin(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgKeyValueU)
	if msg == nil {
		return
	}
	network.OnRoleOnline(msg.Key, uint16(msg.Value))
}

func onRoleOffline(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgKeyValueU)
	if msg == nil {
		return
	}
	network.OnRoleOffline(msg.Key)
}
