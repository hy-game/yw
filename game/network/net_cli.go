package network

import (
	"com/util"
	"errors"
	"game/types"
	"io"
	"pb"

	"com/log"
)

var (
	errAPINotFind = errors.New("api not defined")
	errMsgParse   = errors.New("parser msg error")
)

//Server rpc service
type Server struct {
}

//SrvSrv rpc stream
func (server *Server) SrvSrv(stream pb.SrvService_SrvSrvServer) error {
	sess := types.NewSession(stream)
	recvLoop(1, sess)

	log.Debugf("start stream success:%d", sess.ID)

	return handLoop(SendChanSize, sess)
}

func recvLoop(cacheSize int, s *types.Session) {
	s.RecvChan = make(chan *pb.SrvMsg, cacheSize)

	go func() {
		defer func() {
			log.Debugf("close recv loop")
			close(s.RecvChan)
		}()
		for {
			msg, err := s.Stream.Recv()
			if err == io.EOF { //cli closed
				return
			}
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case s.RecvChan <- msg:
			case <-s.Die:
				return
			}
		}
	}()
}

func handLoop(sendChanSize int, s *types.Session) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			util.PrintStack()
		}
	}()

	s.SendChan = make(chan *pb.SrvMsg, sendChanSize)
	s.Evt = make(chan types.Evt, RoleEvtSize)
	defer func() {
		s.OnClose()
		close(s.SendChan)
		close(s.Die)
	}()
	for {
		select {
		case msg, ok := <-s.RecvChan:
			if !ok {
				return nil
			}
			//msgid/2 data
			if err := cliMsgHandler.handle(uint16(msg.GetID()), msg.Msg, s); err != nil {
				log.Warnf("%s hand msg err:%v", s.Desc(), err)
				return err
			}

		case out := <-s.SendChan:
			if err := s.Stream.Send(out); err != nil {
				log.Warnf("%s send msg err:%v", s.Desc(), err)
				return err
			}
		case evtData := <-s.Evt:
			if err := evtHandle.handle(evtData, s); err != nil {
				log.Warnf("%s hand msg err:%v", s.Desc(), err)
				return err
			}
		}

		if s.IsClose {
			return nil
		}
	}
}
