package network

import (
	"center/setup"
	"com/log"
	"com/mq"
	"github.com/golang/protobuf/proto"
	"pb"
	share "share"
	"strconv"
)

var msgSer *mq.NsqSession

func InitSerNet() {
	msgSer = mq.NewNsqSession(setup.Setup.NSQ, setup.Setup.NSQLookup)
}

func StartSerNet() {
	msgSer.AddConsumer(share.CenterTopic, "ct1")
}

//Register 注册消息
func RegisterHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	msgSer.RegisterMsgHandle(msgID, cf, df)
}

func SendToGm(gameId uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	if gameId == 0{
		return
	}
	err := msgSer.Send(share.GameTopic+strconv.Itoa(int(gameId)), msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to game err:%v", msgId, err)
	}
}

func SendToFt(fightId uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.FightTopic+strconv.Itoa(int(fightId)), msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to fight err:%v", msgId, err)
	}
}

func SendToManage(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.ManageTopic, msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to manage err:%v", msgId, err)
	}
}
