package network

import (
	"com/log"
	"com/mq"
	"github.com/golang/protobuf/proto"
	"manage/setup"
	"pb"
	share "share"
	"strconv"
)

var sesSer *mq.NsqSession

//InitSerNet	初始化服务器间的网络
func InitSerNet() {
	sesSer = mq.NewNsqSession(setup.Setup.NSQ, setup.Setup.NSQLookup)
}

//StartSerNet	开始接收服务器间消息
func StartSerNet() {
	sesSer.AddConsumer(share.ManageTopic, "manage")
}

//Register 注册消息
func RegisterHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	sesSer.RegisterMsgHandle(msgID, cf, df)
}

func Broadcast(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.BroadCastTopic, msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d broadcast err:%v", msgId, err)
	}
}

func SendToGame(gameId uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.GameTopic+strconv.Itoa(int(gameId)), msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to game err:%v", msgId, err)
	}
}

func SendToGameAll(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.GameAllTopic, msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to game all err:%v", msgId, err)
	}
}

func SendToCenter(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.CenterTopic, msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to center err:%v", msgId, err)
	}
}

func SendToFightAll(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.FightTopic, msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to fight all err:%v", msgId, err)
	}
}

func SendToFight(ftid uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	err := sesSer.Send(share.FightTopic+strconv.Itoa(int(ftid)), msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to fight err:%v", msgId, err)
	}
}
