package gnet

import (
	"com/log"
	"com/util"
	"github.com/golang/protobuf/proto"
	"pb"
)

//send 发送数据给客户端,非线程安全
func (s *Session) Send(msgID uint16, data []byte) bool {
	select {
	case s.out <- &PkgWriter{
		msgId: msgID,
		data:  data,
	}:
		return true
	default:
		log.Warnf("%s send queue full", s.String())
		return false
	}
}

//SendPB 发送proto数据给客户端,非线程安全
func (s *Session) SendPB(msgID uint16, msgData proto.Message) bool {
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
		msgStr := ""
		if msgData != nil {
			msgStr = msgData.String()
		}
		if msgID != uint16(pb.MsgIDS2C_S2CHeartBeat) {
			log.Infof("%s send [%d]%s data:%s",
				s.String(), msgID, pb.MsgIDS2C_name[int32(msgID)], msgStr)
		}
		return true
	} else {
		return false
	}
}

//send
func (s *Session) sendLoop(cfg *Config) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			util.PrintStack()
		}
		log.Debug("send loop stop")
	}()

	for {
		select {
		case p := <-s.out:
			s.rawSend(p)
		case <-s.Ctrl:
			return
		}
	}
}

func (s *Session) rawSend(p *PkgWriter) {
	if len(p.data) > sendPacketLimit {
		sendPacketLimit = len(p.data)
		s.outCache = make([]byte, sendPacketLimit+6)
	}
	len := p.Write(s.outCache)
	_, err := s.Conn.Write(s.outCache[:len])
	if err != nil {
		log.Warnf("send data error, err:%v", err)
	}
}
