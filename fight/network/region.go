package network

import (
	"fight/types"
	"pb"
)

//战斗场景的接口
type IRegion interface {
	OnCreate(rgf *RgnForm)                                    //创建时
	OnRoleConnect(r *types.Role)                              //角色连接进来
	OnRoleDisConnect(r *types.Role)                           //角色连接断开
	OnRoleEnter(r *types.Role)                                //角色进入
	OnRoleLeave(r *types.Role, winner uint32)                 //角色退出
	Run()                                                     //运行
	SecLoop()                                                 //每秒运行
	GetFighter(fer uint32) *types.Fighter                     //获取玩家
	OnAreaScriptFinish(uint32, *pb.MsgBattleAreaScriptFinish) //区域脚本结束
}

type Region struct {
	rgf       *RgnForm
	IsStart   bool             //战斗是否已经开始
	IsFinish  bool             //战斗是否已经结束
	MFighters []*types.Fighter //怪物阵营
	gk        *RegionGK
}

func (rg *Region) OnCreate(rfg *RgnForm) {
	rg.MFighters = make([]*types.Fighter, 0)
	rg.gk = &RegionGK{}
	rg.gk.Init(rg)
	rg.rgf = rfg
}
func (rg *Region) OnRoleConnect(r *types.Role) {

}
func (rg *Region) OnRoleDisConnect(r *types.Role) {
	rg.gk.OnRoleDisConnect(r.Guid)
}
func (rg *Region) OnRoleEnter(r *types.Role) {
	r.HadEnter = true
	r.Ses.Send(pb.MsgIDS2C_BattleEnterAck, &pb.MsgBattleEnterAck{RetCode: 1})
	//组装开始消息
	startMsg := &pb.MsgBattleStart{}
	startMsg.TriggerIDs = rg.gk.GetHadTriIDs() //组装已触发区域块
	startMsg.AreaScripts = rg.gk.AreaScripts   //组装已触发的区域脚本
	if rg.IsStart {
		r.Ses.Send(pb.MsgIDS2C_BattleStart, startMsg) //中途进入
		return
	}
	for _, v := range rg.rgf.Roles {
		if !v.HadEnter {
			return
		}
	}
	rg.IsStart = true
	rg.rgf.SendToAll(pb.MsgIDS2C_BattleStart, startMsg) //全部进入
}
func (rg *Region) OnRoleLeave(r *types.Role, winner uint32) {
	r.HadLeave = true
	r.Ses.Send(pb.MsgIDS2C_BattleLeaveAck, &pb.MsgBattleEnterAck{RetCode: 1})
	if rg.IsFinish {
		return
	}
	for _, v := range rg.rgf.Roles {
		if !v.HadLeave {
			return
		}
	}
	rg.IsFinish = true
	finMsg := &pb.MsgBattleFinishData{BattleGuid: rg.rgf.Id, Winner: winner}
	SendToFc(pb.MsgIDS2S_Ft2FcBattleFinish, finMsg)
	rg.rgf.Close()
}
func (rg *Region) Run() {
	if !rg.IsStart {
		return
	}
	if rg.IsFinish {
		return
	}
	//活跃怪
	Fits := make([]*types.Fighter, 0)
	FitMaps := make(map[uint32]*types.Fighter)
	for _, v := range rg.MFighters {
		if v.RealAttr.Hp > 0 {
			Fits = append(Fits, v)
			FitMaps[v.RealAttr.Guid] = v
		}
	}
	rg.MFighters = Fits
	//区域触发
	rg.gk.Run(FitMaps)
}
func (rg *Region) SecLoop() {

}
func (rg *Region) GetFighter(fer uint32) *types.Fighter {
	if fer < 1000 {
		for _, v := range rg.rgf.Roles {
			for _, v2 := range v.Fights {
				if v2.RealAttr.Guid == fer {
					return v2
				}
			}
		}
	} else {
		for _, v := range rg.MFighters {
			if v.RealAttr.Guid == fer {
				return v
			}
		}
	}
	return nil
}
func (rg *Region) OnAreaScriptFinish(roleGuid uint32, finish *pb.MsgBattleAreaScriptFinish) {
	rg.gk.OnAreaScriptFinish(roleGuid, finish)
}
