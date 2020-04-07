package logic

import (
	"com/log"
	"com/mq"
	"gate/setup"
	"github.com/golang/protobuf/proto"
	"pb"
	share "share"
	"strconv"
)

var msgSer *mq.NsqSession

func InitSerNet() {
	msgSer = mq.NewNsqSession(setup.Setup.NSQ, setup.Setup.NSQLookup)
	registerServerMsg()
	msgSer.AddConsumer(share.GateTopic+strconv.Itoa(int(setup.Setup.Id)), "gt1")
	msgSer.AddConsumer(share.BroadCastTopic, "gt"+strconv.Itoa(int(setup.Setup.Id)))
}

//Register 注册消息
func RegisterHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	msgSer.RegisterMsgHandle(msgID, cf, df)
}

func SendToAcc(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.AccountTopic, msgId, msgData, uint16(setup.Setup.Id))
	if err != nil {
		log.Warnf("send msg %d to acc err:%v", msgId, err)
	}
}

func SendToBa(msgId pb.MsgIDS2S, msgData proto.Message) {
	err := msgSer.Send(share.BalanceTopic, msgId, msgData, uint16(setup.Setup.Id))
	if err != nil {
		log.Warnf("send msg %d to balance err:%v", msgId, err)
	}
}