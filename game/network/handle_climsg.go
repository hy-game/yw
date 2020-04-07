package network

import (
	"com/log"
	"game/types"
	"github.com/golang/protobuf/proto"
	"math"
	"pb"
	"time"
)

var cliMsgHandler = newRoute(math.MaxUint16)

//RegisterCliHandle 注册客户端发来的消息处理函数
func RegisterCliHandle(msgID pb.MsgIDC2S, cf func() proto.Message, df func(msg proto.Message, s *types.Session)) {
	cliMsgHandler.register(uint16(msgID), cf, df)
}

type msgHandler struct {
	createFunc func() proto.Message
	handleFunc func(msg proto.Message, s *types.Session)
}

//route 消息处理器
type route struct {
	handlers []*msgHandler
}

//newRoute createRoute
func newRoute(size int) *route {
	r := &route{
		make([]*msgHandler, size),
	}
	return r
}

//register 注册消息
func (r *route) register(msgID uint16, cf func() proto.Message, df func(msg proto.Message, s *types.Session)) {
	n := &msgHandler{
		createFunc: cf,
		handleFunc: df,
	}
	r.handlers[msgID] = n
}

//handle 处理消息
func (r *route) handle(id uint16, data []byte, s *types.Session) error {
	begin := time.Now()
	node, err := r.getHandler(id)
	if err != nil {
		return errAPINotFind
	}

	msg, err := r.parseMsg(node, data)
	if err != nil {
		return errMsgParse
	}

	node.handleFunc(msg, s)

	costTime := time.Now().Sub(begin)
	var msgStr string
	if msg != nil {
		msgStr = msg.String()
	}
	log.Tracef("handle [%d]%s cost:%v, msg:%s", id, pb.MsgIDC2S_name[int32(id)], costTime, msgStr)

	return nil
}

func (r *route) getHandler(id uint16) (n *msgHandler, err error) {
	if int(id) >= len(r.handlers) {
		err = errAPINotFind
		return
	}

	n = r.handlers[id]
	if nil == n || nil == n.createFunc {
		err = errAPINotFind
		return
	}
	return
}

func (r *route) parseMsg(n *msgHandler, data []byte) (msg proto.Message, err error) {
	msg = n.createFunc()
	if msg == nil { //允许只有消息id没内容
		return
	}
	err = proto.Unmarshal(data, msg)
	return
}
