package types

import (
	"github.com/golang/protobuf/proto"
	"pb"
)

//Role 角色数据
type Role struct {
	Guid     uint32
	Ses      *Session
	HadEnter bool //是否已经进入
	HadLeave bool //是否已经离开
	Fights   []*Fighter
	VisFight *Fighter
}

type Evt struct {
}

func (r *Role) OnEvent(e Evt) {

}

func (r *Role) OnPhySync(sync *pb.MsgBattlePhySync) {
	if r.VisFight != nil && r.VisFight.RealAttr.Hp > 0 {
		r.VisFight.RealAttr.PosX = sync.PosX
		r.VisFight.RealAttr.PosY = sync.PosY
		r.VisFight.RealAttr.PosZ = sync.PosZ
		r.VisFight.RealAttr.Yaw = sync.Yaw
	}
}

func NewRole(data *pb.MsgRoleInFight, fts []*Fighter, vft *Fighter) *Role {
	r := &Role{Guid: data.Guid}
	r.Fights = fts
	r.VisFight = vft
	return r
}

func (r *Role) Send(msgID pb.MsgIDS2C, msgData proto.Message) {
	if r.Ses != nil {
		r.Ses.Send(msgID, msgData)
	}
}
