package network

import (
	"fight/configs"
	"fight/types"
	"pb"
)

type RegionGK struct {
	rg          *Region
	gk          *configs.RegionGK
	AllTriLogic map[uint32]*RegionTriggerLogic
	AreaScripts []*pb.MsgBattleAreaScriptStart
}

func (gk *RegionGK) Init(rg *Region) {
	gk.rg = rg
	gk.gk = configs.GetRegionLogic(rg.rgf.StartData.BattleID)
	gk.AllTriLogic = make(map[uint32]*RegionTriggerLogic)
}

func (gk *RegionGK) GetHadTriIDs() []uint32 {
	tIDs := make([]uint32, 0)
	for k := range gk.AllTriLogic {
		tIDs = append(tIDs, k)
	}
	return tIDs
}
func (gk *RegionGK) OnRoleDisConnect(roleGuid uint32) {
	for _, v := range gk.AllTriLogic {
		v.OnRoleDisConnect(roleGuid)
	}
}
func (gk *RegionGK) OnAreaScriptFinish(roleGuid uint32, finish *pb.MsgBattleAreaScriptFinish) {
	for _, v := range gk.AllTriLogic {
		if v.TriggerID == finish.TriggerID && v.AreaID == finish.AreaID {
			v.OnAreaScriptFinish(roleGuid, finish)
			break
		}
	}
}
func (gk *RegionGK) Run(FitMaps map[uint32]*types.Fighter) {
	gk.RunTrigger()
	gk.CreateMonster(FitMaps)
	isHadMonster := false
	if len(FitMaps) > 0 {
		isHadMonster = true
	}
	gk.CheckFinish(isHadMonster)
}

func (gk *RegionGK) RunTrigger() {
	tgs := gk.gk.TriggerList.Trigger
	roles := gk.rg.rgf.Roles
	for _, v2 := range tgs {
		_, find := gk.AllTriLogic[v2.Id]
		if find {
			continue
		}
		for _, v := range roles {
			if v.VisFight == nil {
				continue
			}
			if v2.IsIn(v.VisFight.RealAttr.PosX, v.VisFight.RealAttr.PosZ) {
				rgf := gk.rg.rgf
				area, dif := gk.gk.GetRegionDifficulty(v2.Id, rgf.StartData.BattleDifficulty)
				logic := &RegionTriggerLogic{TriggerID: v2.Id, AreaID: area.Id, Area: area, Dif: dif}
				gk.AllTriLogic[v2.Id] = logic
				triMsg := &pb.MsgBattleTriggerEnter{TriggerID: v2.Id}
				rgf.SendToAll(pb.MsgIDS2C_BattleTriggerEnter, triMsg)
				logic.Enter.RunAreaScript(logic.TriggerID, logic.AreaID, "Enter", logic.Area, gk.rg)
				break
			}
		}
	}
}

func (gk *RegionGK) CreateMonster(FitMaps map[uint32]*types.Fighter) {
	for _, v := range gk.AllTriLogic {
		v.CreateMonster(gk.rg, FitMaps)
	}
}

func (gk *RegionGK) CheckFinish(isHadMonster bool) {
	//战斗结束判断
	fts := gk.gk.FinishTypeList.FinishType
	allMonsterTris := gk.gk.FinishTypeList.AllMonsterTris
	difficulty := gk.rg.rgf.StartData.BattleDifficulty
	for _, v := range fts {
		if v.Id == 1 { //检查是否所有怪物已死亡
			if isHadMonster {
				continue //场景中还有活跃怪
			}
			isAllFinish := true
			for _, v2 := range allMonsterTris[difficulty] {
				logic, find := gk.AllTriLogic[v2]
				if !find || !logic.Exit.isFinish {
					isAllFinish = false
					break
				}
			}
			if !isAllFinish {
				continue //有怪的区域还没触发
			}
			gk.rg.IsFinish = true
			finMsg := &pb.MsgBattleFinishData{BattleGuid: gk.rg.rgf.Id}
			if v.Win {
				for _, v2 := range gk.rg.rgf.Roles {
					finMsg.Winner = v2.Guid
					break
				}
			}
			SendToFc(pb.MsgIDS2S_Ft2FcBattleFinish, finMsg)
			gk.rg.rgf.Close()
			break
		}
	}
}
