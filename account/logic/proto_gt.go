package logic

import (
	"github.com/golang/protobuf/proto"
	"pb"
)

//RegisteMsgHandle 注册消息处理函数
func RegisteRoleMsgHandle() {
	RegisterHandle(pb.MsgIDS2S_Gt2AccLogin, func() proto.Message { return &pb.MsgLogin{} }, onLogin) //登录
	//	RegisterHandle(pb.MsgIDS2S_Gt2AccReConn, func() proto.Message { return &pb.MsgReConnToAcc{} }, onReConn)	//
}

func onLogin(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgLogin)
	if msg == nil {
		return
	}

	loginMgr.Login(msg, serId)
}

//
//func onReConn(msgBase proto.Message, serId uint16) {
//	msg := msgBase.(*pb.MsgReConnToAcc)
//	if msg == nil {
//		return
//	}
//
//	loginMgr.ReConn(msg, serId)
//}
