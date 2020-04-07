package item

import (
	"com/log"
	"game/configs"
	"game/types"
	"math"
	"pb"
)

//Cfg	道具配置
func Cfg(index uint32) *pb.MsgItemCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("GoodsList", index, "")
	if !ok {
		log.Warnf("can not find item cfg:%d", index)
		return nil
	}
	return f.(*pb.MsgItemCfg)
}

type Data struct {
}

//NewData	一些初始化操作
func NewData(r *types.Role) *Data {
	if r.Data.Items == nil {
		r.Data.Items = make(map[uint32]uint32)
	}
	return &Data{}
}

//Add	添加道具
func Add(ID uint32, cnt uint32, opType pb.ESource, r *types.Role) {
	if !checkCfg(ID) {
		log.Warnf("item %d not in cfg", ID)
		return
	}
	if math.MaxUint32 - r.Data.Items[ID] < cnt {
		log.Errorf("%d item %d > maxuint32", r.Guid, ID)
		r.Data.Items[ID] = math.MaxUint32
	} else {
		r.Data.Items[ID] += cnt
	}

	r.Send(pb.MsgIDS2C_S2CItem, &pb.MsgItem{
		ID:  ID,
		Cnt: r.Data.Items[ID],
	})
}

//Prize 奖励，可能是道具，英雄，装备。。。等等
func Prize(item *pb.CPriceItem, opType pb.ESource, r *types.Role) {
	if item.Mtype != pb.EItemType_Good {
		return //不是物品
	}
	if item.Count <= 0 {
		return //没有实际添加物品
	}

	Add(uint32(item.Oriname), uint32(item.Count), opType, r)
}

//Enough	道具是否够
func Enough(ID uint32, cnt uint32, r *types.Role) bool {
	if v, ok := r.Data.Items[ID]; ok {
		return v >= cnt
	} else {
		return false
	}
}

//Get	获得道具数量
func Get(ID uint32, r *types.Role) uint32 {
	if v, ok := r.Data.Items[ID]; ok {
		return v
	} else {
		return 0
	}
}

//DeL	删除道具
func Del(ID uint32, cnt uint32, opType pb.ESource, r *types.Role) bool {
	if Enough(ID, cnt, r) {
		r.Data.Items[ID] -= cnt

		r.Send(pb.MsgIDS2C_S2CItem, &pb.MsgItem{
			ID:  ID,
			Cnt: r.Data.Items[ID],
		})

		return true
	} else {
		return false
	}
}

func checkCfg(ID uint32) bool {
	if cfg := Cfg(ID); cfg != nil {
		return true
	} else {
		return false
	}
}
