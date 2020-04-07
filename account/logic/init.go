package logic

import (
	"pb"
	"sync"
)
var WaitGameInfo sync.WaitGroup

func Init(){
	WaitGameInfo.Add(1)
	ReqGameRoles()
	WaitGameInfo.Wait()
}

func ReqGameRoles(){
	SendToAllGm(pb.MsgIDS2S_Acc2GmGameRolesReq, nil)
}