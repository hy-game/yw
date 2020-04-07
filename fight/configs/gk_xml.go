package configs

import (
	"com/cfgs"
	"com/log"
	"encoding/xml"
	"math"
	"pb"
	"strconv"
	"strings"
)

type RegionGK struct {
	XMLName        xml.Name             `xml:"root"`
	BornBos        RegionBornBos        `xml:"BornBos"`
	TriggerList    RegionTriggerList    `xml:"TriggerList"`
	AreaList       RegionAreaList       `xml:"AreaList"`
	FinishTypeList RegionFinishTypeList `xml:"FinishTypeList"`
}

func (gk *RegionGK) GetRegionDifficulty(triID uint32, dif uint32) (*RegionArea, *RegionDifficulty) {
	//note:只考虑一个触发器只能触发一个区域的情况
	for _, v := range gk.AreaList.Area {
		if v.TriggerId == triID {
			for _, v2 := range v.Difficulty {
				if v2.Level == dif {
					return v, v2
				}
			}
		}
	}
	return nil, nil
}

func castPos(value float64) int32 {
	return (int32)(math.Floor(value*100 + 0.5))
}

func (gk *RegionGK) InitServer() {
	//note:生成中间变量，方便计算
	bPos := strings.Split(gk.BornBos.Pos, ",")
	gk.BornBos.ServerPos[0], _ = strconv.ParseFloat(bPos[0], 10)
	gk.BornBos.ServerPos[1], _ = strconv.ParseFloat(bPos[1], 10)
	gk.BornBos.ServerPos[2], _ = strconv.ParseFloat(bPos[2], 10)
	gk.BornBos.ServerCastPos[0] = castPos(gk.BornBos.ServerPos[0])
	gk.BornBos.ServerCastPos[1] = castPos(gk.BornBos.ServerPos[1])
	gk.BornBos.ServerCastPos[2] = castPos(gk.BornBos.ServerPos[2])
	gk.BornBos.ServerCastYaw = castPos(gk.BornBos.Yaw)
	for _, v := range gk.TriggerList.Trigger {
		tPos := strings.Split(v.Pos, ",")
		v.ServerPos[0], _ = strconv.ParseFloat(tPos[0], 10)
		v.ServerPos[1], _ = strconv.ParseFloat(tPos[1], 10)
		v.ServerPos[2], _ = strconv.ParseFloat(tPos[2], 10)
		if v.Size != "" {
			tSize := strings.Split(v.Size, ",")
			v.ServerSize[0], _ = strconv.ParseFloat(tSize[0], 10)
			v.ServerSize[1], _ = strconv.ParseFloat(tSize[1], 10)
			v.ServerSize[2], _ = strconv.ParseFloat(tSize[2], 10)
		}
	}
	for _, v := range gk.AreaList.Area {
		for _, v2 := range v.Difficulty {
			for _, v3 := range v2.MonsterWave {
				for _, v4 := range v3.Monster {
					mPos := strings.Split(v4.Pos, ",")
					v4.ServerPos[0], _ = strconv.ParseFloat(mPos[0], 10)
					v4.ServerPos[1], _ = strconv.ParseFloat(mPos[1], 10)
					v4.ServerPos[2], _ = strconv.ParseFloat(mPos[2], 10)
					v4.ServerCastPos[0] = castPos(v4.ServerPos[0])
					v4.ServerCastPos[1] = castPos(v4.ServerPos[1])
					v4.ServerCastPos[2] = castPos(v4.ServerPos[2])
					v4.ServerCastYaw = castPos(v4.Yaw)
				}
			}
		}
	}
	//生成玩家结束事件信息
	gk.FinishTypeList.AllMonsterTris = make(map[uint32][]uint32)
	for _, v := range gk.AreaList.Area {
		for _, v2 := range v.Difficulty {
			for _, v3 := range v2.MonsterWave {
				if len(v3.Monster) > 0 {
					gk.FinishTypeList.AllMonsterTris[v2.Level] = append(gk.FinishTypeList.AllMonsterTris[v2.Level], v.TriggerId)
					break
				}
			}
		}
	}
}

//出生点
type RegionBornBos struct {
	Pos           string  `xml:"Pos,attr"`
	Yaw           float64 `xml:"Yaw,attr"`
	ServerPos     [3]float64
	ServerCastPos [3]int32
	ServerCastYaw int32
}

