package network

import (
	"balance/setup"
	"com/mq"
	"com/util"
	"github.com/golang/protobuf/proto"
	"pb"
	"share"
)

var msgSer = make([]*mq.NsqSession, 0)

func InitSerNet() {
	for _, v := range setup.Setup.Nsq {
		msgSer = append(msgSer, mq.NewNsqSession(v.NSQAddr, v.NSQLookup))
	}
}

func StartSerNet() {
	for i := 0; i != len(msgSer); i++{
		msgSer[i].AddConsumer(share.BalanceTopic, "bl"+util.ToString(i))
	}
}

//Register 注册消息
func RegisterHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	for _, v := range msgSer{
		v.RegisterMsgHandle(msgID, cf, df)
	}
}
