package hero

import (
	"com/log"
	"com/util"
	"game/configs"
	"game/types"
	"github.com/golang/protobuf/proto"
	"pb"
)

//Cfg 英雄配置
func CfgHero(index uint32) *pb.MsgHeroCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("HeroAttri", index, "")
	if !ok {
		log.Warnf("not find cfg:%d from HeroAttri.txt", index)
		return nil
	}
	return f.(*pb.MsgHeroCfg)
}

//Cfg 英雄升级配置
func CfgHeroLevel(t uint32, level uint32) *pb.MsgHeroLevelCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("HeroLevel", util.ToString(t)+","+util.ToString(level), "")
	if !ok {
		log.Warnf("not find cfg:%d,%d from HeroLevel.txt", t, level)
		return nil
	}
	return f.(*pb.MsgHeroLevelCfg)
}

//Cfg 英雄升星配置
func CfgHeroStar(t uint32, star uint32) *pb.MsgHeroStarCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("HeroStar", util.ToString(t)+","+util.ToString(star), "")
	if !ok {
		log.Warnf("not find cfg:%d,%d from HeroStar.txt", t, star)
		return nil
	}
	return f.(*pb.MsgHeroStarCfg)
}

//Cfg 英雄升品配置
func CfgHeroPinJie(t uint32, pinJie uint32, job uint32) *pb.MsgHeroPinJieCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("HeroPinJie", util.ToString(t)+","+util.ToString(pinJie)+","+util.ToString(job), "")
	if !ok {
		log.Warnf("not find cfg:%d,%d,%d from HeroPinJie.txt", t, pinJie, job)
		return nil
	}
	return f.(*pb.MsgHeroPinJieCfg)
}

type HeroData struct {
	GuidMax uint32
}

func NewData(r *types.Role) *HeroData {
	if r.Data.Heros == nil {
		r.Data.Heros = make(map[uint32]*pb.MsgHero)
	}

	max := uint32(0)
	for k, _ := range r.Data.Heros {
		if k > max {
			max = k
		}
	}
	return &HeroData{GuidMax: max}
}

func Data(r *types.Role) *HeroData {
	data := r.GetComp(types.TCHero)
	if data == nil {
		return nil
	} else {
		return data.(*HeroData)
	}
}

//Add	添加一个武将
func Add(cfgID uint32, r *types.Role) {
	if cfg := CfgHero(cfgID); cfg == nil {
		log.Warnf("can not find hero cfg :%d", cfgID)
		return
	}

	data := Data(r)
	if data == nil {
		return
	}
	data.GuidMax++
	h := &pb.MsgHero{
		CfgId:  cfgID,
		Guid:   data.GuidMax,
		Star:   1,
		PinJie: 1,
		Level:  1,
		Exp:    0,
	}
	r.Data.Heros[h.Guid] = h

	SendHeroInfo(pb.MsgHeroProto_Add, h, r)
}

//Del 删除一个武将
func Del(guid uint32, r *types.Role) {
	delete(r.Data.Heros, guid)
	msg := pb.MsgHeroProto{
		Op:   pb.MsgHeroProto_Del,
		Guid: guid,
	}
	r.Send(pb.MsgIDS2C_S2CHeroAck, &msg)
}

func SendHeroInfo(op pb.MsgHeroProto_Operator, hero *pb.MsgHero, r *types.Role) {
	msg := pb.MsgHeroProto{
		Op:   op,
		Hero: hero,
	}
	r.Send(pb.MsgIDS2C_S2CHeroAck, &msg)
}

func getHero(guid uint32, r *types.Role) *pb.MsgHero {
	if h, ok := r.Data.Heros[guid]; ok {
		return h
	} else {
		return nil
	}
}

