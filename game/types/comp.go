package types

import (
	"pb"
)

type TypeComp int

const (
	TCHero TypeComp = iota //英雄
	TCItem
	TCFB
	TCMail

	TCMax //最大值
)

type IComp interface {
}

//ICompChgData 如果在给客户端发送角色所有数据前，需要做些修改，就实现该接口
type ICompChgData interface {
	ChangeForCli(data *pb.MsgPlayerData, r *Role) //数据发给客户端前可以做修改
}

//ICompSecLoop 如果需要每秒update，就实现该接口
type ICompSecLoop interface {
	SecLoop(r *Role) //每秒更新
}

//ICompDataReset 如果需要每天定时重置数据，就实现该接口
type ICompDataReset interface {
	OnDataReset(r *Role) //每秒更新
}
