package logic

import (
	"manage/logic/handler"
	"manage/network"
)

//Init 逻辑初始化
func Init() {
	network.InitSerNet()
	handler.RegisteSrvHandle()
	network.StartSerNet()
}
