package logic

import (
	"fcenter/network"
	"pb"

	"github.com/golang/protobuf/proto"
)

func initFtMsg() {
	network.NsqRegisterSrvHandle(pb.MsgIDS2S_Ft2FcBattleCreateAck, func() proto.Message { return &pb.MsgBattleStartData{} }, onBattleCreateAck)
	network.NsqRegisterSrvHandle(pb.MsgIDS2S_Ft2FcBattleFinish, func() proto.Message { return &pb.MsgBattleFinishData{} }, onBattleFinish)
}

func onBattleFinish(msgBase proto.Message, serID uint16) {
	msgCast := msgBase.(*pb.MsgBattleFinishData)
	gBtMgr.OnFinishBattle(msgCast)
}

func onBattleCreateAck(msgBase proto.Message, serID uint16) {
	msgCast := msgBase.(*pb.MsgBattleStartData)
	gBtMgr.OnCreateBattleAck(msgCast)
}
