package cache

import "com/log"

var defaultRedis *Redis

func Init(cfg *RedisCfg) {
	if defaultRedis != nil {
		log.Errorf("default Redis already init, you can init a new one, use NewRedis")
		return
	}

	defaultRedis = NewRedis(cfg)
}

func Execute(cmd string) interface{} {
	return defaultRedis.ExecuteSync(cmd)
}