//触发区域
type RegionTrigger struct {
	Id         uint32  `xml:"Id,attr"`
	StartState bool    `xml:"StartState,attr"`
	Type       string  `xml:"Type,attr"` //Box,Sphere
	Pos        string  `xml:"Pos,attr"`
	Size       string  `xml:"Size,attr"`
	Yaw        float64 `xml:"Yaw,attr"`
	Radius     float64 `xml:"Radius,attr"`
	ServerPos  [3]float64
	ServerSize [3]float64
}

func (tri *RegionTrigger) IsIn(rolePosX, rolePosZ int32) bool {
	if tri.Type == "Box" {
		newRolePos := [3]float64{(float64)(rolePosX)/100 - tri.ServerPos[0], 0, (float64)(rolePosZ)/100 - tri.ServerPos[2]}
		newRolePos[0] = newRolePos[0]*math.Cos(tri.Yaw-360) + newRolePos[2]*math.Sin(tri.Yaw-360)
		newRolePos[2] = newRolePos[0]*math.Sin(tri.Yaw-360) + newRolePos[2]*math.Cos(tri.Yaw-360)
		if newRolePos[0] >= -tri.ServerSize[0]/2 && newRolePos[0] <= tri.ServerSize[0]/2 {
			if newRolePos[2] >= -tri.ServerSize[2]/2 && newRolePos[0] <= tri.ServerSize[2]/2 {
				return true
			}
		}
	} else if tri.Type == "Sphere" {
		newRolePos := [3]float64{(float64)(rolePosX)/100 - tri.ServerPos[0], 0, (float64)(rolePosZ)/100 - tri.ServerPos[2]}
		newRolePos[0] = newRolePos[0] * newRolePos[0]
		newRolePos[2] = newRolePos[2] * newRolePos[2]
		if newRolePos[0]+newRolePos[2] < tri.Radius*tri.Radius {
			return true
		}
	}
	return false
}

type RegionTriggerList struct {
	Trigger []*RegionTrigger `xml:"Trigger"`
}

type RegionScript struct {
	Id    uint32 `xml:"Id,attr"`
	Path  string `xml:"Path,attr"`
	Param string `xml:"Param,attr"`
}

type RegionScriptList struct {
	Type   string          `xml:"Type"`
	Script []*RegionScript `xml:"Script"`
}

type RegionMonsterAI struct {
	ExcludeMonsterGID string `xml:"ExcludeMonsterGID,attr"`
}

//触发怪
type RegionMonster struct {
	Name          string              `xml:"name,attr"`
	Pos           string              `xml:"Pos,attr"`
	Yaw           float64             `xml:"Yaw,attr"`
	ScriptList    []*RegionScriptList `xml:"ScriptList"`
	AI            RegionMonsterAI     `xml:"ai"`
	ServerPos     [3]float64
	ServerCastPos [3]int32
	ServerCastYaw int32
}

type RegionMonsterWave struct {
	Id      uint32           `xml:"Id,attr"`
	Monster []*RegionMonster `xml:"monster"`
}

type RegionDifficulty struct {
	Level       uint32               `xml:"level,attr"`
	MonsterWave []*RegionMonsterWave `xml:"MonsterWave"`
}

type RegionArea struct {
	Id         uint32              `xml:"Id,attr"`
	TriggerId  uint32              `xml:"TriggerId,attr"`
	ScriptList []*RegionScriptList `xml:"ScriptList"`
	Difficulty []*RegionDifficulty `xml:"Difficulty"`
}

type RegionAreaList struct {
	Area []*RegionArea `xml:"Area"`
}

type RegionFinishType struct {
	Id         uint32 `xml:"Id,attr"`
	DescribeID uint32 `xml:"DescribeID,attr"`
	Win        bool   `xml:"Win,attr"`
}

type RegionFinishTypeList struct {
	FinishType     []*RegionFinishType `xml:"FinishType"`
	AllMonsterTris map[uint32][]uint32 //遍历所有怪需要的Trigger
}

func init() {
	cfgs.RegisterHandleExt("RegionLogic/", ".xml", RegionGKFileHandler)
}

// 解析配置
func RegionGKFileHandler(source *pb.MsgBytesCfg, store *cfgs.SConfigStore) {
	gk := &RegionGK{}
	err := xml.Unmarshal(source.Filebody, gk)
	gk.InitServer()
	log.Debug(source.Filename, err)
	store.XmlMap[source.Filename] = gk
}

// 读取配置
func GetRegionLogic(RegionID uint32) *RegionGK {
	cfg := Config()
	if cfg == nil {
		return nil
	}
	key := "RegionLogic/" + strconv.FormatUint((uint64)(RegionID), 10) + "_gk"
	return cfg.XmlMap[key].(*RegionGK)
}
