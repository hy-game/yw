/*
角色相关的数据库操作
*/
package role

import (
	"com/log"
	"database/sql"
	"fmt"
	"game/gmdb"
	"github.com/golang/protobuf/proto"
	"pb"
)

var MaxGuidInDb uint32

func GenGuid() uint32 {
	MaxGuidInDb++
	return MaxGuidInDb
}

//DbLoadRole	-------------------加载角色数据-----------------------------
type DbLoadRole struct {
	RoleSesId uint32
	Acc       string
}

func dbLoadRole(acc string, roleSesId uint32) {
	t := &DbLoadRole{
		RoleSesId: roleSesId,
		Acc:       acc,
	}
	gmdb.DBRole.Read(t)
}

func (t *DbLoadRole) Run(conn *sql.DB) {
	var guid, level uint32
	var name, acc string
	data := make([]byte, 0)

	err := conn.QueryRow("SELECT * FROM role_in_db WHERE acc = ?", t.Acc).Scan(&guid, &acc, &name, &level, &data)
	if err != nil {
		log.Debugf("insert role %d", t.Acc)
		guidMax := GenGuid()
		_, err := conn.Exec(fmt.Sprintf("INSERT INTO role_in_db(guid, acc, name, level) VALUES(%d, '%s', '%d', %d)", guidMax, t.Acc, guidMax, 1))
		if err != nil {
			log.Debugf("create role %s err:%v", t.Acc, err)
			return
		}

		err = conn.QueryRow(fmt.Sprintf("SELECT * FROM role_in_db WHERE acc = '%s'", t.Acc)).Scan(&guid, &acc, &name, &level, &data)
		if err != nil {
			log.Warnf("load role %s err:%v", t.Acc, err)
			return
		}
	}

	msg := &pb.MsgPlayerData{}
	if len(data) > 0 {
		err = proto.Unmarshal(data, msg)
		if err != nil {
			log.Warnf("unmarshal fail :%d", guid)
			return
		}
	}
	msg.Guid = guid
	msg.Level = level
	msg.Name = name
	msg.Account = acc

	Mgr4Role.OnLoadRoleData(msg, t.RoleSesId)
	log.Debugf("load role success [%d]%s", msg.Guid, msg.Account)
}

//DbSaveRole	---------------------保存角色数据------------------------
type DbSaveRole struct {
	Guid  uint32 `gorm:"primary_key"`
	Acc   string `gorm:"index:idx_acc"`
	Name  string
	Level uint32
	Data  []byte `gorm:"type:blob(104857600)"`
}

func (t *DbSaveRole) Run(conn *sql.DB) {
	_, err := conn.Exec("UPDATE role_in_db SET name = ?, level = ?, data = ? WHERE guid = ?", t.Name, t.Level, t.Data, t.Guid)
	if err != nil {
		log.Warnf("save role %s err:%v", t.Guid, err)
		return
	}
	log.Debugf("save role success [%d]%s", t.Guid, t.Acc)
}
