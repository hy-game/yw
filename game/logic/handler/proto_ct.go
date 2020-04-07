/*
center 发过来的消息处理
*/
package handler

import (
	"game/network"
	"game/types"
	"github.com/golang/protobuf/proto"
	"pb"
)

func initCenterMsgHandle() {
	network.RegisterSrvHandle(pb.MsgIDS2S_CtForwardToRole, func() proto.Message { return &pb.MsgForwardToRole{} }, onForwardMsg)		//转发消息
	network.RegisterSrvHandle(pb.MsgIDS2S_Ct2GmSendMail, func() proto.Message { return &pb.MsgMail{} }, onRecvMail)					//收到邮件
}

func onForwardMsg(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgForwardToRole)
	if msg == nil {
		return
	}
	types.PostEvt(msg.RoleGuid, types.Evt{
		Type: types.ForwardToRole,
		Data: msg,
	})
}

func onRecvMail(msgBase proto.Message, serId uint16)  {
	msg := msgBase.(*pb.MsgMail)
	if msg == nil {
		return
	}
	types.PostEvt(msg.RoleGuid, types.Evt{
		Type: types.RecvMail,
		Data: msg,
	})
}