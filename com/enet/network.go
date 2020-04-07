package enet

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"sync"
	"time"
)

var (
	shutDown  = make(chan struct{})
	waitGroup sync.WaitGroup
)

//Config 网络配置
type Config struct {
	ReadDeadline time.Duration //time.Second * 1500
	OutChanSize  int           //128
	InChanSize   int

	RpmLimit            int
	EvtChanSize         int
	ReadSocketBuffSize  int
	WriteSocketBuffSize int //32767
	RecvPackegLimit     uint32
	SendPacketLimit     uint32
}

//StartTCPServer 开始tcp服务
func StartTCPServer(listenEndPoint string, netDlgt INet) {
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
		conn.SetReadBuffer(netDlgt.GetCfg().ReadSocketBuffSize)
		// // set socket write buffer
		conn.SetWriteBuffer(netDlgt.GetCfg().WriteSocketBuffSize)

		var s Connection
		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Error("cannot get remote address:", err)
			return
		}
		s.Ip = net.ParseIP(host)
		log.Infof("new connection from:%v port:%v", host, port)

		s.Cfg = netDlgt.GetCfg()
		netDlgt.OnCreateSession(&s)

		go s.start(conn)
	}
}

//Close 关闭，并等待所有goroutine退出
func Close() {
	close(shutDown)
	log.Debugf("start wait %v", waitGroup)
	waitGroup.Wait()
}

//DailTCP 连接服务器，
//todo connection 协程可能没启动，但返回的session就已经在发送数据了
func DailTCP(rAddr string, localAddr string, netDlgt INet) (ISession, error) {
	tcpRAddr, err := net.ResolveTCPAddr("tcp4", rAddr)
	checkError(err)

	var tcpLAddr *net.TCPAddr

	if localAddr != "" {
		tcpLAddr, err = net.ResolveTCPAddr("tcp4", localAddr)
		checkError(err)
	}

	conn, err := net.DialTCP("tcp", tcpLAddr, tcpRAddr)
	if err != nil {
		log.Warnf("dailtcp[%v] error：%v", rAddr, err)
		return nil, err
	}

	c := Connection{}
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("cannot get remote address:", err)
		return nil, err
	}
	c.Ip = net.ParseIP(host)
	log.Infof("connection to:%v port:%v", host, port)

	c.Cfg = netDlgt.GetCfg()
	netDlgt.OnCreateSession(&c)

	go c.start(conn)

	return c.Ses, err
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
