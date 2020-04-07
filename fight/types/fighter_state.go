package types

import "pb"

type FighterStateMachine struct {
	Owner    *Fighter
	StateMap []IFighterState
	PreState IFighterState
	CurState IFighterState
}

func (fsm *FighterStateMachine) Init(owner *Fighter) {
	fsm.Owner = owner
	fsm.StateMap = make([]IFighterState, len(pb.EFighterState_name))
}

func (fsm *FighterStateMachine) AddState(eType pb.EFighterState, state IFighterState) {
	fsm.StateMap[eType] = state
}

func (fsm *FighterStateMachine) SetCurrentState(eType pb.EFighterState) {
	fsm.PreState = fsm.StateMap[eType]
	fsm.CurState = fsm.StateMap[eType]
}

func (fsm *FighterStateMachine) ChangeState(eType pb.EFighterState) {
	state := fsm.StateMap[eType]
	if state != nil {
		if fsm.CurState != nil {
			fsm.CurState.Exit(fsm.Owner)
			fsm.PreState = fsm.CurState
		}
		fsm.CurState = state
		fsm.CurState.Enter(fsm.Owner)
	}
}

func (fsm *FighterStateMachine) Execute() {
	fsm.CurState.Execute(fsm.Owner)
}

type IFighterState interface {
	Enter(*Fighter)
	Execute(*Fighter)
	Exit(*Fighter)
	OnEvent(*Fighter, *pb.MsgFighterStateEvent)
}

//////////////////////////////////////////////////////////////////////////////

type FsDefault struct {
}

func (fs *FsDefault) Enter(*Fighter) {

}
func (fs *FsDefault) Execute(*Fighter) {

}
func (fs *FsDefault) Exit(*Fighter) {

}
func (fs *FsDefault) OnEvent(owner *Fighter, event *pb.MsgFighterStateEvent) {
	if event.EType == pb.EFighterStateEvent_FSE_TargetState {
		owner.FSM.ChangeState(event.Target)
	}
}

//////////////////////////////////////////////////////////////////////////////
