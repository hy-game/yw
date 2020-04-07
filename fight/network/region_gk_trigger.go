package network

import (
	"com/log"
	"fight/configs"
	"fight/types"
	"pb"
)

type RegionTriggerLogic struct {
	TriggerID     uint32
	AreaID        uint32 //
	Area          *configs.RegionArea
	Dif           *configs.RegionDifficulty
	Wave          uint32
	List          []uint32
	MonsterFinish bool
	Enter         RegionScriptLogic
	Exit          RegionScriptLogic
}

func (logic *RegionTriggerLogic) OnRoleDisConnect(roleGuid uint32) {
	logic.Enter.AddRole(roleGuid)
	logic.Exit.AddRole(roleGuid)
}

func (logic *RegionTriggerLogic) OnAreaScriptFinish(roleGuid uint32, finish *pb.MsgBattleAreaScriptFinish) {
	if finish.Type == "Enter" {
		logic.Enter.AddRole(roleGuid)
	} else if finish.Type == "Exit" {
		logic.Exit.AddRole(roleGuid)
	}
}

func (logic *RegionTriggerLogic) CreateMonster(rg *Region, FitMaps map[uint32]*types.Fighter) {
	if !logic.Enter.isFinish {
		return
	}
	if logic.MonsterFinish {
		return
	}
	allDie := true
	for _, v2 := range logic.List {
		if FitMaps[v2] != nil {
			allDie = false
			break
		}
	}
	if !allDie {
		return
	}
	logic.MonsterFinish = true
	log.Warn("monster killed", logic.List)
	for _, v3 := range logic.Dif.MonsterWave {
		if v3.Id > logic.Wave {
			//再来一波怪v3
			logic.MonsterFinish = false
			logic.Wave = v3.Id
			logic.List = nil
			createMsg := &pb.MsgBattleMonsterCreate{AreaID: logic.AreaID, WaveID: logic.Wave}
			for k4, v4 := range v3.Monster {
				//创建v4怪物
				initAttr := &pb.MsgFighter{
					ID:      0,
					Guid:    logic.AreaID*10000000 + logic.Dif.Level*1000000 + logic.Wave*10000 + (uint32(k4 + 1)),
					Atk:     100,
					Def:     100,
					Hp:      100,
					MaxHp:   100,
					Speed:   800,
					PosX:    v4.ServerCastPos[0],
					PosY:    v4.ServerCastPos[1],
					PosZ:    v4.ServerCastPos[2],
					Yaw:     v4.ServerCastYaw,
					Visible: true,
				}
				newMonster := CreateFighter(initAttr, true)
				newMonsterMsg := &pb.MsgBattleMonster{Fighter: initAttr}
				logic.List = append(logic.List, newMonster.RealAttr.Guid)                             //加入区域怪列表
				rg.MFighters = append(rg.MFighters, newMonster)                                       //加入场景怪列表
				createMsg.List = append(createMsg.List, newMonsterMsg)                                //加入发送列表
				newMonsterMsg.AI = &pb.MsgBattleMonsterAI{ExcludeMonsterGID: v4.AI.ExcludeMonsterGID} //附加怪物AI
				//附加怪物脚本让客户端自己运行，服务器不等待结果
				for _, v5 := range v4.ScriptList {
					oneList := &pb.MsgBattleScriptList{Type: v5.Type}
					newMonsterMsg.Scripts = append(newMonsterMsg.Scripts, oneList)
					for _, v6 := range v5.Script {
						oneScript := &pb.MsgBattleScript{ID: v6.Id, Path: v6.Path, Param: v6.Param}
						oneList.Scripts = append(oneList.Scripts, oneScript)
					}
				}
			}
			log.Warn("monster create", logic.List)
			rg.rgf.SendToAll(pb.MsgIDS2C_BattleMonsterCreate, createMsg)
			break
		}
	}
	if logic.MonsterFinish {
		logic.Exit.RunAreaScript(logic.TriggerID, logic.AreaID, "Exit", logic.Area, rg)
	}
}
