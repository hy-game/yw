package types

import (
	"github.com/golang/protobuf/proto"
	"pb"
)

var offLineEvt = make([]func(*Role), 0)

//AddOfflineEvt 添加角色的下线处理事件
func AddOfflineEvt(f func(*Role)) {
	offLineEvt = append(offLineEvt, f)
}

//Role	角色数据
type Role struct {
	Acc    string
	Guid   uint32
	Ses    *Session
	Comps  map[TypeComp]IComp
	Data   *pb.MsgPlayerData
	Battle *pb.MsgBattleStartData
}

//GetComp	获取组件
func (r *Role) GetComp(t TypeComp) IComp {
	return r.Comps[t]
}

func (r *Role) onDisconnect() {
	for _, v := range offLineEvt {
		v(r)
	}
}

//Send	发送数据
func (r *Role) Send(msgID pb.MsgIDS2C, msgData proto.Message) {
	if r.Ses != nil {
		r.Ses.Send(msgID, msgData)
	}
}
