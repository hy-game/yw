package configs

import (
	"com/cfgs"
	"com/log"
	"pb"
	"sync/atomic"
)

type value struct {
	atomic.Value
}

var config value
var loadFinish = make(chan bool)

//初始化配置信息 有服务器启动时调用
func Init() {
	//再请求配置
	log.Infof("Config all load begin...")
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
	config.update(in)

	log.Infof("Config all load finish...")
	<-loadFinish
}

//获取服务器玩法配置
func Config() *cfgs.SConfigStore {
	iv := config.Load()
	if iv == nil {
		return nil
	}
	return iv.(*cfgs.SConfigStore)
}

//获取特定行 特定列的配置数据
func GetValue(table_name string, index interface{}, column string) (interface{}, bool) {
	cfg := Config()
	if cfg == nil {
		return nil, false
	}
	return cfg.GetValue(table_name, index, column)
}
