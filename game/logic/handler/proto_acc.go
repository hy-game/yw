package handler

import (
	"game/logic/role"
	"game/network"
	"github.com/golang/protobuf/proto"
	"pb"
)

//注意：必须写注释，这些函数在处理manage消息的协程调用
func initAccMsgHandle() {
	network.RegisterSrvHandle(pb.MsgIDS2S_Acc2GmGameRolesReq, func() proto.Message { return nil }, onAccInit)    //常规更新配置
}

func onAccInit(msgBase proto.Message, serId uint16) {
	role.Mgr4Role.AccReqGameRoles()
}
