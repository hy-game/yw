package network

import (
	"game/types"
	"math"
	"pb"
)

type evtHandler struct {
	handleFunc func(e types.Evt, s *types.Session, r *types.Role)
}

var evtHandle = newEvtRoute(math.MaxUint16)

//RegisterCliHandle 注册事件消息处理函数
func RegisterEvtHandle(msgID pb.MsgIDC2S, df func(e types.Evt, s *types.Session, r *types.Role)) {
	evtHandle.register(int(msgID), df)
}

//route 消息处理器
type evtRoute struct {
	handlers []*evtHandler
}

//newRoute createRoute
func newEvtRoute(size int) *evtRoute {
	r := &evtRoute{make([]*evtHandler, size)}
	return r
}

//register 注册消息
func (r *evtRoute) register(msgID int, df func(e types.Evt, s *types.Session, r *types.Role)) {
	n := &evtHandler{
		handleFunc: df,
	}
	r.handlers[msgID] = n
}

//handle 处理消息
func (r *evtRoute) handle(e types.Evt, s *types.Session) error {
	node, err := r.getHandler(e.Type)
	if err != nil {
		return errAPINotFind
	}

	node.handleFunc(e, s, s.Role)
	return nil
}

func (r *evtRoute) getHandler(id types.GameEvent) (n *evtHandler, err error) {
	if int(id) >= len(r.handlers) {
		err = errAPINotFind
		return
	}

	n = r.handlers[id]
	return
}
