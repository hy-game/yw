package mq

import (
	"com/log"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"pb"
	"time"
)

var (
	ErrAPINotFind = errors.New("api not register")
	ErrMsgParser  = errors.New("parser msg error")
)

type msgHandlerCt struct {
	createFunc func() proto.Message
	handleFunc func(msg proto.Message, serId uint16)
}

//route 消息处理器
type routeNsq struct {
	handlers []*msgHandlerCt
}

//newRoute createRoute
func newRouteCt(size int) *routeNsq {
	r := &routeNsq{
		handlers: make([]*msgHandlerCt, size),
	}
	return r
}

//register 注册消息
func (r *routeNsq) register(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	n := &msgHandlerCt{
		createFunc: cf,
		handleFunc: df,
	}
	r.handlers[msgID] = n
}

func (r *routeNsq) getHandler(id uint16) (n *msgHandlerCt, err error) {
	if int(id) >= len(r.handlers) {
		err = ErrAPINotFind
		return
	}

	n = r.handlers[id]
	if nil == n || nil == n.createFunc {
		err = ErrAPINotFind
		return
	}
	return
}

//handle 处理消息
func (r *routeNsq) handle(id uint16, data []byte, serId uint16) error {
	begin := time.Now()
	node, err := r.getHandler(id)
	if err != nil {
		return ErrAPINotFind
	}

	msg, err := r.parseMsg(node, data)
	if err != nil {
		return ErrMsgParser
	}

	node.handleFunc(msg, serId)

	costTime := time.Now().Sub(begin)
	var msgStr string
	if msg != nil {
		msgStr = msg.String()
	}
	if id > uint16(pb.MsgIDS2S_S2STraceEnd) {
		log.Tracef("handle [%d]%s cost:%v from:%d %s", id, pb.MsgIDS2S_name[int32(id)], costTime, serId, msgStr)
	}
	return nil
}

func (r *routeNsq) parseMsg(n *msgHandlerCt, data []byte) (msg proto.Message, err error) {
	msg = n.createFunc()
	if msg == nil { //允许只有消息id没内容
		return
	}
	err = proto.Unmarshal(data, msg)
	return
}
