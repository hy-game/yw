package configs

import (
	"center/network"
	"com/cfgs"
	"com/log"
	"pb"
	"strings"
	"sync/atomic"
)

type value struct {
	atomic.Value
}

var config value
var yyact value

var loadFinish = make(chan bool)

//模块初始化 启服时调用
func Init() {
	//先注册配置加载后事件
	cfgs.RegisterAfterFunc(AfterLoad)
	//再请求配置
	network.SendToManage(pb.MsgIDS2S_Ct2MaConfigReq, nil)
	network.SendToManage(pb.MsgIDS2S_Ct2MaConfigYYAct, nil)
	log.Infof("Config all load begin...")
	loadFinish <- true
	loadFinish <- true
	close(loadFinish)
}

func (v *value) update(in *pb.MsgOriginalCfgs) {
	loading := cfgs.NewCfgStore()
	loading.Reload(in)
	v.Store(loading)
}

//收到manage功能配置 更新配置
func UpdateCfg(in *pb.MsgOriginalCfgs) {
	init := false
	if config.Load() == nil {
		init = true
	}
	config.update(in)

	log.Infof("Config all load finish...")
	if init {
		<-loadFinish
	}
}

//收到manage运营配置 更新配置
func UpdateYYAct(in *pb.MsgOriginalCfgs) {
	init := false
	if yyact.Load() == nil {
		init = true
	}
	yyact.update(in)

	log.Infof("yyact config all load finish...")
	if init {
		<-loadFinish
	}
}

//获取功能配置
func Config() *cfgs.SConfigStore {
	iv := config.Load()
	if iv == nil {
		return nil
	}
	return iv.(*cfgs.SConfigStore)
}

//获取运营活动配置
func YYAct() *cfgs.SConfigStore {
	iv := yyact.Load()
	if iv == nil {
		return nil
	}
	return iv.(*cfgs.SConfigStore)
}

//解析完配置后 对部分数据做特殊处理
func AfterLoad(cfg *cfgs.SConfigStore) {
	giftPackGloup(cfg)
}

//礼包码组特殊处理的地方
func giftPackGloup(cfg *cfgs.SConfigStore) {
	cfg.MapTable("GiftPack/channelgroup", "", func(source interface{}) bool {
		cgcfg := source.(*pb.MsgGiftChannelGroup)
		if len(cgcfg.Chans) <= 0 {
			return false
		}
		strs := strings.Split(cgcfg.Chans[0], ",")
		cgcfg.Chans = make([]string, len(strs))
		cgcfg.Chans = append(cgcfg.Chans, strs...)
		return true
	})
}
