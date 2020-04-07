package item

import (
	"com/log"
	"game/configs"
	"game/types"
	"pb"
)

//
func getSpecialItem(index pb.ESpecialItem) uint32 {
	cfg := configs.Config()
	if cfg == nil {
		return 0
	}
	f, ok := cfg.GetValue("SpecialItem", pb.ESpecialItem_name[int32(index)], "")
	if !ok {
		log.Warnf("can not find Special item cfg:%d", index)
		return 0
	}
	return f.(*pb.MsgSpecialItemCfg).ID
}

//-------------------金币------------------
//AddJinBi 添加金币
func AddJinBi(cnt uint32, opType pb.ESource, r *types.Role) {
	Add(getSpecialItem(pb.ESpecialItem_JinBi), cnt, opType, r)
}

//DelJinBi	减少金币
func DelJinBi(cnt uint32, opType pb.ESource, r *types.Role) {
	Del(getSpecialItem(pb.ESpecialItem_JinBi), cnt, opType, r)
}

//GetJinBi	角色拥有的金币
func GetJinBi(r *types.Role) uint32 {
	return Get(getSpecialItem(pb.ESpecialItem_JinBi), r)
}
