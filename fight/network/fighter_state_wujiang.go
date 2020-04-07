package network

import (
	"fight/types"
	"pb"
)

type FsIdle struct {
	types.FsDefault
}

type FsMove struct {
	types.FsDefault
}

type FsAttack struct {
	types.FsDefault
}

type FsBeHit struct {
	types.FsDefault
}

type FsAbnormal struct {
	types.FsDefault
}

func (fs *FsAbnormal) OnEvent(owner *types.Fighter, event *pb.MsgFighterStateEvent) {
}

type FsNavigation struct {
	types.FsDefault
}

type FsDead struct {
	types.FsDefault
}

type FsVictory struct {
	types.FsDefault
}
