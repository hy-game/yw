package logic

import (
	"github.com/jinzhu/gorm"
	"pb"
)

func Migrate(conn *gorm.DB) error{
	return conn.AutoMigrate(&pb.MsgPlayerData{}).Error
}

func Init(){
	RegisteMsgHandle()
	initCliMsg()
}