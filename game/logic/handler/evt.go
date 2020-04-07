/*
事件处理
*/
package handler

import (
	"game/logic/comps/mail"
	"game/logic/role"
	"game/network"
	"game/types"
	"pb"
)

//注意：这里要写注释。未绑定前(LoadRoleData, ReConn)，r==nil,
func registEvtHanele() {
	network.RegisterEvtHandle(types.LoadRoleData, onLoadData)              //角色数据加载完成，登录成功
	network.RegisterEvtHandle(types.ReConn, onEvtReConn)                   //重连
	network.RegisterEvtHandle(types.KickOut, onKickOut)                    //剔掉角色
	network.RegisterEvtHandle(types.RoleReLogin, onRoleReLogin)            //角色异地登录处理
	network.RegisterEvtHandle(types.SecLoop, onSecLoop)                    //角色每秒更新
	network.RegisterEvtHandle(types.ForwardToRole, onEvtForwardMsg)         //转发center的消息
	network.RegisterEvtHandle(types.GMOrder, onGMOrder)                    //角色Gm命令
	network.RegisterEvtHandle(types.BattleCreateAck, onEvtBattleCreateAck) //创建战斗返回
	network.RegisterEvtHandle(types.BattleFinish, onEvnBattleFinish)       //战斗结束
	network.RegisterEvtHandle(types.RecvMail, onEvtRecvMail)       //战斗结束
}

func onLoadData(e types.Evt, s *types.Session, r *types.Role) {
	role.Online(types.GetRoleForConnEvt(&e), s)
}

func onEvtReConn(e types.Evt, s *types.Session, r *types.Role) {
	role.ReConnSuccess(types.GetRoleForConnEvt(&e), s)
}

func onEvtForwardMsg(e types.Evt, s *types.Session, r *types.Role){
	msg := e.Data.(*pb.MsgForwardToRole)
	if msg == nil {
		return
	}
	s.SendByte(msg.MsgID, msg.Data)
}

func onEvtBattleCreateAck(e types.Evt, s *types.Session, r *types.Role) {

}

func onEvnBattleFinish(e types.Evt, s *types.Session, r *types.Role) {
	role.FinishBattle(r, e.Data.(*pb.MsgBattleFinishData))
}

func onKickOut(e types.Evt, s *types.Session, r *types.Role) {
	s.Close()
}

func onSecLoop(e types.Evt, s *types.Session, r *types.Role) {
	if r == nil {
		return
	}
	role.SecLoop(r)
}

func onRoleReLogin(e types.Evt, s *types.Session, r *types.Role) {
	data := e.Data.(*pb.MsgReLogin)
	if data == nil {
		return
	}
	//todo 是否需要发个消息给客户端
	s.Close()

	s.Role = nil
	r.Ses = nil

	role.Mgr4Role.RoleReLogin(data.Acc, data.SesID)
}

func onGMOrder(e types.Evt, s *types.Session, r *types.Role) {
	msg, ok := e.Data.(*pb.MsgGMTask)
	if !ok {
		return
	}
	if r == nil {
		return
	}
	handleGMOrder(msg, s, r)
}

func onEvtRecvMail(e types.Evt, s *types.Session, r *types.Role) {
	msg, ok := e.Data.(*pb.MsgMail)
	if !ok {
		return
	}

	mail.Add(msg, r)
}
