package logic

import (
	"fcenter/network"

	"pb"

	"github.com/golang/protobuf/proto"
)

func initGmMsg() {
	network.NsqRegisterSrvHandle(pb.MsgIDS2S_Gm2FcBattleSearchReq, func() proto.Message { return &pb.MsgUint{} }, onBattleSearchReq)
	network.NsqRegisterSrvHandle(pb.MsgIDS2S_Gm2FcBattleCreateReq, func() proto.Message { return &pb.MsgBattleStartData{} }, onBattleCreateReq)
}

func onBattleSearchReq(msgBase proto.Message, serID uint16) {

}

func onBattleCreateReq(msgBase proto.Message, serID uint16) {
	msgCast := msgBase.(*pb.MsgBattleStartData)
	msgCast.BattleCreateServer = (uint32)(serID)
	if gBtMgr.OnCreateBattle(msgCast) {
		network.NsqSendToRandFt(pb.MsgIDS2S_Fc2FtBattleCreateReq, msgCast)
	}
}
