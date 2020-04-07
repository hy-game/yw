package orm

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

//IOrmTask 使用gorm的异步任务接口
type IOrmTask interface {
	Run(*gorm.DB)
}

//Orm gorm封装
type Orm struct {
	tasks          chan IOrmTask
	syncSess       chan *gorm.DB
	syncSingleSess *gorm.DB
	sess           *gorm.DB
	cfg            *DbCfg

	closeChan chan struct{}
}

//Init 初始化
func (m *Orm) Init(cfg *DbCfg, asyncTaskSize int32, syncSessCnt int32, migrate func(*gorm.DB) error) bool {
	m.cfg = cfg
	conStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		cfg.Usr, cfg.Pswd, cfg.Host, cfg.Port, "", cfg.Charset)

	sess, err := gorm.Open("mysql", conStr)
	if err != nil {
		log.Panicf("connect db %s error:%v", conStr, err)
		return false
	}

	sess.LogMode(true)
	sess.SingularTable(true)

	err = sess.Exec(fmt.Sprintf("CREATE DATABASE  IF NOT EXISTS %s", cfg.Db)).Error
	if err != nil {
		return false
	}

	err = sess.Exec(fmt.Sprintf("use %s", cfg.Db)).Error
	if err != nil {
		return false
	}

	err = migrate(sess)
	if err != nil {
		log.Panicf("migrate error:%v", err)
		return false
	}

	conStrWithDb := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		cfg.Usr, cfg.Pswd, cfg.Host, cfg.Port, cfg.Db, cfg.Charset)

	m.syncSingleSess, err = gorm.Open("mysql", conStrWithDb)
	if err != nil {
		log.Panicf("connect db %s error:%v", conStrWithDb, err)
		return false
	}

	m.sess = sess
	m.tasks = make(chan IOrmTask, asyncTaskSize)
	m.closeChan = make(chan struct{})
	m.syncSess = make(chan *gorm.DB, syncSessCnt)

	go m.run(m.sess)

	for i := 0; i < cap(m.syncSess); i++ {
		sess, err := gorm.Open("mysql", conStrWithDb)
		if err != nil {
			log.Panicf("connect db %s error:%v", conStrWithDb, err)
			return false
		}
		m.syncSess <- sess
	}

	log.Infof("db init to :%v", m.cfg)
	return true
}

//Close 关闭
func (m *Orm) Close() {
	log.Info("gorm close")
	close(m.closeChan)
}

//Post 投递异步任务
func (m *Orm) Post(t IOrmTask) {
	m.tasks <- t
}

//Execute 执行同步任务
func (m *Orm) Execute(f func(sess *gorm.DB) error) error {
	sess := <-m.syncSess
	defer func() {
		m.syncSess <- sess
	}()

	return f(sess)
}

//SingleSyncExe 单线程执行同步任务
func (m *Orm) SingleSyncExe(f func(sess *gorm.DB) error) error {
	return f(m.syncSingleSess)
}

func (m *Orm) run(conn *gorm.DB) {
	defer func() {
		conn.Close()
		m.syncSingleSess.Close()
		for v := range m.syncSess {
			v.Close()
		}
	}()

	for {
		select {
		case t := <-m.tasks:
			t.Run(conn)
		case <-m.closeChan:
			return
		}
	}
}
