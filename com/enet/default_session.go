package enet

import (
	"com/evt"
	"com/log"
	"github.com/golang/protobuf/proto"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	sesIDGen int32
)

//默认处理
type Session struct {
	Conn *Connection
	Id   uint32
}

func (s *Session) String() string {
	return "ses id " + strconv.Itoa(int(s.Id))
}

func (s *Session) SendPB(msgID uint16, msg proto.Message) bool {
	return s.Conn.SendPB(msgID, msg)
}

func (s *Session) OnConnect() {
	s.Id = uint32(atomic.AddInt32(&sesIDGen, 1))
	addSession(s.Id, s)
	log.Infof("%s %s connect", s.String(), s.Conn.Ip.String())
}
func (s *Session) Close() {
	s.Conn.Close()
}
func (s *Session) OnClosed() {
	removeSession(s.Id)
	log.Infof("%s disconnect", s.String())
}

func (s *Session) OnRecvInvalidMsg(reader IPkgReader) {
	log.Infof("%s recv invalid msg [%d]", s.String(), reader.GetMsgID())
}

func (s *Session) OnEvent(event evt.Evt) {

}

func (s *Session) TraceRecv(msgId uint16, msg proto.Message, duration time.Duration) {
	msgStr := ""
	if msg != nil {
		msgStr = msg.String()
	}
	log.Tracef("%s recv [%d] usetime:%v data:%s",
		s.String(), msgId, duration, msgStr)
}

func (s *Session) TraceSend(msgId uint16, msg proto.Message) {
	msgStr := ""
	if msg != nil {
		msgStr = msg.String()
	}
	log.Tracef("%s send [%d] data:%s",
		s.String(), msgId, msgStr)
}

//sess mgr
var sess sync.Map

func addSession(id uint32, s *Session) {
	sess.Store(id, s)
}

func removeSession(id uint32) {
	sess.Delete(id)
}

func GetSession(id uint32) *Session {
	s, ok := sess.Load(id)
	if ok {
		return s.(*Session)
	}
	return nil
}
