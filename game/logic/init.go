package logic

import (
	"com/log"
	"database/sql"
	"game/configs"
	"game/gmdb"
	"game/logic/handler"
	"game/logic/role"
	"game/network"
	"game/setup"
)

//Init 初始化,在数据库初始化前
func Init() {
	network.Init()
	handler.Init()
	network.StartSerNet()
	role.Init()
	configs.Init()
}

//InitAfterDBInit 数据库初始化后再次调用该函数
func InitAfterDBInit() {
	getMaxGuidInDb()
}

const (
	MaxPlayerCnt = 0xffffff
)

func getMaxGuidInDb() {
	gmdb.DBRole.SyncExe(func(conn *sql.DB) {
		err := conn.QueryRow("SELECT guid FROM role_in_db WHERE TRUNCATE(guid / ?, 0) = ? ORDER BY guid DESC LIMIT 1", MaxPlayerCnt, setup.Setup.ID).Scan(&role.MaxGuidInDb)
		if err != nil {
			log.Info("get max role guid err %v", err)
			role.MaxGuidInDb = setup.Setup.ID * MaxPlayerCnt
		}
		log.Tracef("get max guid in db:%d", role.MaxGuidInDb)
	})
}
