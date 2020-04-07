package logic

import (
	"fcenter/network"
	"pb"
	"sync"
)

type battleManager struct {
	btguid uint32
	mtx    sync.Mutex
	link   map[uint32]uint32 //玩家战斗连接
	list   map[uint32]*pb.MsgBattleStartData
}

var gBtMgr = &battleManager{
	link: make(map[uint32]uint32),
	list: make(map[uint32]*pb.MsgBattleStartData),
}

func (btm *battleManager) OnCreateBattle(start *pb.MsgBattleStartData) bool {
	btm.mtx.Lock()
	defer btm.mtx.Unlock()
	for _, v := range start.Roles {
		if v.Guid == 0 {
			return false
		}
		_, isfind := btm.link[v.Guid]
		if isfind {
			//todo: reconnect
			continue
			//return false
		}
	}
	btm.btguid++
	start.BattleGuid = btm.btguid
	btm.list[start.BattleGuid] = start
	for _, v := range start.Roles {
		btm.link[v.Guid] = btm.btguid
	}
	return true
}

func (btm *battleManager) QueryRoleBattle(roleGuid uint32) (ftId uint32, btmGuid uint32) {
	btm.mtx.Lock()
	defer btm.mtx.Unlock()

	if btGuid, ok := btm.link[roleGuid]; ok {
		if bt := btm.list[btGuid]; bt != nil {
			btmGuid = bt.BattleGuid
			ftId = bt.FtID
		}
	}
	return
}

func (btm *battleManager) OnCreateBattleAck(start *pb.MsgBattleStartData) bool {
	btm.mtx.Lock()
	defer btm.mtx.Unlock()
	battle := btm.list[start.BattleGuid]
	if battle == nil {
		return false
	}
	battle.FtID = start.FtID
	if start.BattleCreateServer > 0 { //create by gm
		network.NsqSendToGm((uint16)(start.BattleCreateServer), pb.MsgIDS2S_Fc2GmBattleCreateAck, start)
	} else { //create by room

	}
	return true
}

func (btm *battleManager) OnFinishBattle(finish *pb.MsgBattleFinishData) bool {
	btm.mtx.Lock()
	defer btm.mtx.Unlock()
	battle := btm.list[finish.BattleGuid]
	if battle == nil {
		return false
	}
	finish.BattleID = battle.BattleID
	for _, v := range battle.Roles {
		finish.Guid = v.Guid
		network.NsqSendToGm((uint16)(v.ServerID), pb.MsgIDS2S_Fc2GmBattleFinish, finish)
	}
	return true
}
