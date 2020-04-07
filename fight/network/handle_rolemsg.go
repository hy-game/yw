package network

import (
	"errors"
	"fight/types"
	"github.com/golang/protobuf/proto"
	"math"
	"pb"
)

var (
	errAPINotFind = errors.New("api not defined")
	errMsgParse   = errors.New("parser msg error")
)

var roleMsgRoute = newRoute(math.MaxUint16)

//RegisterCliMsgHandle	注册客户端消息处理函数
func RegisterCliMsgHandle(msgID pb.MsgIDC2S, cf func() proto.Message, df func(msg proto.Message, s *types.Role, rgn *RgnForm)) {
	roleMsgRoute.register(uint16(msgID), cf, df)
}

//HandleRoleMsg 客户端消息处理
func HandleRoleMsg(msgID uint32, msg []byte, role *types.Role, rgn *RgnForm) error {
	return roleMsgRoute.handle(uint16(msgID), msg, role, rgn)
}

type msgHandler struct {
	createFunc func() proto.Message
	handleFunc func(msg proto.Message, r *types.Role, rgn *RgnForm)
}

//route 消息处理器
type route struct {
	handlers []*msgHandler
}

//newRoute createRoute
func newRoute(size int) *route {
	r := &route{}
	r.init(size)
	return r
}

//register 注册消息
func (r *route) register(msgID uint16, cf func() proto.Message, df func(msg proto.Message, s *types.Role, rgn *RgnForm)) {
	n := &msgHandler{
		createFunc: cf,
		handleFunc: df,
	}
	r.handlers[msgID] = n
}

//handle 处理消息
func (r *route) handle(id uint16, data []byte, s *types.Role, rgn *RgnForm) error {
	//begin := time.Now()
	node, err := r.getHandler(id)
	if err != nil {
		return errAPINotFind
	}

	msg, err := r.parseMsg(node, data)
	if err != nil {
		return errMsgParse
	}

	node.handleFunc(msg, s, rgn)

	//if IsTraceCliMsg() {
	//	costTime := time.Now().Sub(begin)
	//	if id != uint16(pb.MsgIDC2S_C2SHeartBeat) { //心跳
	//		var msgStr string
	//		if msg != nil {
	//			msgStr = msg.String()
	//		}
	//		log.WithFields(log.Fields{
	//			"api":     id,
	//			"apiName": pb.MsgIDC2S_name[int32(id)],
	//			"cli":     s.Desc(),
	//			"data":    msgStr,
	//			"cost":    costTime,
	//		}).Debug("Done")
	//	}
	//}
	return nil
}

func (r *route) init(size int) {
	r.handlers = make([]*msgHandler, size)
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
