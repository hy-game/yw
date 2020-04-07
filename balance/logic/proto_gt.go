package logic

import (
	"github.com/golang/protobuf/proto"
	"balance/network"
	"pb"
)

//RegisteMsgHandle 注册消息处理函数
func RegisteMsgHandle() {
	 network.RegisterHandle(pb.MsgIDS2S_GtHeartBeat, func() proto.Message { return &pb.MsgGtHeartBeat{} }, onGateInit)
}

func onGateInit(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgGtHeartBeat)
	if msg == nil {
		return
	}
	updateGateCli(msg.EndPoint, msg.RoleCnt)
}
