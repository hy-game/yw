package handler

import (
	"center/configs"
	"center/network"
	"github.com/golang/protobuf/proto"
	"pb"
)

func RegisteManageMsgHandle() {
	network.RegisterHandle(pb.MsgIDS2S_MsgBroadcastCfgs, func() proto.Message { return &pb.MsgOriginalCfgs{} }, onManageCfg)
	network.RegisterHandle(pb.MsgIDS2S_Ma2CtGMOrder, func() proto.Message { return &pb.MsgGMTask{} }, onManageGMOrder)
	network.RegisterHandle(pb.MsgIDS2S_MsgBroadcastYYAct, func() proto.Message { return &pb.MsgOriginalCfgs{} }, onManageYYAct)
}

func onManageCfg(msgBase proto.Message, serId uint16) {
	msg, ok := msgBase.(*pb.MsgOriginalCfgs)
	if !ok || msg == nil {
		return
	}
	go configs.UpdateCfg(msg)
}

func onManageGMOrder(msgBase proto.Message, serId uint16) {
	msg, ok := msgBase.(*pb.MsgGMTask)
	if !ok || msg == nil {
		return
	}
}

func onManageYYAct(msgBase proto.Message, serId uint16) {
	msg, ok := msgBase.(*pb.MsgOriginalCfgs)
	if !ok || msg == nil {
		return
	}
	go configs.UpdateYYAct(msg)
}
