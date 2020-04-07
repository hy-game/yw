package configs

import "com/cfgs"

//获取服务器玩法配置
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
