package fb

/*
本来抽象了个玩法interface，来做统一的进入战斗和战斗结果处理，后面考虑其实每个玩法进入和结果的消息内容都不一样，没必要统一
*/
import (
	"com/log"
	"game/configs"
	"game/logic/comps/item"
	"game/types"
	"pb"
)

//Cfg 英雄配置
func CfgFB(index uint32) *pb.MsgFBCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("FBCfg", index, "")
	if !ok {
		log.Warnf("not find cfg:%d from FBCfg.txt", index)
		return nil
	}
	return f.(*pb.MsgFBCfg)
}

func CfgRegion(index uint32) *pb.MsgRegionList {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("RegionList", index, "")
	if !ok {
		log.Warnf("not find cfg:%d from RegionList.txt", index)
		return nil
	}
	return f.(*pb.MsgRegionList)
}

type Data struct {
}

func NewData(r *types.Role) *Data {
	if r.Data.FB == nil {
		r.Data.FB = &pb.MsgFBData{
			Cur:     0,
			Histroy: make(map[uint32]*pb.MsgFBHistroy),
		}
	}
	return &Data{}
}

func CanEnter(regionCfg *pb.MsgRegionList, r *types.Role)(pb.EFormationType, bool) {
	//if r.Data.FB.Cur+1 < regionCfg.ID {
	//	return pb.EFormationType_EFBNormal, false
	//}
	return pb.EFormationType_EFBNormal, true
}

func OnFinish(r *types.Role, data *pb.MsgBattleFinishData){
	data.Items = append(data.Items, &pb.CPriceItem{Mtype: pb.EItemType_Good, Oriname: 4, Count: 100})
	data.Items = append(data.Items, &pb.CPriceItem{Mtype: pb.EItemType_Good, Oriname: 3, Count: 33})
	data.Items = append(data.Items, &pb.CPriceItem{Mtype: pb.EItemType_Good, Oriname: 5, Count: 55})
	data.Items = append(data.Items, &pb.CPriceItem{Mtype: pb.EItemType_Good, Oriname: 6, Count: 66})
	for _, v := range data.Items {
		item.Prize(v, pb.ESource_System, r)
	}
}