package logic

import (
	"com/log"
	"fight/configs"
	"fight/network"
	"fight/setup"
	"pb"

	"github.com/golang/protobuf/proto"
)

func initFcMsgHandle() {
	network.RegisterSrvHandle(pb.MsgIDS2S_Fc2FtBattleCreateReq, func() proto.Message { return &pb.MsgBattleStartData{} }, onBattleCreateReq) //请求创建战斗场景
	network.RegisterSrvHandle(pb.MsgIDS2S_MsgBroadcastCfgs, func() proto.Message { return &pb.MsgOriginalCfgs{} }, onManageCfg)              //常规更新配置
}

func onBattleCreateReq(msgBase proto.Message, serID uint16) {
	msgCast := msgBase.(*pb.MsgBattleStartData)
	if msgCast == nil {
		return
	}

	rg := network.MgrForRegion.Create(msgCast)
	if rg == nil {
		log.Warnf("create region err:%s", msgCast.String())
		return
	}

	network.MgrForRegion.Add(rg)
	msgCast.FtID = setup.Setup.ID
	network.SendToFc(pb.MsgIDS2S_Ft2FcBattleCreateAck, msgCast)
}

func onManageCfg(msgBase proto.Message, serId uint16) {
	msg, ok := msgBase.(*pb.MsgOriginalCfgs)
	if !ok || msg == nil {
		return
	}
	go configs.UpdateCfg(msg)
}
