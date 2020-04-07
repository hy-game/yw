package gnet

import (
	"com/log"
	"pb"

	"github.com/golang/protobuf/proto"
)

type EvtType int

const (
	EvtLoginSuccess EvtType = iota
)

type Evt struct {
	Type EvtType
	Data proto.Message
}

func (s *Session) PostEvt(evtType EvtType, param proto.Message) {
	evt := Evt{
		Type: evtType,
		Data: param,
	}
	s.evt <- evt
}

func (s *Session) onEvent(evt Evt) {
	switch evt.Type {
	case EvtLoginSuccess: //登录
		msg := evt.Data.(*pb.MsgLoginAck)
		s.login(msg)
	}
}

func (s *Session) login(msg *pb.MsgLoginAck) {
	if msg == nil || msg.Data == nil {
		log.Warnf("data is nil when login to game")
		return
	}
	s.startStreamGm(uint16(msg.GameID), msg.Data.Acc)
	data, err := proto.Marshal(msg.Data)
	if err != nil {
		log.Warnf("marshal err:%v when login to game", err)
		return
	}
	s.forwardToGame(uint16(pb.MsgIDC2S_C2SLogin), data)
}
