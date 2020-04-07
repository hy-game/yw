package enet

import (
	"com/evt"
	"github.com/golang/protobuf/proto"
	"time"
)

type INet interface {
	OnCreateSession(s *Connection)
	MsgRegistry(msgID uint16, cf func() proto.Message, df func(msg proto.Message, s ISession))
	GetCfg() *Config
}

type IPkgReader interface {
	GetMsgID() uint16
	GetSeqNum() uint32
	GetData() []byte
}

type IPkgWriter interface {
	Write(retCache []byte) int
}

type IPkg interface {
	NewReader([]byte, ISession) (IPkgReader, error)
	NewWriter(msgID uint16, data []byte, s ISession) IPkgWriter
}

type ISession interface {
	String() string

	OnConnect()
	OnClosed()

	OnRecvInvalidMsg(reader IPkgReader)
	OnEvent(event evt.Evt)

	TraceRecv(msgId uint16, msg proto.Message, duration time.Duration)
	TraceSend(msgId uint16, msg proto.Message)
}
