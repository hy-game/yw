package gmdb

import (
	"com/database/db"
	"com/log"
	"game/configs"
	"pb"
)

func Init(){
	initDbMisc()
	initDbRole()
}

//Cfg	配置
func cfg(dbType string) *pb.MsgDBCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("dbcfg", dbType, "")
	if !ok {
		log.Warnf("can not find db cfg:%d", dbType)
		return nil
	}
	return f.(*pb.MsgDBCfg)
}

var DBRole db.DB
func initDbRole() {
	dbMsg := cfg("role")
	if dbMsg == nil {
		log.Panicf("can not find db cfg:%d", 1)
		return
	}

	if err := DBRole.Init(dbMsg.User, dbMsg.Password, dbMsg.IP, uint16(dbMsg.Port), dbMsg.Name); err != nil {
		log.Panicf("init db err")
	}
}

var DBMisc db.DB
func initDbMisc() {
	dbMsg := cfg("misc")
	if dbMsg == nil {
		log.Panicf("can not find db cfg:%d", 1)
		return
	}

	if err := DBMisc.Init(dbMsg.User, dbMsg.Password, dbMsg.IP, uint16(dbMsg.Port), dbMsg.Name); err != nil {
		log.Panicf("init db err")
	}
}