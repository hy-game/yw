package handler

import (
	"github.com/golang/protobuf/proto"
	"manage/logic/configs"
	"manage/network"
	"manage/web/controllers"
	"pb"
	share "share"
)

//registeSrvHandle	注册消息处理，必须写注释
func RegisteSrvHandle() {
	network.RegisterHandle(pb.MsgIDS2S_Gm2MaConfigReq, func() proto.Message { return nil }, onGameReqConfig)               //game请求配置
	network.RegisterHandle(pb.MsgIDS2S_Gm2MaConfigYYAct, func() proto.Message { return nil }, onGameYYActConfig)           //game请求运营活动
	network.RegisterHandle(pb.MsgIDS2S_Gm2MaRoleData, func() proto.Message { return &pb.MsgPlayerData{} }, onGameRoleData) //角色数据

	network.RegisterHandle(pb.MsgIDS2S_Ct2MaConfigReq, func() proto.Message { return nil }, onCenterReqConfig)             //center请求配置
	network.RegisterHandle(pb.MsgIDS2S_Ct2MaConfigYYAct, func() proto.Message { return nil }, onCenterYYActConfig)         //center请求运营活动
	network.RegisterHandle(pb.MsgIDS2S_Ct2MaRLData, func() proto.Message { return &pb.MsgRankListPack{} }, onCenterRLData) //center请求排行榜数据

	network.RegisterHandle(pb.MsgIDS2S_Ft2MaConfigReq, func() proto.Message { return nil }, onFightReqConfig) //center请求配置

}

func onGameReqConfig(msgBase proto.Message, serId uint16) {
	configs.BroadCastToServer(share.GameTopic, serId, false)
}

func onCenterReqConfig(msgBase proto.Message, serId uint16) {
	configs.BroadCastToServer(share.CenterTopic, serId, false)
}

func onGameYYActConfig(msgBase proto.Message, serId uint16) {
	configs.BroadCastToServer(share.GameTopic, serId, true)
}

func onCenterYYActConfig(msgBase proto.Message, serId uint16) {
	configs.BroadCastToServer(share.CenterTopic, serId, true)
}

func onFightReqConfig(msgBase proto.Message, serId uint16) {
	configs.BroadCastToServer(share.FightTopic, serId, false)
}

func onGameRoleData(msgBase proto.Message, serId uint16) {
	data, ok := msgBase.(*pb.MsgPlayerData)
	if !ok || data == nil {
		return
	}
	ch, ok := controllers.WebMgr.Param.(chan *pb.MsgPlayerData)
	if ok {
		ch <- data
		close(ch)
	}
}

func onCenterRLData(msgBase proto.Message, serId uint16) {
	data, ok := msgBase.(*pb.MsgRankListPack)
	if !ok || data == nil {
		return
	}
	ch, ok := controllers.WebMgr.Param.(chan *pb.MsgRankListPack)
	if ok {
		ch <- data
		close(ch)
	}
}
