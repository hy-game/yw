package enet

import (
	"com/util"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func (s *Connection) OnRecv(data []byte) {
	defer util.PanicPrintStack(s, data)

	p, err := s.Pkg.NewReader(data, s.Ses)
	if err != nil {
		log.Warnf("read packet session:%s, err:%v", s.Ses.String(), err)
		s.Close()
		return
	}

	msgID := p.GetMsgID()
	seqNum := p.GetSeqNum() //only c2s

	if seqNum != 0 && seqNum != s.PkgCnt {
		log.Errorf("%s sequeue num err: %d should be %d", s.Ses.String(), seqNum, s.PkgCnt)
		s.Close()
		return
	}

	if !s.Handler.Handle(msgID, p.GetData(), s.Ses) {
		s.Ses.OnRecvInvalidMsg(p)
	}
}

//send 发送数据给客户端,非线程安全
func (s *Connection) Send(msgID uint16, data []byte) bool {
	p := s.Pkg.NewWriter(msgID, data, s.Ses)

	select {
	case s.out <- p:
		return true
	default:
		log.WithFields(log.Fields{"id": s.Ses.String(), "ip": s.Ip}).Warn("send queue full")
		return false
	}
}

func (s *Connection) SendPkg(writer IPkgWriter) bool {
	select {
	case s.out <- writer:
		return true
	default:
		log.WithFields(log.Fields{"id": s.Ses.String(), "ip": s.Ip}).Warn("send queue full")
		return false
	}
}

//SendPB 发送proto数据给客户端,非线程安全
func (s *Connection) SendPB(msgID uint16, msgData proto.Message) bool {
	var b []byte
	if msgData != nil {
		var err error
		b, err = proto.Marshal(msgData)
		if err != nil {
			log.Warnf("send pb, marshal error:%v", err)
			return false
		}
	}

	if s.Send(msgID, b) {
		s.Ses.TraceSend(msgID, msgData)
		return true
	} else {
		return false
	}
}
