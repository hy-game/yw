package gnet

import (
	"com/util"
	"context"
	"crypto/rc4"
	"gate/service"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"pb"
	"strconv"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/metadata"

	log "github.com/sirupsen/logrus"
)

var (
	sesID           = uint32(1)
	sendPacketLimit = 1024 * 1024
)

const (
	packetHeadLen = 4
)

const (
	SesClose   = 0x00000001
	FightClose = 0x00000010
)

//Session 客户端和gate的网络会话
type Session struct {
	Id   uint32
	Conn net.Conn

	In       chan []byte
	out      chan *PkgWriter
	evt      chan Evt
	MQ       chan pb.SrvMsg
	outCache []byte

	StreamGm pb.SrvService_SrvSrvClient
	StreamFt pb.SrvService_SrvSrvClient

	ConnTime    time.Time
	LastPkgTime time.Time
	PkgCnt      uint32
	PkgCnt1Min  int
	Ip          net.IP

	Ctrl     chan struct{}
	fightDie chan struct{}
	flag     util.Flag
	EnCoder  *rc4.Cipher
	DeCoder  *rc4.Cipher
}

func (s *Session) String() string {
	return strconv.Itoa(int(s.ID())) + "_" + s.Ip.String()
}
func (s *Session) ID() uint32 {
	return s.Id
}

func (s *Session) OnConnect() {
	s.Id = atomic.AddUint32(&sesID, 1)
	AddCliSession(s.Id, s)

	log.Infof("%s connect", s.String())
}

func (s *Session) OnClosed() {
	log.Infof("%s disconnect", s.String())
	RemoveCliSession(s.Id)
}

//close 关闭,非线程安全,只能在消息里调用
func (s *Session) Close() {
	s.flag.Add(SesClose)
}

//closeToFt 关闭到fight的网络会话
func (s *Session) CloseToFt() {
	s.flag.Add(FightClose)
}

//start recv loop
func (s *Session) start(conn net.Conn, cfg *Config) {
	s.evt = make(chan Evt, 100)
	s.out = make(chan *PkgWriter, cfg.OutChanSize)

	s.outCache = make([]byte, sendPacketLimit+6)
	go s.sendLoop(cfg)

	waitGroup.Add(1)
	go s.mainLoop(cfg)

	s.recvLoop(cfg)
}

//main
func (s *Session) mainLoop(cfg *Config) {
	log.Debug("main loop start")
	defer func() {
		waitGroup.Done()
		log.Debugf("main loop stop %v", waitGroup)

		if err := recover(); err != nil {
			log.Error(err)
			util.PrintStack()
		}
	}()

	s.ConnTime = time.Now()
	s.LastPkgTime = s.ConnTime

	s.MQ = make(chan pb.SrvMsg, cfg.EvtChanSize)
	tick := time.NewTicker(time.Minute)

	defer func() {
		s.OnClosed()
		if s.StreamGm != nil {
			s.StreamGm.CloseSend()
		}
		if s.StreamFt != nil {
			s.StreamFt.CloseSend()
			close(s.fightDie)
		}
		close(s.Ctrl)
	}()

	s.OnConnect()

	for {
		select {
		case cliMsg, ok := <-s.In:
			if !ok {
				return
			}

			s.PkgCnt++
			s.PkgCnt1Min++
			s.LastPkgTime = time.Now()

			s.onRecvCliMsg(cliMsg)
		case gmMsg := <-s.MQ:
			s.onRecvGameMsg(gmMsg)
		case e := <-s.evt:
			s.onEvent(e)
		case <-tick.C:
			s.check1Min(cfg)
		case <-shutDown:
			s.Close()
		}

		if s.flag.Has(FightClose) {
			if s.StreamFt != nil {
				s.StreamFt.CloseSend()
				close(s.fightDie)
			}
			s.StreamFt = nil
		}

		if s.flag.Has(SesClose) {
			return
		}
	}
}

func (s *Session) check1Min(cfg *Config) {
	defer func() {
		s.PkgCnt1Min = 0
	}()

	if cfg.RpmLimit > 0 && s.PkgCnt1Min > cfg.RpmLimit {
		s.Close()

		log.WithFields(log.Fields{
			"id":      s.String(),
			"cnt1min": s.PkgCnt1Min,
			"total":   s.PkgCnt,
		}).Error("RPM")
	}
}

//forwardToGame	转发数据到game
func (s *Session) forwardToGame(msgId uint16, msgData []byte) {
	msg := &pb.SrvMsg{
		ID:  uint32(msgId),
		Msg: msgData,
	}
	if s.StreamGm == nil {
		log.Errorf("stream no open")
		return
	}

	if err := s.StreamGm.Send(msg); err != nil {
		log.Errorf("forward to game:%v", err)
		return
	}
	log.Debugf("%s forward to game:%d", s.String(), msgId)
}

func (s *Session) forwardToFight(msgID uint16, msgData []byte) {
	msg := &pb.SrvMsg{ID: uint32(msgID), Msg: msgData}
	if s.StreamFt == nil {
		log.Errorf("stream no open")
		return
	}

	if err := s.StreamFt.Send(msg); err != nil {
		log.Errorf("forward to fight:%v", err)
		return
	}
	log.Debugf("%s forward to fight:%d", s.String(), msgID)
}

