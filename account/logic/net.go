package logic

import (
	"account/setup"
	"com/log"
	"com/mq"
	"github.com/golang/protobuf/proto"
	"pb"
	"share"
	"strconv"
)

var msgSer *mq.NsqSession

func InitSerNet() {
	msgSer = mq.NewNsqSession(setup.Setup.NSQ, setup.Setup.NSQLookup)
	RegisteRoleMsgHandle()
	RegisteGameMsgHandle()
	msgSer.AddConsumer(share.AccountTopic, "acc1")
	msgSer.AddConsumer(share.BroadCastTopic, "acc1")
}

//Register 注册消息
func RegisterHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	msgSer.RegisterMsgHandle(msgID, cf, df)
}

func SendToGt(gateId uint16, msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.GateTopic+strconv.Itoa(int(gateId)), msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to gate err:%v", msgId, err)
	}
}

func SendToAllGm(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.GameAllTopic, msgId, msgData, 0)
	if err != nil {
		log.Warnf("send msg %d to game err:%v", msgId, err)
	}
}
