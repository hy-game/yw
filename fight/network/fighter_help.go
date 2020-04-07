package network

import (
	"fight/types"
	"pb"
)

func CreateFighter(initAttr *pb.MsgFighter, isMonster bool) *types.Fighter {
	one := &types.Fighter{}
	one.Init(initAttr)
	if isMonster { //只是迎合客户端，服务器处理状态机都应该一致
		one.FSM.AddState(pb.EFighterState_FS_Idle, &MsIdle{})
		one.FSM.AddState(pb.EFighterState_FS_Move, &MsMove{})
		one.FSM.AddState(pb.EFighterState_FS_Attack, &MsAttack{})
		one.FSM.AddState(pb.EFighterState_FS_BeHit, &MsBeHit{})
		one.FSM.AddState(pb.EFighterState_FS_Abnormal, &MsAbnormal{})
		one.FSM.AddState(pb.EFighterState_FS_Navigation, &MsNavigation{})
		one.FSM.AddState(pb.EFighterState_FS_Dead, &MsDead{})
		one.FSM.AddState(pb.EFighterState_FS_Victory, &MsVictory{})
		one.FSM.SetCurrentState(pb.EFighterState_FS_Idle)
	} else {
		one.FSM.AddState(pb.EFighterState_FS_Idle, &FsIdle{})
		one.FSM.AddState(pb.EFighterState_FS_Move, &FsMove{})
		one.FSM.AddState(pb.EFighterState_FS_Attack, &FsAttack{})
		one.FSM.AddState(pb.EFighterState_FS_BeHit, &FsBeHit{})
		one.FSM.AddState(pb.EFighterState_FS_Abnormal, &FsAbnormal{})
		one.FSM.AddState(pb.EFighterState_FS_Navigation, &FsNavigation{})
		one.FSM.AddState(pb.EFighterState_FS_Dead, &FsDead{})
		one.FSM.AddState(pb.EFighterState_FS_Victory, &FsVictory{})
		one.FSM.SetCurrentState(pb.EFighterState_FS_Idle)
	}
	return one
}
