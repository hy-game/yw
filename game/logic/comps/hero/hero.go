package hero

import (
	"com/log"
	"game/logic/comps/item"
	"game/types"
	"pb"
)

func heroLevelUp(msg *pb.MsgHeroUpReq, r *types.Role, h *pb.MsgHero) {
	exp := uint32(0)
	for _, v := range msg.Cost {
		itemCfg := item.Cfg(v.ID)
		if itemCfg == nil {
			log.Warnf("can not find item:%d", v.ID)
			return
		}

		if !item.Enough(v.ID, v.Cnt, r) {
			return
		}

		exp += itemCfg.Param1 * v.Cnt
		if heroAddExp(exp, h) {
			item.Del(v.ID, v.Cnt, pb.ESource_HeroUp, r)
		}
	}
}

func heroAddExp(exp uint32, h *pb.MsgHero) bool {
	heroCfg := CfgHero(h.CfgId)
	if heroCfg == nil {
		return false
	}
	starCfg := CfgHeroStar(heroCfg.Type, h.Star)
	if starCfg == nil {
		return false
	}
	if h.Level >= starCfg.LevelLimit {
		return false
	}

	levelCfg := CfgHeroLevel(heroCfg.Type, h.Level)
	if levelCfg == nil {
		return false
	}

	for levelCfg.NeedExp < h.Exp+exp {
		h.Level++
		exp = exp - (levelCfg.NeedExp - h.Exp)
		h.Exp = 0

		starCfg := CfgHeroStar(heroCfg.Type, h.Star)
		if starCfg == nil {
			return true
		}
		if h.Level >= starCfg.LevelLimit {
			return true
		}

		levelCfg = CfgHeroLevel(heroCfg.Type, h.Level)
		if levelCfg == nil {
			return true
		}
	}
	h.Exp += exp
	return true
}

func heroStarUp(msg *pb.MsgHeroStarUpReq, r *types.Role, h *pb.MsgHero) {
	heroCfg := CfgHero(h.CfgId)
	if heroCfg == nil {
		return
	}
	if msg.SpeceilCost != nil {
		//升星卡
		//todo
	} else if msg.CostHero != nil {
		starCfg := CfgHeroStar(heroCfg.Type, h.Star)
		if starCfg == nil {
			return
		}

		if item.GetJinBi(r) < starCfg.JinBi {
			return
		}

		cnt := uint32(0)
		for _, v := range msg.CostHero {
			hero := getHero(v, r)
			if hero == nil {
				return
			}

			if h.Star < starCfg.HeroStar {
				return
			}
			cnt++
		}

		if cnt < starCfg.HeroCnt {
			return
		}

		for _, v := range msg.CostHero {
			Del(v, r)
		}
		item.DelJinBi(starCfg.JinBi, pb.ESource_HeroUp, r)

		h.Star++
	}
}

func heroPinJieUp(msg *pb.MsgHeroPinJieUpReq, r *types.Role, h *pb.MsgHero) {
	switch msg.Type {
	case pb.MsgHeroPinJieUpReq_Normal:
		heroPinJieUpNormal(r, h)
	case pb.MsgHeroPinJieUpReq_ShengPinKa: //todo 升品卡
	case pb.MsgHeroPinJieUpReq_ZhiShengKa:
	}
}

func heroPinJieUpNormal(r *types.Role, h *pb.MsgHero){
	heroCfg := CfgHero(h.CfgId)
	if heroCfg == nil {
		return
	}

	pinJieCfg := CfgHeroPinJie(heroCfg.Type, h.PinJie, heroCfg.Job)
	if pinJieCfg == nil {
		return
	}

	if !item.Enough(pinJieCfg.Cost, pinJieCfg.CostCnt, r) {
		return
	}

	if item.GetJinBi(r) < pinJieCfg.JinBi {
		return
	}

	item.Del(pinJieCfg.Cost, pinJieCfg.CostCnt, pb.ESource_HeroUp, r)
	item.DelJinBi(pinJieCfg.JinBi, pb.ESource_HeroUp, r)

	h.PinJie++
}