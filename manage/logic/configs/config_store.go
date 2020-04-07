package configs

import (
	clog "com/log"
	"github.com/golang/protobuf/proto"
	"log"
	"manage/network"
	"pb"
	share "share"
	"sync/atomic"
)

//保存原始配置文件的结构
type mapCfgs map[string]*pb.MsgOriginalCfgs

//正在使用的配置文件
var origincfg atomic.Value

//正在读取的配置文件
var loadingcfg mapCfgs

//广播所有配置
func broadCastCfg() {
	iv := origincfg.Load()
	if iv == nil {
		clog.Warnf("broadcast config error")
		return
	}
	cfgs := iv.(mapCfgs)
	ucfg, uok := cfgs[topicUtil]
	for svr, cfg := range cfgs {
		switch svr {
		case share.GameTopic:
			{
				sendCfg := proto.Clone(cfg)
				if uok {
					proto.Merge(sendCfg, ucfg)
				}
				network.SendToGameAll(pb.MsgIDS2S_MsgBroadcastCfgs, sendCfg)
			}
		case share.CenterTopic:
			{
				sendCfg := proto.Clone(cfg)
				if uok {
					proto.Merge(sendCfg, ucfg)
				}
				network.SendToCenter(pb.MsgIDS2S_MsgBroadcastCfgs, sendCfg)
			}
		case share.FightTopic:
			{
				sendCfg := proto.Clone(cfg)
				if uok {
					proto.Merge(sendCfg, ucfg)
				}
				network.SendToFightAll(pb.MsgIDS2S_MsgBroadcastCfgs, sendCfg)
			}
		case topicYYAct:
			{
				network.SendToGameAll(pb.MsgIDS2S_MsgBroadcastYYAct, cfg)
				network.SendToCenter(pb.MsgIDS2S_MsgBroadcastYYAct, cfg)
			}
		case topicUtil:
		default:
			log.Fatalf("no handle server %v", svr)
		}
	}
}

func broadCastYYAct() {
	iv := origincfg.Load()
	if iv == nil {
		clog.Warnf("broadcast config error")
		return
	}
	cfgs := iv.(mapCfgs)
	cfg, ok := cfgs[topicYYAct]
	if ok {
		network.SendToGameAll(pb.MsgIDS2S_MsgBroadcastYYAct, cfg)
		network.SendToCenter(pb.MsgIDS2S_MsgBroadcastYYAct, cfg)
	}
}

//指定服务器请求配置
func BroadCastToServer(svr string, id uint16, yyact bool) {
	iv := origincfg.Load()
	if iv == nil {
		return
	}
	tpc := svr
	if yyact {
		tpc = topicYYAct
	}
	cfgs := iv.(mapCfgs)
	cfg, ok := cfgs[tpc]
	if !ok {
		log.Printf("no handle server %v from request", svr)
		return
	}
	msgid := pb.MsgIDS2S_MsgBroadcastCfgs
	if yyact {
		msgid = pb.MsgIDS2S_MsgBroadcastYYAct
	}
	sendcfg := proto.Clone(cfg)
	if !yyact {
		ucfg, uok := cfgs[topicUtil]
		if uok {
			proto.Merge(sendcfg, ucfg)
		}
	}
	switch svr {
	case share.GameTopic:
		{
			if id == 0 {
				network.SendToGameAll(msgid, sendcfg)
			} else {
				network.SendToGame(id, msgid, sendcfg)
			}
		}
	case share.CenterTopic:
		{
			network.SendToCenter(msgid, sendcfg)
		}
	case share.FightTopic:
		{
			if id == 0 {
				network.SendToFightAll(msgid, sendcfg)
			} else {
				network.SendToFight(id, msgid, sendcfg)
			}
		}
	case topicUtil:
	default:
		log.Fatalf("no handle server %v", svr)
	}
}

func GetCfgs(svr string) *pb.MsgOriginalCfgs {
	iv := origincfg.Load()
	if iv == nil {
		return nil
	}

	cfgs := iv.(mapCfgs)
	cfg, ok := cfgs[svr]
	if !ok {
		return nil
	}
	return cfg
}
