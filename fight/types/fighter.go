package types

import (
	"github.com/golang/protobuf/proto"
	"pb"
)

type Fighter struct {
	InitAttr *pb.MsgFighter //初始属性
	RealAttr *pb.MsgFighter //实际属性
	FSM      *FighterStateMachine
}

func (ft *Fighter) Init(initAttr *pb.MsgFighter) {
	ft.InitAttr = proto.Clone(initAttr).(*pb.MsgFighter)
	ft.RealAttr = proto.Clone(initAttr).(*pb.MsgFighter)
	ft.FSM = &FighterStateMachine{}
	ft.FSM.Init(ft)
}
