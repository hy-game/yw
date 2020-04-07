package network

import (
	"com/log"
	"com/mq"
	"pb"
	share "share"
	"strconv"

	"github.com/golang/protobuf/proto"
)

var msgSer *mq.NsqSession

// NsqInit :
func NsqInit(nsqdAddr string, nsqLkup []string) {
	msgSer = mq.NewNsqSession(nsqdAddr, nsqLkup)
}

// NsqListen :
func NsqListen() {
	msgSer.AddConsumer(share.FightCenterTopic, "fc")
}

// NsqRegisterSrvHandle : 注册消息到默认handler
func NsqRegisterSrvHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	msgSer.RegisterMsgHandle(msgID, cf, df)
}

// NsqSendToFt :
func NsqSendToFt(ftID uint16, msgID pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.FightTopic+strconv.Itoa(int(ftID)), msgID, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to ft %d err:%v", msgID, ftID, err)
	}
}

// NsqSendToRandFt :
func NsqSendToRandFt(msgID pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.FightTopic, msgID, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to ft %d err:%v", msgID, -1, err)
	}
}

// NsqSendToGm :
func NsqSendToGm(gmID uint16, msgID pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.GameTopic+strconv.Itoa(int(gmID)), msgID, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to gm %d err:%v", msgID, gmID, err)
	}
}
