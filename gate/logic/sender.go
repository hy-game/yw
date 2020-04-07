package logic

import (
	"gate/gnet"
	"github.com/golang/protobuf/proto"
	"pb"
)

//SendToCli	发送数据给客户端
func SendToCli(cliSesId uint32, msgId pb.MsgIDS2C, msg proto.Message) {
	s := gnet.GetCliSession(cliSesId)
	if s != nil {
		s.SendPB(uint16(msgId), msg)
	}
}
