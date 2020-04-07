package network

import (
	"com/log"
	"com/mq"
	"game/setup"
	"share"

	"pb"
	"strconv"

	"github.com/golang/protobuf/proto"
)

var sesSer *mq.NsqSession

//RegisterCtHandle 注册服务器间的消息处理函数
func RegisterSrvHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	sesSer.RegisterMsgHandle(msgID, cf, df)
}

//SendToCt 发消息到center
func SendToCt(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.CenterTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to ct err:%v", msgId, err)
	}
}

//SendToAcc 发消息到account
func SendToAcc(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.AccountTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to acc err:%v", msgId, err)
	}
}

//SendToFt 发消息到fighter
func SendToFt(fightId uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.FightTopic+strconv.Itoa(int(fightId)), msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to fight err:%v", msgId, err)
	}
}

// SendToFc :
func SendToFc(msgID pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.FightCenterTopic, msgID, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to fight err:%v", msgID, err)
	}
}

//Broadcast	服务器间消息广播
func Broadcast(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.BroadCastTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d broadcast err:%v", msgId, err)
	}
}

func SendToManage(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.ManageTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d broadcast err:%v", msgId, err)
	}
}
