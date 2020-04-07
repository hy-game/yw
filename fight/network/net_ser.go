package network

import (
	"com/log"
	"com/mq"
	"fight/setup"
	"pb"
	share "share"
	"strconv"

	"github.com/golang/protobuf/proto"
)

var msgSer *mq.NsqSession

//InitSerNet	初始化
func InitSerNet() {
	msgSer = mq.NewNsqSession(setup.Setup.NSQ, setup.Setup.NSQLookup)

}

//StartSerNet	开始接收消息
func StartSerNet() {
	msgSer.AddConsumer(share.FightTopic+strconv.Itoa(int(setup.Setup.ID)), "ft")
	msgSer.AddConsumer(share.FightTopic, "ft")
	msgSer.AddConsumer(share.FightAllTopic, "ft"+strconv.Itoa(int(setup.Setup.ID)))
}

//Register 注册消息到默认handler
func RegisterSrvHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	msgSer.RegisterMsgHandle(msgID, cf, df)
}

func SendToCt(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.CenterTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to ct err:%v", msgId, err)
	}
}

// SendToFc :
func SendToFc(msgID pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.FightCenterTopic, msgID, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to ct err:%v", msgID, err)
	}
}

func SendToGm(gameId uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.GameTopic+strconv.Itoa(int(gameId)), msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to game err:%v", msgId, err)
	}
}

//Broadcast	服务器间消息广播
func Broadcast(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.BroadCastTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d broadcast err:%v", msgId, err)
	}
}

//Fight给Manage发消息
func SendToMa(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.ManageTopic, msgId, msgData, uint16(setup.Setup.ID))
	if err != nil {
		log.Warnf("send msg %d to game err:%v", msgId, err)
	}
}
