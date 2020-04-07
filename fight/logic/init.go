package logic

import (
	"fight/configs"
	"fight/network"
	"pb"
)

// Init 初始化
func Init() {
	network.InitSerNet()
	initFcMsgHandle()
	initRoleMsgHandle()
	network.StartSerNet()
	network.SendToMa(pb.MsgIDS2S_Ft2MaConfigReq, nil)
	configs.Init()
}
