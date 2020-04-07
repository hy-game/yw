package logic

import (
	"com/log"
	"gate/gnet"
	"github.com/golang/protobuf/proto"
	"pb"
)

func initCliMsg() {
	gnet.RegistryCliMsg(pb.MsgIDC2S_C2SInit, func() proto.Message { return &pb.MsgKeyExchange{} }, onNetInit)               //初始化
	gnet.RegistryCliMsg(pb.MsgIDC2S_C2SAccInfoReq, func() proto.Message { return &pb.MsgStr{} }, onAccInfoReq) //重连
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

func onAccInfoReq(msgBase proto.Message, s *gnet.Session) {
	//msg := msgBase.(*pb.MsgStr)
	s.SendPB(uint16(pb.MsgIDS2C_S2CAccInfoAck), &pb.MsgAccInfo{GateEndPoint:getGateEndPoint()})
}