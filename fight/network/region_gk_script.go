package network

import (
	"fight/configs"
	"fight/types"
	"pb"
)

type RegionScriptLogic struct {
	sendTime    int64
	isFinish    bool
	FinishRoles map[uint32]bool
}

func (sl *RegionScriptLogic) AddWait(roleMap map[uint32]*types.Role) {
	sl.isFinish = false
	sl.FinishRoles = make(map[uint32]bool)
	for k, v := range roleMap {
		if v.HadEnter && v.Ses != nil {
			sl.FinishRoles[k] = false
		}
	}
}

func (sl *RegionScriptLogic) AddRole(roleGuid uint32) {
	if !sl.isFinish {
		sl.FinishRoles[roleGuid] = true
		for _, v := range sl.FinishRoles {
			if v == false {
				return
			}
		}
		sl.isFinish = true //所有用户结束
	}
}

func (sl *RegionScriptLogic) RunAreaScript(TriggerID, AreaID uint32, sType string, area *configs.RegionArea, rg *Region) {
	rgf := rg.rgf
	if sl.isFinish {
		return
	}
	if sl.sendTime > 0 {
		if rgf.nowTime > sl.sendTime+5000 {
			sl.isFinish = true //等待超时
		}
	} else {
		sl.sendTime = rgf.nowTime
		sl.isFinish = true
		for _, vScript := range area.ScriptList {
			if vScript.Type == sType {
				sl.AddWait(rgf.Roles)
				scriptMsg := &pb.MsgBattleAreaScriptStart{TriggerID: TriggerID, AreaID: AreaID, Type: sType}
				for _, vScriptOne := range vScript.Script {
					oneScript := &pb.MsgBattleScript{ID: vScriptOne.Id, Path: vScriptOne.Path, Param: vScriptOne.Param}
					scriptMsg.Scripts = append(scriptMsg.Scripts, oneScript)
				}
				rgf.SendToAll(pb.MsgIDS2C_BattleAreaScriptStart, scriptMsg)
				rg.gk.AreaScripts = append(rg.gk.AreaScripts, scriptMsg) //保存已触发的区域块
				break
			}
		}
	}
}
