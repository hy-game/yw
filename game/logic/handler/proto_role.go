/*
gate转发的客户端消息处理
*/
package handler

import (
	"game/logic/comps/fb"
	"game/logic/comps/hero"
	"game/logic/role"
	"game/network"
	"game/types"
	"pb"

	"github.com/golang/protobuf/proto"
)

//注意：必须写注释  这些函数在每个角色的协程调用
func initCliMsgHandle() {
	network.RegisterCliHandle(pb.MsgIDC2S_C2SLogin, func() proto.Message { return &pb.MsgLogin{} }, onLogin)                             //登录请求
	network.RegisterCliHandle(pb.MsgIDC2S_C2SReConn, func() proto.Message { return &pb.MsgReConn{} }, onReConn)                          //重连请求
	network.RegisterCliHandle(pb.MsgIDC2S_BattleCreateReq, func() proto.Message { return &pb.MsgBattleCreateReq{} }, onBattleCreateReq)  //请求进入战斗场景
	network.RegisterCliHandle(pb.MsgIDC2S_BattleEnterReqToGs, func() proto.Message { return &pb.MsgBattleEnterReq{} }, onBattleEnterReq) //请求进入战斗场景
	network.RegisterCliHandle(pb.MsgIDC2S_BattleLeaveReqToGs, func() proto.Message { return &pb.MsgBattleLeaveReq{} }, onBattleLeaveReq) //请求进入战斗场景
	network.RegisterCliHandle(pb.MsgIDC2S_C2SHero, func() proto.Message { return &pb.MsgHeroProto{} }, onHero)                           //请求进入战斗场景
}

func onBattleLeaveReq(msg proto.Message, s *types.Session) {
	msgCast := msg.(*pb.MsgBattleLeaveReq)
	if s.Role.Battle != nil {
		s.Role.Send(pb.MsgIDS2C_BattleLeaveAck, &pb.MsgBattleLeaveAck{RetCode: 1})
		role.FinishBattle(s.Role, &pb.MsgBattleFinishData{Winner: msgCast.Winner})
		s.Role.Battle = nil //战斗结束了
	} else {
		s.Role.Send(pb.MsgIDS2C_BattleLeaveAck, &pb.MsgBattleLeaveAck{RetCode: 101})
	}
}

func onBattleEnterReq(msg proto.Message, s *types.Session) {
	if s.Role.Battle != nil {
		s.Role.Send(pb.MsgIDS2C_BattleEnterAck, &pb.MsgBattleEnterAck{RetCode: 1})
		s.Role.Send(pb.MsgIDS2C_BattleStart, &pb.MsgBattleStart{})
	} else {
		s.Role.Send(pb.MsgIDS2C_BattleEnterAck, &pb.MsgBattleEnterAck{RetCode: 101})
	}
}

func onLogin(msgBase proto.Message, s *types.Session) {
	msg := msgBase.(*pb.MsgLogin)
	if msg == nil {
		return
	}

	role.Mgr4Role.ReqRoleData(msg, s.ID)
}

func onReConn(msgBase proto.Message, s *types.Session) {
	msg := msgBase.(*pb.MsgReConn)
	if msg == nil {
		return
	}

	role.Mgr4Role.ReqReConn(msg, s.ID)
}

func onBattleCreateReq(msgBase proto.Message, s *types.Session) {
	msgRecv := msgBase.(*pb.MsgBattleCreateReq)
	if msgRecv == nil {
		return
	}

	if s.Role.Battle != nil {
		//todo: 自动结算上次战斗
	}

	regionCfg := fb.CfgRegion(msgRecv.BattleID)
	if regionCfg == nil {
		return
	}

	if _, ok := role.CanEnterBattle(regionCfg, s.Role); ok{
		role.CreateBattle(regionCfg, s.Role, msgRecv.BattleID, msgRecv.Fighters)
	}else{

	}
}

func onHero(msgBase proto.Message, s *types.Session) {
	if s.Role == nil {
		return
	}

	hero.OnProto(msgBase, s.Role)
}
