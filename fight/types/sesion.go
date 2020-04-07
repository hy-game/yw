package types

import (
	"com/log"
	"github.com/golang/protobuf/proto"
	"pb"
	"strconv"
)

//Session	网络会话
type Session struct {
	guid   uint32
	Die    chan struct{}
	stream pb.SrvService_SrvSrvServer
}

//NewSession	新建一个网络会话
func NewSession(stream pb.SrvService_SrvSrvServer) *Session {
	s := &Session{
		Die:    make(chan struct{}),
		stream: stream,
	}
	return s
}

//Desc 描述
func (s *Session) Desc() string {
	return strconv.Itoa(int(s.guid))
}

//Send 发送消息给客户端
func (s *Session) Send(msgId pb.MsgIDS2C, msg proto.Message) {
	if s.stream == nil {
		log.Warnf("%s session closed", s.Desc())
		return
	}

	var b []byte
	var err error

	if msg != nil {
		b, err = proto.Marshal(msg)
		if err != nil {
			log.Warnf("%s marshal err %v when send:%s", s.Desc(), err, msg.String())
			return
		}
	}
	gameMsg := &pb.SrvMsg{
		Msg: b,
		ID:  uint32(msgId),
	}

	if err := s.stream.Send(gameMsg); err != nil {
		log.Warnf("%s send msg err:%v", s.Desc(), err)
		return
	}

	var msgStr string
	if msg != nil {
		msgStr = msg.String()
	}
	log.Tracef("send [%d]%s to %s data:%s", msgId, pb.MsgIDS2C_name[int32(msgId)], s.Desc(), msgStr)
}

//Close	主动关闭网络会话
func (s *Session) Close() {
	s.Send(pb.MsgIDS2C_Ft2GtKickRole, nil)
}
