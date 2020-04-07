package ranklist

import (
	"bytes"
	"center/configs"
	"center/network"
	"com/database/db"
	"com/log"
	"database/sql"
	"github.com/golang/protobuf/proto"
	"pb"
	"share"
	"sync"
	"time"
)

//排行榜数据
var ranklists = make(map[pb.ERankListType]*rankList)

//排行榜操作
var oper = make(chan *pb.MsgRanklistHandle)

//获取排行榜配置
func GetCfg(typ pb.ERankListType) *pb.MsgRankListCfg {
	cfg := configs.Config()
	if cfg == nil {
		return nil
	}
	f, ok := cfg.GetValue("RankList", typ, "")
	if !ok {
		log.Warnf("Get RankList Cfg with type [%v] fail in RankList.txt", typ)
		return nil
	}
	return f.(*pb.MsgRankListCfg)
}

//初始化排行榜
func Init() {
	Load()
	go run()
}

func run() {
	c := time.Tick(time.Minute)
	for {
		select {
		case <-c:
			{
				for _, v := range ranklists {
					v.run()
				}
			}
		case task := <-oper:
			{
				handle_task(task)
			}
		}
	}
}

func handle_task(task *pb.MsgRanklistHandle) {
	switch task.Oper {
	case pb.MsgRanklistHandle_Update: //更新排行榜数据
		{
			rl, ok := ranklists[task.Type]
			if !ok {
				rl = newRankList(task.Type)
				ranklists[task.Type] = rl
			}
			if rl == nil {
				log.Warnf("ranklist task create fail with type [%v]", task.Oper)
				return
			}
			rl.update(task.Data, task.Force)
			//强制刷新排行榜会存盘
			if task.Force {
				save()
			}
		}
	case pb.MsgRanklistHandle_Pack: //打包排行榜数据
		{
			rl, ok := ranklists[task.Type]
			if !ok {
				sendRet(task, nil)
				return
			}
			id := uint32(0)
			if task.Data != nil {
				id = task.Data.Id
			}
			msg := rl.pack(id)
			sendRet(task, msg)
		}
	default:
		log.Warnf("ranklist task no handle with type [%v]", task.Oper)
	}
}

func sendRet(task *pb.MsgRanklistHandle, msg proto.Message) {
	switch task.Server {
	case share.ManageTopic:
		network.SendToManage(pb.MsgIDS2S(task.Msgid), msg)
	default:
		network.SendToRole(task.Data.Id, pb.MsgIDS2C(task.Msgid), msg)
	}
}

//注册排行榜任务
func PushTask(task *pb.MsgRanklistHandle) {
	oper <- task
}

//存盘
func save() {
	sql := bytes.NewBuffer(nil)
	sql.WriteString("REPLACE INTO `rank_lists` VALUES")
	first := true
	params := make([]interface{}, 0)
	for k, v := range ranklists {
		if !first {
			sql.WriteString(",")
		} else {
			first = false
		}
		msg := &pb.MsgRankListPack{
			Type: k,
		}
		for _, rl := range v.data {
			msg.Data = append(msg.Data, rl)
		}
		b, err := proto.Marshal(msg)
		if err != nil {
			log.Errorf("marshal ranklist with type %v data err:%v", k, err)
			continue
		}
		sql.WriteString("(?, ?)")
		params = append(params, k, b)
	}
	if first {
		return //没有排行榜
	}
	sql.WriteString(";")

	save := &RankListSave{
		Sql:    sql.String(),
		Params: params,
	}
	db.Write(save)
}

type RankListSave struct {
	Sql    string
	Params []interface{}
}

func (t *RankListSave) Run(conn *sql.DB) {
	_, err := conn.Exec(t.Sql, t.Params...)
	if err != nil {
		log.Warnf("Save RankList data error %v", err)
		return
	}
	log.Debugf("Save RankList data Finish")
}

//读盘
func Load() {
	load := &RankLists{}
	load.wg.Add(1)
	db.Read(load)
	load.wg.Wait()
}

type RankLists struct {
	Type uint32 `gorm:"primary_key;auto_increment:false"`
	Data []byte `gorm:"type:blob(655350)"`
	wg   sync.WaitGroup
}

func (t *RankLists) Run(conn *sql.DB) {
	rows, err := conn.Query("SELECT * FROM `rank_lists`;")
	if err != nil {
		log.Warnf("Load RankList Data Error [%v]", err)
	}
	ranklists = make(map[pb.ERankListType]*rankList)
	for rows.Next() {
		var typ uint32
		data := make([]byte, 0)
		rows.Scan(&typ, &data)
		rl := &rankList{
			typ: pb.ERankListType(typ),
		}
		cfg := GetCfg(rl.typ)
		if cfg == nil {
			log.Warnf("unmarshal ranklist getcfg fail :%v", typ)
			continue
		}
		if len(data) > 0 {
			msg := &pb.MsgRankListPack{}
			err := proto.Unmarshal(data, msg)
			if err != nil {
				log.Warnf("unmarshal ranklist data fail :%v", typ)
				continue
			}
			for _, v := range msg.Data {
				rl.data = append(rl.data, v)
				rl.kv()
			}
		}
		if cfg.RefreshSpan > 0 {
			rl.next = uint64(time.Now().Unix()) + uint64(cfg.RefreshSpan)
		}
		ranklists[rl.typ] = rl
	}
	t.wg.Done()
}
