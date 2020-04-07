package enet

import (
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var (
	errAPINotRegist = errors.New("api not regist")
)

type msgHandler struct {
	createFunc func() proto.Message
	handleFunc func(msg proto.Message, s ISession)
}

//Route 消息处理器
type Route struct {
	handlers []*msgHandler
}

//NewRoute createRoute
func NewRoute(size int) *Route {
	r := &Route{}
	r.init(size)
	return r
}

//Register 注册消息
func (r *Route) Register(msgID uint16, cf func() proto.Message, df func(msg proto.Message, s ISession)) {
	n := &msgHandler{
		createFunc: cf,
		handleFunc: df,
	}
	r.handlers[msgID] = n
}

//Handle 处理消息
func (r *Route) Handle(id uint16, data []byte, s ISession) bool {
	begin := time.Now()
	node, err := r.getHandler(id)
	if err != nil {
		return false
	}

	msg, err := r.parseMsg(node, data)
	if err != nil {
		log.Warnf("parser msg %d error:%v", id, err)
		return true
	}

	node.handleFunc(msg, s)

	cost := time.Now().Sub(begin)
	s.TraceRecv(id, msg, cost)

	return true
}

func (r *Route) init(size int) {
	r.handlers = make([]*msgHandler, size)
}

func (r *Route) getHandler(id uint16) (n *msgHandler, err error) {
	if int(id) >= len(r.handlers) {
		err = errAPINotRegist
		return
	}

	n = r.handlers[id]
	if nil == n || nil == n.createFunc {
		err = errAPINotRegist
		return
	}
	return
}

func (r *Route) parseMsg(n *msgHandler, data []byte) (msg proto.Message, err error) {
	msg = n.createFunc()
	if msg == nil { //允许只有消息id没内容
		return
	}
	err = proto.Unmarshal(data, msg)
	return
}
