package cache

import (
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type RedisCfg struct {
	Pswd         string
	Host         string
	Port         int
	Db           int
	MaxTask      int
	SyncSessSize int
}

type Redis struct {
	conn      redis.Conn
	tasks     chan IRedisTask
	cfg       *RedisCfg
	syncSess  chan redis.Conn
	closeChan chan struct{}
}

func NewRedis(cfg *RedisCfg) *Redis {
	m := &Redis{
		tasks:     make(chan IRedisTask, cfg.MaxTask),
		closeChan: make(chan struct{}),
		syncSess:  make(chan redis.Conn, cfg.SyncSessSize),
		cfg:       cfg,
	}
	addr := m.cfg.Host + strconv.Itoa(m.cfg.Port)

	m.conn = NewRedisConn(addr, cfg.Pswd, cfg.Db)
	if m.conn == nil {
		return nil
	}
	//sync sess
	for i := 0; i < cap(m.syncSess); i++ {
		c := NewRedisConn(addr, cfg.Pswd, cfg.Db)
		if c == nil {
			return nil
		}
		m.syncSess <- c
	}

	go m.run(m.conn)

	return m
}

func NewRedisConn(endPoint, pswd string, index int) redis.Conn {
	c, err := redis.Dial("tcp", endPoint)
	if err != nil {
		log.Errorf("connect redis %s error:%v", endPoint, err)
		return nil
	}
	return c
}

func (m *Redis) ExecuteSync(cmd string) interface{} {
	s := <-m.syncSess
	r, err := s.Do(cmd)
	if err != nil {
		log.Errorf("redis command[%s] err:%v", cmd, err)
		return nil
	}

	m.syncSess <- s

	return r
}

func (m *Redis) Post(t IRedisTask) {
	m.tasks <- t
}

//-----------------------------------------
func (m *Redis) ping() bool {
	r, err := redis.String(m.conn.Do("Ping"))
	if err != nil || r != "Pong" {
		log.Errorf("redis disconnect:%v", err)
		return false
	}
	return true
}

func (m *Redis) run(conn redis.Conn) {
	t := time.NewTicker(time.Minute * 1)
	defer func() {
		t.Stop()
		conn.Close()
	}()

	for {
		select {
		case t := <-m.tasks:
			Run(conn)
		case <-t.C:
			if !m.ping() {
				//	m.conn = NewRedisConn()
			}
		case <-m.closeChan:
			return
		}
	}
}