//startStreamGm 开启到game的流
func (s *Session) startStreamGm(gameID uint16, acc string) {
	// 连接到已选定game服务器
	conn := service.Get("game", gameID)
	if conn == nil {
		log.Errorf("cannot get game service, id:%d", gameID)
		return
	}

	cli := pb.NewSrvServiceClient(conn)
	// 开启到游戏服的流
	mtdata := metadata.New(map[string]string{"acc": acc})
	ctx := metadata.NewOutgoingContext(context.Background(), mtdata)
	stream, err := cli.SrvSrv(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	s.StreamGm = stream
	log.Debugf("%s start to game%d stream success", acc, gameID)
	// 读取GAME返回消息的goroutine
	go func(s *Session) {
		for {
			in, err := s.StreamGm.Recv()
			if err == io.EOF { // 流关闭
				log.Debug(err)
				return
			}
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case s.MQ <- *in:
			case <-s.Ctrl:
				return
			}
		}
	}(s)
}

func (s *Session) startStreamFt(roleGuid int, ftID uint16, btGuid int, isReConn string) {
	log.Debugf("%d start stream ft %d", roleGuid, ftID)

	conn := service.Get("fight", ftID)
	if conn == nil {
		log.Errorf("cannot get fight service, id:%d", ftID)
		return
	}

	cli := pb.NewSrvServiceClient(conn)
	// 开启到fight服的流
	mtdata := metadata.New(map[string]string{"guid": strconv.Itoa(roleGuid), "region": strconv.Itoa(btGuid), "reconn": isReConn})
	ctx := metadata.NewOutgoingContext(context.Background(), mtdata)
	stream, err := cli.SrvSrv(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	s.flag.Del(FightClose)
	s.fightDie = make(chan struct{})
	s.StreamFt = stream

	// 读取fight返回消息的goroutine
	go func(sess *Session) {
		for {
			in, err := sess.StreamFt.Recv()
			if err == io.EOF { // 流关闭
				log.Debug(err)
				return
			}
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case sess.MQ <- *in:
			case <-sess.fightDie:
				return
			}
		}
	}(s)
}

func (s *Session) ReConn(msg *pb.MsgReConn) {
	s.startStreamGm(uint16(msg.GameID), msg.Acc)
	b, err := proto.Marshal(msg)
	if err != nil {
		s.Close()
	}
	s.forwardToGame(uint16(pb.MsgIDC2S_C2SReConn), b)

	if msg.FightID != 0 && msg.BattleGuid != 0 {
		s.startStreamFt(int(msg.Guid), uint16(msg.FightID), int(msg.BattleGuid), "1")
	}
}

func (s *Session) OnEnterBattle(msg *pb.MsgBattleEnterReq) {
	if msg.FtID > 0 {
		if s.StreamFt == nil {
			s.startStreamFt(int(msg.Guid), uint16(msg.FtID), int(msg.BattleGuid), "1")
			if s.StreamFt == nil {
				sendMsg := &pb.MsgBattleEnterAck{RetCode: 2} //无法连接战斗服务器
				s.SendPB(uint16(pb.MsgIDS2C_BattleEnterAck), sendMsg)
				return
			}
		}
		b, _ := proto.Marshal(msg)
		s.forwardToFight(uint16(pb.MsgIDC2S_BattleEnterReqToFs), b)
	} else {
		if s.StreamGm == nil {
			sendMsg := &pb.MsgBattleEnterAck{RetCode: 3} //无法连接游戏服务器
			s.SendPB(uint16(pb.MsgIDS2C_BattleEnterAck), sendMsg)
			return
		}
		b, _ := proto.Marshal(msg)
		s.forwardToGame(uint16(pb.MsgIDC2S_BattleEnterReqToGs), b)
	}
}

func (s *Session) OnLeaveBattle(msg *pb.MsgBattleLeaveReq) {
	if msg.FtID > 0 {
		if s.StreamFt == nil {
			s.startStreamFt(int(msg.Guid), uint16(msg.FtID), int(msg.BattleGuid), "1")
			if s.StreamFt == nil {
				sendMsg := &pb.MsgBattleLeaveAck{RetCode: 2} //无法连接战斗服务器
				s.SendPB(uint16(pb.MsgIDS2C_BattleLeaveAck), sendMsg)
				return
			}
		}
		b, _ := proto.Marshal(msg)
		s.forwardToFight(uint16(pb.MsgIDC2S_BattleLeaveReqToFs), b)
	} else {
		if s.StreamGm == nil {
			sendMsg := &pb.MsgBattleLeaveAck{RetCode: 3} //无法连接游戏服务器
			s.SendPB(uint16(pb.MsgIDS2C_BattleLeaveAck), sendMsg)
			return
		}
		b, _ := proto.Marshal(msg)
		s.forwardToGame(uint16(pb.MsgIDC2S_BattleLeaveReqToGs), b)
	}
}
