/*
客户端消息处理
*/
package logic

import (
	"fight/network"
	"fight/types"
	"pb"

	"github.com/golang/protobuf/proto"
)

//注意：必须写注释，这些函数是角色所在场景协程调用的
func initRoleMsgHandle() {
	network.RegisterCliMsgHandle(pb.MsgIDC2S_BattleEnterReqToFs, func() proto.Message { return &pb.MsgBattleEnterReq{} }, onBattleEnterReq) //客户端正式连接fighter进入场景
	network.RegisterCliMsgHandle(pb.MsgIDC2S_BattleLeaveReqToFs, func() proto.Message { return &pb.MsgBattleLeaveReq{} }, onBattleLeaveReq) //客户端正式连接fighter进入场景
	network.RegisterCliMsgHandle(pb.MsgIDC2S_BattlePhySync, func() proto.Message { return &pb.MsgBattlePhySync{} }, onBattlePhySync)
	network.RegisterCliMsgHandle(pb.MsgIDC2S_BattleDebugDamage, func() proto.Message { return &pb.MsgBattleDebugDamage{} }, onBattleDebugDamage)
	network.RegisterCliMsgHandle(pb.MsgIDC2S_BattleAreaScriptFinish, func() proto.Message { return &pb.MsgBattleAreaScriptFinish{} }, onBattleAreaScriptFinish)
}

func onBattleAreaScriptFinish(msg proto.Message, s *types.Role, rgn *network.RgnForm) {
	msgCast := msg.(*pb.MsgBattleAreaScriptFinish)
	rgn.Battle.OnAreaScriptFinish(s.Guid, msgCast)
}

func onBattleDebugDamage(msg proto.Message, s *types.Role, rgn *network.RgnForm) {
	msgCast := msg.(*pb.MsgBattleDebugDamage)
	ft := rgn.Battle.GetFighter(msgCast.Defer)
	if ft != nil {
		if msgCast.Hurt > 0 {
			ft.RealAttr.Hp -= uint32(msgCast.Hurt)
		} else {
			ft.RealAttr.Hp += uint32(-msgCast.Hurt)
		}
		if ft.RealAttr.Hp < 0 {
			ft.RealAttr.Hp = 0
		}
	}
}

func onBattlePhySync(msg proto.Message, s *types.Role, rgn *network.RgnForm) {
	msgCast := msg.(*pb.MsgBattlePhySync)
	s.OnPhySync(msgCast)
}

func onBattleLeaveReq(msg proto.Message, r *types.Role, rgn *network.RgnForm) {
	msgCast := msg.(*pb.MsgBattleLeaveReq)
	rgn.Battle.OnRoleLeave(r, msgCast.Winner)
}

//加载完成，进入场景
func onBattleEnterReq(msgBase proto.Message, r *types.Role, rgn *network.RgnForm) {
	rgn.Battle.OnRoleEnter(r)
}
