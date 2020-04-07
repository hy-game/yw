package configs

import (
	"com/cfgs"
	"com/log"
	"game/network"
	"pb"
	"sync/atomic"
)

type value struct {
	atomic.Value
}

var config value
var yyact value
var loadFinish = make(chan bool)

//初始化配置信息 有服务器启动时调用
func Init() {
	//先注册配置加载后事件
	//cfgs.RegisterAfterFunc(AfterLoad)
	//再请求配置
	network.SendToManage(pb.MsgIDS2S_Gm2MaConfigReq, nil)
	network.SendToManage(pb.MsgIDS2S_Gm2MaConfigYYAct, nil)
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

//更新服务器基本配置 由manage通知调用
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

//更新运营活动配置 由manage通知调用
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
