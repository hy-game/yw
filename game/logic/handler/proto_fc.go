/*
fighter发来的消息处理
*/
package handler

import (
	"game/network"
	"game/types"
	"pb"

	"github.com/golang/protobuf/proto"
)

//注意：写注释
func initFtMsgHandle() {
	network.RegisterSrvHandle(pb.MsgIDS2S_Fc2GmBattleCreateAck, func() proto.Message { return &pb.MsgBattleStartData{} }, onBattleCreateAck)
	network.RegisterSrvHandle(pb.MsgIDS2S_Fc2GmBattleFinish, func() proto.Message { return &pb.MsgBattleFinishData{} }, onBattleFinish)
}

func onBattleFinish(msgBase proto.Message, serID uint16) {
	msgCast := msgBase.(*pb.MsgBattleFinishData)
	e := types.Evt{Type: types.BattleFinish, Data: msgCast}
	if !types.PostEvt(msgCast.Guid, e){
		types.PostOfflineOp(msgCast.Guid, e)
	}
}

func onBattleCreateAck(msgBase proto.Message, serID uint16) {
	msgCast := msgBase.(*pb.MsgBattleStartData)
	msgSend := &pb.MsgBattleCreateAck{BtStart: msgCast}
	for _, v := range msgCast.Roles {
		roleSes := types.GetRoleSes(v.SesID)
		if roleSes != nil {
			msgSend.Guid = v.Guid
			roleSes.Send(pb.MsgIDS2C_BattleCreateAck, msgSend)
		}
	}
}
