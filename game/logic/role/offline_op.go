package role

import (
	"com/log"
	"database/sql"
	"fmt"
	"game/gmdb"
	"game/types"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

type OffOp int

const (
	Add OffOp = iota
	Load
)

type OfflineOpParam struct {
	op       OffOp
	evt      types.Evt
	roleGuid uint32
}

type OfflineOpMgr struct {
	ops     chan *OfflineOpParam
	loading map[uint32]int64
}

func newOfflineOperator() *OfflineOpMgr {
	m := &OfflineOpMgr{
		ops:     make(chan *OfflineOpParam, 1000),
		loading: make(map[uint32]int64),
	}
	go m.run()

	return m
}

var Mgr4OfflineOp = newOfflineOperator()

func (m *OfflineOpMgr) PostOfflineEvt(evt types.Evt, roleGuid uint32) {
	m.ops <- &OfflineOpParam{
		op:       Add,
		evt:      evt,
		roleGuid: roleGuid,
	}
}

func (m *OfflineOpMgr) Load(roleGuid uint32) {
	m.ops <- &OfflineOpParam{
		op:       Load,
		roleGuid: roleGuid,
	}
}

func (m *OfflineOpMgr) run() {
	for {
		select {
		case p := <-m.ops:
			switch p.op {
			case Add:
				m.onAdd(p)
			case Load:
				m.onLoad(p)
			}
		}
	}
}

func (m *OfflineOpMgr) onAdd(p *OfflineOpParam) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Warnf("create uuid err when offlineop add %v", err)
		return
	}

	gmdb.DBMisc.Write(&dbSaveOfflineOp{
		RoleGuid: p.roleGuid,
		OpGuid:   u,
		Type:     p.evt.Type,
		Data:     p.evt.Data,
	})
}

const (
	LoadCD = 5
)

func (m *OfflineOpMgr) onLoad(p *OfflineOpParam) {
	now := time.Now().Unix()
	if v, ok := m.loading[p.roleGuid]; ok && now-v < LoadCD {
		return
	}
	m.loading[p.roleGuid] = now

	gmdb.DBMisc.Read(&dbLoadOfflineOp{
		RoleGuid: p.roleGuid,
	})
}

type OfflineOp struct {
	UID  string `gorm:"primary_key"`
	Role uint32 `gorm:"index:idx_guid"`
	Type uint32
	Data []byte `gorm:"type:blob(1048576)"`
}

//DbLoadRole	-------------------加载角色数据-----------------------------
type dbLoadOfflineOp struct {
	RoleGuid  uint32
	RoleSesID uint32
}

func (t *dbLoadOfflineOp) Run(conn *sql.DB) {
	rows, err := conn.Query("SELECT * FROM offline_op WHERE role = ?", t.RoleGuid)
	if err != nil {
		log.Warnf("load offline op %d err:%v", t.RoleGuid, err)
		return
	}

	go func() {
		tDel := &dbDelOfflineOp{
			ids: make([]string, 0),
		}
		defer func() {
			defer rows.Close()
			gmdb.DBMisc.Write(tDel)
		}()

		uuidStr := ""
		var op OfflineOp
		for rows.Next() {
			err = rows.Scan(&uuidStr, &op.Role, &op.Type, &op.Data)
			if err != nil {
				log.Warnf("load offline op scan err:%v", err)
				continue
			}

			var msg proto.Message
			if len(op.Data) > 0 {
				err = proto.Unmarshal(op.Data, msg)
				if err != nil {
					log.Warnf("load offline op unmarshal err:%v", err)
					continue
				}
			}
			ok := types.PostToSes(t.RoleSesID, types.Evt{
				Type: types.GameEvent(op.Type),
				Data: msg,
			})
			if !ok {
				return
			} else {
				tDel.ids = append(tDel.ids, uuidStr)
			}
		}
	}()
}

type dbDelOfflineOp struct {
	ids []string
}

func (t *dbDelOfflineOp) Run(conn *sql.DB) {
	params := strings.Join(t.ids, "','")
	sqlStr := fmt.Sprintf("DELETE FROM offline_op WHERE uid IN ('%s')", params)
	_, err := conn.Exec(sqlStr)
	if err != nil {
		log.Warnf("del offline op err:%v", err)
	}
}

//dbSaveOfflineOp	---------------------保存角色离线操作数据------------------------
type dbSaveOfflineOp struct {
	RoleGuid uint32
	OpGuid   uuid.UUID
	Type     types.GameEvent
	Data     proto.Message
}

func (t *dbSaveOfflineOp) Run(conn *sql.DB) {
	_, err := conn.Exec("INSERT INTO offline_op VALUES (?, ?, ?, ?)", t.OpGuid.String(), t.RoleGuid, t.Type, t.Data)
	if err != nil {
		log.Warnf("save offline op %d err:%v", t.RoleGuid, err)
		return
	}
}
