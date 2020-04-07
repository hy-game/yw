package types

import (
	"com/log"
	"com/util"
	"github.com/golang/protobuf/proto"
	"pb"
)

//Session 网络会话
type Session struct {
	ID       uint32
	Stream   pb.SrvService_SrvSrvServer
	RecvChan chan *pb.SrvMsg
	SendChan chan *pb.SrvMsg

	Evt     chan Evt
	Die     chan struct{}
	IsClose bool
	Role    *Role
	desc 	string
}

var (
	sessID = uint32(0)
)

//NewSession	新建一个网络会话
func NewSession(stream pb.SrvService_SrvSrvServer) *Session {
	sessID++

	s := &Session{
		Stream: stream,
		ID:     sessID,
		Die:    make(chan struct{}),
		desc:   util.ToString(sessID),
	}
	addRoleSes(s.ID, s)

	return s
}

//PostEvt	投递一个事件
func (s *Session) postEvt(e Evt) {
	select {
	case s.Evt <- e:
	default:
		log.Warnf("post evt to %s fail evt chan is full", s.Desc())
	}
}

//Desc 描述
func (s *Session) Desc() string {
	return s.desc
}

//Send Send msg to cli
func (s *Session) Send(msgID pb.MsgIDS2C, msg proto.Message) {
	var b []byte
	var err error

	if msg != nil {
		b, err = proto.Marshal(msg)
		if err != nil {
			log.Warnf("%s marshal err %v when send:%s", s.Desc(), err, msg.String())
			return
		}
	}
	s.SendByte(uint32(msgID), b)

	var msgStr string
	if msg != nil {
		msgStr = msg.String()
	}
	log.Tracef("send [%d]%s to %s data:%s", msgID, pb.MsgIDS2C_name[int32(msgID)], s.Desc(), msgStr)
}

func (s *Session)SendByte(msgID uint32, data []byte){
	gameMsg := &pb.SrvMsg{
		Msg: data,
		ID:  msgID,
	}
	select {
	case s.SendChan <- gameMsg:
	default:
		log.Warnf("Send chan is full %s", s.Desc())
	}
}

//Close	关闭连接
func (s *Session) Close() {
	s.Send(pb.MsgIDS2C_Gm2GtKickRole, nil)
}

//OnClose	离线处理 非线程安全
func (s *Session) OnClose() {
	if s.Role != nil {
		s.Role.onDisconnect()
	}
	delRoleSes(s.ID)
	log.Infof("%s closed", s.Desc())
}
