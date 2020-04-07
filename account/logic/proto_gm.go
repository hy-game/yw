package logic

import (
	"github.com/golang/protobuf/proto"
	"pb"
)

func RegisteGameMsgHandle() {
	RegisterHandle(pb.MsgIDS2S_GmHeartBeat, func() proto.Message { return &pb.MsgGameHeartBeat{} }, onGameHeartBeat) //game心跳
	RegisterHandle(pb.MsgIDS2S_Gm2AccClearRole, func() proto.Message { return &pb.MsgROnOffLine{} }, onGameClearRole)
	RegisterHandle(pb.MsgIDS2S_Gm2AccGameRolesAck, func() proto.Message { return &pb.MsgGameRoles{} }, onGameRoles)
}

func onGameHeartBeat(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgGameHeartBeat)
	if msg == nil {
		return
	}

	loginMgr.PostEvt(evtParam{
		op:    OpGameInfo,
		gmSrv: msg,
		serId: serId,
	})
}

func onGameClearRole(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgROnOffLine)
	if msg == nil {
		return
	}

	loginMgr.PostEvt(evtParam{
		op:       OpRoleClear,
		roleInfo: msg,
		serId:    serId,
	})
}

func onGameRoles(msgBase proto.Message, serId uint16) {
	msg := msgBase.(*pb.MsgGameRoles)
	if msg == nil {
		return
	}

	loginMgr.PostEvt(evtParam{
		op:       	OpGameRoles,
		gmRoles: 	msg,
		serId:    	serId,
	})
}