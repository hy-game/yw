package service

import (
	"com/log"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type Client struct {
	conn       *grpc.ClientConn
	updateTime int64
}

func (c *Client) isActive() bool {
	return time.Now().Unix()-c.updateTime < 20
}

type Service struct {
	client map[uint16]*Client
}

func newService() *Service {
	return &Service{client: make(map[uint16]*Client)}
}

type Mgr struct {
	services map[string]*Service
	mtx      sync.Mutex
}

//NewServiceMgr	创建新的服务管理器
func NewServiceMgr() *Mgr {
	m := &Mgr{
		services: make(map[string]*Service),
	}
	return m
}

//Init	初始化
func (m *Mgr) Init(needs []string) {
	for _, v := range needs {
		m.services[v] = newService()
	}
}

//Add	添加服务
func (m *Mgr) Add(name string, id uint16, endPoint string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if s, ok := m.services[name]; ok && s != nil {
		if s.client[id] == nil || !s.client[id].isActive() {
			conn, err := grpc.Dial(endPoint, grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				log.Errorf("create rpc conn:[%s %d %s] err:%v", name, id, endPoint, err)
				return
			}

			if s.client[id] != nil && s.client[id].conn != nil {
				s.client[id].conn.Close()
			}
			s.client[id] = &Client{
				updateTime: time.Now().Unix(),
				conn:       conn,
			}

			log.Infof("add service %s:%d=%s", name, id, endPoint)
		} else {
			s.client[id].updateTime = time.Now().Unix()
		}
	}
}

func (m *Mgr) Clear() {
	for _, v := range m.services {
		for _, cli := range v.client {
			if cli.conn != nil {
				cli.conn.Close()
			}
		}
	}
}

//Get	获取服务
func (m *Mgr) Get(name string, id uint16) *grpc.ClientConn {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if s, ok := m.services[name]; ok && s != nil && s.client != nil && s.client[id] != nil && s.client[id].isActive() {
		return s.client[id].conn
	}
	return nil
}

var (
	service = NewServiceMgr()
	once    sync.Once
)

//Init	初始化
func Init(needs ...string) {
	once.Do(func() {
		service.Init(needs)
	})
}

//Add	添加服务
func Add(name string, id uint16, endPoint string) {
	service.Add(name, id, endPoint)
}

//Get	获取服务
func Get(name string, id uint16) *grpc.ClientConn {
	return service.Get(name, id)
}

func Clear() {
	service.Clear()
}
