package logic

import (
	"center/logic/handler"
	"center/network"
)

func Init() {
	network.InitSerNet()
	handler.RegisteGameMsgHandle()
	handler.RegisteManageMsgHandle()
	network.StartSerNet()
}

func InitAfterDB(){

}
