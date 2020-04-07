package logic

import (
	"com/log"
	"gate/gnet"
	"pb"
	"time"

	"github.com/golang/protobuf/proto"
)

//注意：必须写注释，这些函数在处理每个角色消息的协程调用，每个角色一个协程
func initCliMsg() {
	gnet.RegistryCliMsg(pb.MsgIDC2S_C2SHeartBeat, func() proto.Message { return &pb.MsgHeartBeat{} }, onHeartBeat)          //心跳
	gnet.RegistryCliMsg(pb.MsgIDC2S_C2SInit, func() proto.Message { return &pb.MsgKeyExchange{} }, onNetInit)               //初始化
	gnet.RegistryCliMsg(pb.MsgIDC2S_C2SLogin, func() proto.Message { return &pb.MsgLogin{} }, onLogin)                      //登录
	gnet.RegistryCliMsg(pb.MsgIDC2S_C2SReConn, func() proto.Message { return &pb.MsgReConn{} }, onReConn)                   //重连
	gnet.RegistryCliMsg(pb.MsgIDC2S_BattleEnterReq, func() proto.Message { return &pb.MsgBattleEnterReq{} }, onEnterBattle) //重连
	gnet.RegistryCliMsg(pb.MsgIDC2S_BattleLeaveReq, func() proto.Message { return &pb.MsgBattleLeaveReq{} }, onLeaveBattle) //重连
}

func onLeaveBattle(msg proto.Message, s *gnet.Session) {
	msgCast := msg.(*pb.MsgBattleLeaveReq)
	s.OnLeaveBattle(msgCast)
}

func onEnterBattle(msg proto.Message, s *gnet.Session) {
	msgCast := msg.(*pb.MsgBattleEnterReq)
	s.OnEnterBattle(msgCast)
}

func onHeartBeat(msgBase proto.Message, ses *gnet.Session) {
	ses.SendPB(uint16(pb.MsgIDS2C_S2CHeartBeat), &pb.MsgHeartBeat{Time: time.Now().UnixNano()})
}

func onNetInit(msgBase proto.Message, ses *gnet.Session) {
	log.Tracef("%s recv init ", ses.String())
	//msg, ok := msgBase.(*pb.MsgKeyExchange)
	//if !ok {
	//	log.Warnf("msg type err %s", util.FuncCaller(1))
	//	return
	//}
	//var err error
	//
	//
	//enCodeA, enCodea := dh.Exchange()
	//enCodeKey := dh.GetKey(enCodea, big.NewInt(msg.Decode))
	//
	//ses.EnCoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v", enCodeKey)))
	//if err != nil {
	//	log.Warnf("make rc4 key err:%v", err)
	//	return
	//}
	//
	//deCodeA, deCodea := dh.Exchange()
	//deCodeKey := dh.GetKey(deCodea, big.NewInt(msg.Encode))
	//ses.DeCoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v", deCodeKey)))
	//if err != nil {
	//	log.Warnf("make rc4 key err:%v", err)
	//	return
	//}

	ses.SendPB(uint16(pb.MsgIDS2C_S2CInit), &pb.MsgKeyExchange{
		Encode: 0, // deCodeA.Int64(),
		Decode: 0, //enCodeA.Int64(),
	})
}

//todo 客户端点击登录需要加个CD,并且5秒超时，不能一直等待
func onLogin(msgBase proto.Message, ses *gnet.Session) {
	msg := msgBase.(*pb.MsgLogin)
	if msg == nil {
		return
	}

	msg.SesOnGt = ses.ID()
	SendToAcc(pb.MsgIDS2S_Gt2AccLogin, msg)
}

func onReConn(msgBase proto.Message, ses *gnet.Session) {
	if ses.StreamGm != nil {
		return
	}
	msg := msgBase.(*pb.MsgReConn)
	if msg == nil {
		return
	}

	ses.ReConn(msg)
}
