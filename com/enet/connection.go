package enet

import (
	"com/evt"
	"com/util"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	packetHeadLen = 4
)

//Session 客户端和gate的网络会话
type Connection struct {
	Ses     ISession
	Pkg     IPkg
	Handler *Route
	Cfg     *Config

	Ctrl chan struct{}
	Conn net.Conn
	Ip   net.IP
	In   chan []byte
	out  chan IPkgWriter

	ConnTime    time.Time
	LastPkgTime time.Time
	PkgCnt      uint32
	PkgCnt1Min  int
	IsClosed    bool

	MQ   chan evt.Evt
	wait sync.WaitGroup
}

//Close 关闭,非线程安全,只能在消息里调用
func (s *Connection) Close() {
	s.IsClosed = true
}

func (s *Connection) IsClose() bool {
	return s.IsClosed
}

//start recv loop
func (s *Connection) start(conn net.Conn) {
	defer util.PanicPrintStack()
	defer conn.Close()

	s.Conn = conn
	s.In = make(chan []byte, s.Cfg.InChanSize) //no cache
	defer func() {
		close(s.In)
		//		log.Debug("recv loop stop")
	}()
	s.Ctrl = make(chan struct{})

	s.wait.Add(1)
	go s.sendLoop()

	waitGroup.Add(1)
	go s.mainLoop()

	s.recvLoop()
}

//main
func (s *Connection) mainLoop() {
	//	log.Debug("main loop start")
	defer func() {
		waitGroup.Done()
		//		log.Debugf("main loop stop %v", waitGroup)
		util.PanicPrintStack()
	}()

	s.ConnTime = time.Now()
	s.LastPkgTime = s.ConnTime

	s.MQ = make(chan evt.Evt, s.Cfg.EvtChanSize)
	tick := time.NewTicker(time.Minute)

	defer func() {
		s.Ses.OnClosed()
		close(s.Ctrl)
	}()

	s.wait.Wait()
	s.Ses.OnConnect()

	for {
		select {
		case cliMsg, ok := <-s.In:
			if !ok {
				return
			}

			s.PkgCnt++
			s.PkgCnt1Min++
			s.LastPkgTime = time.Now()

			s.OnRecv(cliMsg)
		case evt := <-s.MQ:
			s.Ses.OnEvent(evt)
		case <-tick.C:
			s.check1Min()
		case <-shutDown:
			s.Close()
		}

		if s.IsClose() {
			return
		}
	}
}

//recv
func (s *Connection) recvLoop() {
	//	log.Debug("recv loop start")
	header := make([]byte, packetHeadLen)

	for {
		if s.Cfg.ReadDeadline > 0 {
			s.Conn.SetReadDeadline(time.Now().Add(time.Second * s.Cfg.ReadDeadline))
		}
		n, err := io.ReadFull(s.Conn, header)
		if err != nil {
			netErr, _ := err.(*net.OpError)
			log.Warnf("read header failed :%s, err:%v, %v, size:%d", s.Ses.String(), err, netErr, n)
			return
		}
		//header是数据的大小，不包括header本身
		size := binary.BigEndian.Uint32(header)
		if size > s.Cfg.RecvPackegLimit {
			log.Warnf("read data len out of limit :%s, size%v", s.Ses.String(), size)
			return
		}
		data := make([]byte, size)
		n, err = io.ReadFull(s.Conn, data)
		if err != nil {
			log.Warnf("read data failed :%s, err:%v, size:%d", s.Ses.String(), err, n)
			return
		}

		select {
		case s.In <- data:
		case <-s.Ctrl:
			log.Warnf("connection close by logic :%s", s.Ses.String())
			return
		}
	}
}

//send
func (s *Connection) sendLoop() {
	//	log.Debug("send loop start")

	defer func() {
		util.PanicPrintStack()
		//		log.Debug("send loop stop")
	}()

	s.out = make(chan IPkgWriter, s.Cfg.OutChanSize)

	outCache := make([]byte, s.Cfg.SendPacketLimit+4)
	s.wait.Done()

	for {
		select {
		case p := <-s.out:
			s.rawSend(p, outCache)
		case <-s.Ctrl:
			return
		}
	}
}

func (s *Connection) rawSend(p IPkgWriter, cache []byte) {
	len := p.Write(cache)
	_, err := s.Conn.Write(cache[:len])
	if err != nil {
		log.Warnf("send data error, err:%v", err)
	}
}

func (s *Connection) check1Min() {
	defer func() {
		s.PkgCnt1Min = 0
	}()

	if s.Cfg.RpmLimit > 0 && s.PkgCnt1Min > s.Cfg.RpmLimit {
		s.Close()

		log.WithFields(log.Fields{
			"id":      s.Ses.String(),
			"cnt1min": s.PkgCnt1Min,
			"total":   s.PkgCnt,
		}).Error("RPM")
	}
}
