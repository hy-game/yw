/*
服务器间的广播消息处理
*/
package handler

import (
	"game/configs"
	"game/network"
	"game/types"
	"github.com/golang/protobuf/proto"
	"pb"
)

//注意：必须写注释，这些函数在处理manage消息的协程调用
func initManageMsgHandle() {
	network.RegisterSrvHandle(pb.MsgIDS2S_MsgBroadcastCfgs, func() proto.Message { return &pb.MsgOriginalCfgs{} }, onManageCfg)    //常规更新配置
	network.RegisterSrvHandle(pb.MsgIDS2S_Ma2GmGMOrder, func() proto.Message { return &pb.MsgGMTask{} }, onManageGMOrder)          //GM命令
	network.RegisterSrvHandle(pb.MsgIDS2S_MsgBroadcastYYAct, func() proto.Message { return &pb.MsgOriginalCfgs{} }, onManageYYAct) //运营活动更新
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

	e := types.Evt{
		Type: types.GMOrder,
		Data: msg,
	}
	if !types.PostEvt(uint32(msg.MPlayerId), e){
		types.PostOfflineOp(uint32(msg.MPlayerId), e)
	}
}

func onManageYYAct(msgBase proto.Message, serId uint16) {
	msg, ok := msgBase.(*pb.MsgOriginalCfgs)
	if !ok || msg == nil {
		return
	}
	go configs.UpdateYYAct(msg)
}
