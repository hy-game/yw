package logic

import (
	"com/log"
	"gate/gnet"
	"gate/service"
	"github.com/golang/protobuf/proto"
	"pb"
)

//注意：必须写注释，这些函数在处理account消息的协程调用
func registerServerMsg() {
	RegisterHandle(pb.MsgIDS2S_Acc2GtLoginAck, func() proto.Message { return &pb.MsgLoginAck{} }, onLoginAck)        //登录验证返回
	RegisterHandle(pb.MsgIDS2S_GmHeartBeat, func() proto.Message { return &pb.MsgGameHeartBeat{} }, onGameHeartBeat) //game 心跳
	RegisterHandle(pb.MsgIDS2S_FtHeartBeat, func() proto.Message { return &pb.MsgFtHeartBeat{} }, onFtHeartBeat)     //game 心跳
}

//-------------------账号服务器消息-----------
func onLoginAck(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgLoginAck)
	if msg == nil {
		return
	}

	if msg.Ret != pb.LoginCode_LCSuccess {
		msgAck := &pb.MsgLoginForCli{Ret: msg.Ret}
		SendToCli(msg.Data.SesOnGt, pb.MsgIDS2C_S2CLoginAck, msgAck)
	} else {
		onCheckSuccess(msg)
	}
}

func onCheckSuccess(msg *pb.MsgLoginAck) {
	s := gnet.GetCliSession(msg.Data.SesOnGt)
	if s == nil {
		log.Debug("s == nil when login check success")
		return
	}
	s.PostEvt(gnet.EvtLoginSuccess, msg)
}

//------------------------game直接发送到gate的消息-----------------------------
func onGameHeartBeat(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgGameHeartBeat)
	if msg == nil {
		return
	}
	service.Add("game", serId, msg.EndPoint)
}

//------------------------fight直接发送到gate的消息-----------------------------
func onFtHeartBeat(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgFtHeartBeat)
	if msg == nil {
		return
	}
	service.Add("fight", serId, msg.EndPoint)
}
