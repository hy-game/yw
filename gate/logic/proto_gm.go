package logic

import (
	"gate/gnet"
	"pb"

	"github.com/golang/protobuf/proto"
)

//initGameMsg game通过流发来的消息处理，这些函数在玩家主协程调用的
func initGameMsg() {
	gnet.RegistryGameMsg(pb.MsgIDS2C_Gm2GtKickRole, func() proto.Message { return nil }, onGameKickOutRole)  //踢玩家下线
	gnet.RegistryGameMsg(pb.MsgIDS2C_Ft2GtKickRole, func() proto.Message { return nil }, onFightKickOutRole) //踢玩家下线
}

func onGameKickOutRole(msgBase proto.Message, ses *gnet.Session) {
	ses.Close()
}

func onFightKickOutRole(msgBase proto.Message, ses *gnet.Session) {
	ses.CloseToFt()
}
