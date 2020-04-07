package ranklist

import (
	"com/log"
	"pb"
	"sort"
	"time"
)

type rankList struct {
	typ   pb.ERankListType               //排行榜类型
	data  []*pb.MsgRankListData          //排行榜数据
	help  map[uint32]int32               //排行榜id和rank的映射
	cache map[uint32]*pb.MsgRankListData //排行榜临时数据 一段时间后同步
	next  uint64                         //下次同步排行榜时间
}

func (rl rankList) Len() int {
	return len(rl.data)
}

func (rl rankList) Less(i, j int) bool {
	return rl.data[i].Score > rl.data[j].Score
}

func (rl rankList) Swap(i, j int) {
	rl.data[i], rl.data[j] = rl.data[j], rl.data[i]
}

func (rl *rankList) run() {
	if rl.next <= 0 {
		return //实时刷新的排行榜
	}
	cur := uint64(time.Now().Unix())
	if cur >= rl.next {
		//执行更新
		add := rl.rebuild()
		if add > 0 {
			rl.next = cur + uint64(add)
		}
	}
}

func (rl *rankList) update(data *pb.MsgRankListData, force bool) {
	rl.cache[data.Id] = data
	if force || rl.next <= 0 {
		//实时刷新
		rl.rebuild()
	}
}

func (rl *rankList) rebuild() uint32 {
	cfg := GetCfg(rl.typ)
	if cfg == nil {
		log.Warnf("Create Ranklist with type [%v] without cfg", rl.typ)
		return 0
	} else {
		if rl.merge(cfg.MaxLen) {
			sort.Sort(rl)
			if len(rl.data) > int(cfg.MaxLen) {
				rl.data = rl.data[:cfg.MaxLen]
			}
			rl.kv()
		}
		return cfg.RefreshSpan
	}
}

func (rl *rankList) merge(max int32) bool {
	if len(rl.cache) == 0 {
		return false
	}
	for k, v := range rl.cache {
		pos, ok := rl.help[k]
		if ok {
			rl.data[pos] = v
		} else {
			rl.data = append(rl.data, v)
		}
	}
	//清空cache
	rl.cache = make(map[uint32]*pb.MsgRankListData)
	return true
}

func (rl *rankList) kv() {
	rl.help = make(map[uint32]int32)
	//fmt.Printf("len [%v]", len(rl.data))
	//fmt.Printf("cap [%v]", cap(rl.data))
	for i, v := range rl.data {
		rl.help[v.Id] = int32(i)
	}
}

func (rl *rankList) pack(guid uint32) *pb.MsgRankListPack {
	cfg := GetCfg(rl.typ)
	if cfg == nil {
		log.Warnf("Pack Ranklist with type [%v] without cfg", rl.typ)
		return nil
	}
	msg := &pb.MsgRankListPack{
		Type: rl.typ,
	}
	//我的排名
	rank, ok := rl.help[guid]
	if ok {
		msg.MyRank = rank
	} else {
		msg.MyRank = -1 //表示未上榜
	}
	//排行榜数据
	for i, v := range rl.data {
		if i >= int(cfg.ShowLen) {
			break
		}
		msg.Data = append(msg.Data, v)
	}
	return msg
}

func newRankList(typ pb.ERankListType) *rankList {
	cfg := GetCfg(typ)
	if cfg == nil {
		log.Warnf("Create Ranklist with type [%v] without cfg", typ)
		return nil
	}
	rl := &rankList{
		typ:   typ,
		data:  make([]*pb.MsgRankListData, 0, cfg.MaxLen),
		help:  make(map[uint32]int32),
		cache: make(map[uint32]*pb.MsgRankListData),
	}

	if cfg.RefreshSpan > 0 {
		rl.next = uint64(time.Now().Unix()) + uint64(cfg.RefreshSpan)
	}

	return rl
}
