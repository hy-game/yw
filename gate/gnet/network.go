package gnet

import (
	"com/util"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	shutDown  = make(chan struct{})
	waitGroup sync.WaitGroup
)

//Config 网络配置
type Config struct {
	ReadDeadline        time.Duration //time.Second * 1500
	OutChanSize         int           //128
	ReadSocketBuffSize  int
	WriteSocketBuffSize int //
	RpmLimit            int
	RecvPkgLenLimit     uint32
	EvtChanSize         int
}

//StartTCPServer 开始tcp服务
func StartTCPServer(listenEndPoint string, cfg *Config) {
	addr, err := net.ResolveTCPAddr("tcp4", listenEndPoint)
	checkError(err)

	listener, err := net.ListenTCP("tcp", addr)
	checkError(err)
	log.Info("start listen on:", listener.Addr())

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warn("accept failed:", err)
			continue
		}

		// set socket read buffer
		conn.SetReadBuffer(cfg.ReadSocketBuffSize)
		// // set socket write buffer
		conn.SetWriteBuffer(cfg.WriteSocketBuffSize)

		go handleClient(conn, cfg)
	}
}

func handleClient(conn net.Conn, cfg *Config) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			util.PrintStack()
		}
	}()

	var s Session
	s.Conn = conn
	defer s.Conn.Close()

	s.In = make(chan []byte) //no cache
	defer func() {
		close(s.In)
		log.Debug("recv loop stop")
	}()

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("cannot get remote address:", err)
		return
	}
	s.Ip = net.ParseIP(host)
	log.Infof("new connection from:%v port:%v", host, port)

	s.Ctrl = make(chan struct{})

	s.start(conn, cfg)
}

//Close 关闭，并等待所有goroutine退出
func Close() {
	close(shutDown)
	log.Debugf("start wait %v", waitGroup)
	waitGroup.Wait()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

/************************************************************/
var (
	CliSess sync.Map
	CliCnt  int32
)

func AddCliSession(cliSesID uint32, c *Session) {
	CliSess.Store(cliSesID, c)
	atomic.AddInt32(&CliCnt, 1)
}

func RemoveCliSession(cliSesID uint32) {
	CliSess.Delete(cliSesID)
	atomic.AddInt32(&CliCnt, -1)
}

func GetCliSessCnt() int32 {
	return atomic.LoadInt32(&CliCnt)
}

func GetCliSession(cliSesID uint32) *Session {
	s, ok := CliSess.Load(cliSesID)
	if ok {
		return s.(*Session)
	}
	return nil
}
