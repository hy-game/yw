package app

import (
	"center/configs"
	"center/logic/mail"
	"center/logic/ranklist"
	"com/database/db"
	"com/database/orm"
	"com/log"
	"game/logic/role"
	"github.com/jinzhu/gorm"
	"pb"
)

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

func initDbRole() {
	dbMsg := cfg("role")
	if dbMsg == nil {
		log.Panicf("can not find db cfg:%d", 1)
		return
	}
	dbCfg := &orm.DbCfg{
		Driver:  "mysql",
		Usr:     dbMsg.User,
		Pswd:    dbMsg.Password,
		Host:    dbMsg.IP,
		Port:    int(dbMsg.Port),
		Db:      dbMsg.Name,
		Charset: "utf8",
	}
	ok := orm.Init(dbCfg, 10, 1, func(i *gorm.DB) error {
		return i.AutoMigrate(
			&RoleInDb{},
		).Error
	})
	if !ok {
		log.Panic("init orm error")
	}
	orm.Close()
}

func initDbMisc() {
	dbMsg := cfg("misc")
	if dbMsg == nil {
		log.Panicf("can not find db cfg:%d", 1)
		return
	}
	dbCfg := &orm.DbCfg{
		Driver:  "mysql",
		Usr:     dbMsg.User,
		Pswd:    dbMsg.Password,
		Host:    dbMsg.IP,
		Port:    int(dbMsg.Port),
		Db:      dbMsg.Name,
		Charset: "utf8",
	}
	ok := orm.Init(dbCfg, 10, 1, func(i *gorm.DB) error {
		return i.AutoMigrate(&mail.DBMail{},
			&role.OfflineOp{},
			&ranklist.RankLists{},
		).Error
	})
	if !ok {
		log.Panic("init orm error")
	}
	orm.Close()

	if err := db.Init(dbMsg.User, dbMsg.Password, dbMsg.IP, uint16(dbMsg.Port), dbMsg.Name); err != nil {
		log.Panicf("init db err")
	}
}
