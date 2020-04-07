package gnet

import (
	"com/log"
	"encoding/binary"
	"io"
	"math"
	"pb"
	"time"

	"github.com/golang/protobuf/proto"
)

func RegistryCliMsg(msgID pb.MsgIDC2S, cf func() proto.Message, df func(msg proto.Message, s *Session)) {
	cliMsgHandler.Register(uint16(msgID), cf, df)
}

func RegistryGameMsg(msgID pb.MsgIDS2C, cf func() proto.Message, df func(msg proto.Message, s *Session)) {
	gmMsgHandler.Register(uint16(msgID), cf, df)
}

//recv
func (s *Session) recvLoop(cfg *Config) {
	log.Debug("recv loop start")
	header := make([]byte, packetHeadLen)

	for {
		if cfg.ReadDeadline > 0 {
			s.Conn.SetReadDeadline(time.Now().Add(time.Second * cfg.ReadDeadline))
		}
		n, err := io.ReadFull(s.Conn, header)
		if err != nil {
			log.Warnf("%s read header failed, err:%v, size:%d", s.String(), err, n)
			return
		}
		//header是数据的大小，不包括header本身
		size := binary.BigEndian.Uint32(header)
		if size > cfg.RecvPkgLenLimit {
			log.Warnf("%s read data len out of limit, size%v", s.String(), size)
			return
		}
		data := make([]byte, size)
		n, err = io.ReadFull(s.Conn, data)
		if err != nil {
			log.Warnf("%s read data failed:, err:%v, size:%d", s.String(), err, n)
			return
		}

		select {
		case s.In <- data:
		case <-s.Ctrl:
			log.Warnf("%s connection close by logic", s.String())
			return
		}
	}
}

func (s *Session) onRecvCliMsg(data []byte) {
	p, err := NewReader(data)
	if err != nil {
		log.Warnf("%s read packet err:%v", s.String(), err)
		s.Close()
		return
	}

	msgID := p.GetMsgID()
	seqNum := p.GetSeqNum() //only c2s

	if seqNum != 0 && seqNum != s.PkgCnt {
		log.Errorf("%s sequeue num err: %d should be %d", s.String(), seqNum, s.PkgCnt)
		s.Close()
		return
	}

	if msgID > uint16(pb.MsgIDC2S_C2SGameMax) {
		s.forwardToFight(msgID, p.GetData())
	} else if msgID > uint16(pb.MsgIDC2S_C2SGateMax) {
		s.forwardToGame(msgID, p.GetData())
	} else {
		if !cliMsgHandler.Handle(msgID, p.GetData(), s) {
			log.Infof("%s recv invalid msg from cli [%d]", s.String(), msgID)
		}
	}
}

var gmMsgHandler = NewRoute(math.MaxUint16)
var cliMsgHandler = NewRoute(math.MaxUint16)

func (s *Session) onRecvGameMsg(msg pb.SrvMsg) {
	if msg.ID > uint32(pb.MsgIDS2C_Gm2GtMsgIdMax) {
		s.Send(uint16(msg.ID), msg.Msg)
		log.Debugf("%s forward to cli %d", s.String(), msg.ID)
	} else {
		if !gmMsgHandler.Handle(uint16(msg.ID), msg.Msg, s) {
			log.Warnf("%s recv invalid msg from game [%d]", s.String(), msg.ID)
		}
	}
}