//OnProto 客户端发来的协议处理
func OnProto(msgBase proto.Message, r *types.Role) {
	msg := msgBase.(*pb.MsgHeroProto)
	if msg == nil {
		return
	}
	switch msg.Op {
	case pb.MsgHeroProto_LevelUp:
		subMsg := msg.UpReq
		if subMsg == nil {
			log.Warnf("recv MsgHeroProto without UpReq when LevelUp:%d", r.Guid)
			return
		}
		levelUp(subMsg, r)
	case pb.MsgHeroProto_PinJieUp:
		subMsg := msg.PinJieReq
		if subMsg == nil {
			log.Warnf("recv MsgHeroProto with without when PinJieUp:%d", r.Guid)
			return
		}
		pinJieUp(subMsg, r)
	case pb.MsgHeroProto_StarUp:
		subMsg := msg.StarReq
		if subMsg == nil {
			log.Warnf("recv MsgHeroProto without UpReq when StarUp:%d", r.Guid)
			return
		}
		starUp(subMsg, r)
	}
}

func levelUp(msg *pb.MsgHeroUpReq, r *types.Role) {
	h := getHero(msg.Guid, r)
	if h == nil {
		log.Warnf("can not find hero :%d", msg.Guid)
		return
	}
	heroLevelUp(msg, r, h)
	SendHeroInfo(pb.MsgHeroProto_LevelUp, h, r)
}

func pinJieUp(msg *pb.MsgHeroPinJieUpReq, r *types.Role) {
	h := getHero(msg.Guid, r)
	if h == nil {
		log.Warnf("can not find hero :%d", msg.Guid)
		return
	}
	heroPinJieUp(msg, r, h)
	SendHeroInfo(pb.MsgHeroProto_PinJieUp, h, r)
}

func starUp(msg *pb.MsgHeroStarUpReq, r *types.Role) {
	h := getHero(msg.Guid, r)
	if h == nil {
		log.Warnf("can not find hero :%d", msg.Guid)
		return
	}
	heroStarUp(msg, r, h)
	SendHeroInfo(pb.MsgHeroProto_StarUp, h, r)
}

func MakeFighter(heroID uint32, role *types.Role)*pb.MsgFighter{
	h := getHero(heroID, role)
	if h == nil {
		log.Warnf("role %d has not hero %d when %s", role.Guid, heroID, util.FuncCaller(2))
		return nil
	}

	hCfg := CfgHero(h.CfgId)
	if hCfg == nil {
		return nil
	}

	hLevelCfg := CfgHeroLevel(hCfg.Type, h.Level)
	if hLevelCfg == nil {
		return nil
	}

	hPinJieCfg := CfgHeroPinJie(hCfg.Type, h.PinJie, hCfg.Job)
	if hPinJieCfg == nil {
		return nil
	}

	hStarCfg := CfgHeroStar(hCfg.Type, h.Star)
	if hStarCfg == nil {
		return nil
	}

	rate := hLevelCfg.Rate + hPinJieCfg.Rate + hStarCfg.Rate
	f := &pb.MsgFighter{Attrs:make([]uint32, int(pb.AttrType_AttrMax))}
	f.ID = h.CfgId

	f.Attrs[uint32(pb.AttrType_AtkPhy)] = uint32(hCfg.Atk + hCfg.AtkGrow*rate)
	f.Attrs[uint32(pb.AttrType_AtkMagic)] = uint32(hCfg.Atk + hCfg.AtkGrow*rate)
	f.Attrs[uint32(pb.AttrType_DefPhy)] = uint32(hCfg.DefPhy + hCfg.DefPhyGrow*rate)
	f.Attrs[uint32(pb.AttrType_DefMagic)] = uint32(hCfg.DefMagic + hCfg.DefMagicGrow*rate)
	f.Attrs[uint32(pb.AttrType_InitHp)] = uint32(hCfg.HP + hCfg.HPGrow*rate)
	f.Attrs[uint32(pb.AttrType_AtkSpeed)] = uint32(hCfg.AtkSpeed)
	f.Attrs[uint32(pb.AttrType_MoveSpeed)] = uint32(hCfg.Speed)

	index := 0
	size := len(hPinJieCfg.Attr)
	for i := pb.AttrType_Crit;  i <= pb.AttrType_AtkFeng && index < size; i++ {
		f.Attrs[uint32(i)] += uint32(hPinJieCfg.Attr[index])
		index++
	}

	index = 0
	size = len(hStarCfg.Attr)
	for i := pb.AttrType_Crit;  i <= pb.AttrType_AtkFeng && index < size; i++ {
		f.Attrs[uint32(i)] += uint32(hStarCfg.Attr[index])
		index++
	}

	return f
}