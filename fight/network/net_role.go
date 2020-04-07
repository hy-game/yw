package network

import (
	"com/log"
	"errors"
	"fight/types"
	"google.golang.org/grpc/metadata"
	"io"
	"pb"
	"strconv"
)

var (
	errIncorrectGameMsgType = errors.New("incorrect game msg type")
	errRegionNotCreate      = errors.New("region not create")
)

//Server rpc service
type Server struct {
}

//GameSrv rpc stream
func (server *Server) SrvSrv(stream pb.SrvService_SrvSrvServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return errIncorrectGameMsgType
	}

	if len(md["guid"]) == 0 {
		log.Error("cannot read key:guid from metadata")
		return errIncorrectGameMsgType
	}
	roleid, err := strconv.Atoi(md["guid"][0])
	if err != nil {
		log.Warnf("recv guid %v err:%v", md["guid"][0], err)
		return errIncorrectGameMsgType
	}
	guid := uint32(roleid)

	if len(md["region"]) == 0 {
		log.Error("cannot read key:region from metadata")
		return errIncorrectGameMsgType
	}
	rgId, err := strconv.Atoi(md["region"][0])
	if err != nil {
		log.Warnf("recv region %v err:%v", md["region"][0], err)
		return errIncorrectGameMsgType
	}

	if len(md["reconn"]) == 0 {
		log.Error("cannot read key:region from metadata")
		return errIncorrectGameMsgType
	}
	reconn, err := strconv.Atoi(md["reconn"][0])
	if err != nil {
		log.Warnf("recv region %v err:%v", md["reconn"][0], err)
		return errIncorrectGameMsgType
	}

	ses := types.NewSession(stream)
	rg, ok := MgrForRegion.roleEnter(uint32(rgId), guid, ses, reconn > 0)
	if !ok || rg == nil {
		log.Warnf("region %d not create", rgId)
		return errRegionNotCreate
	}
	defer MgrForRegion.roleLeave(uint32(rgId), guid)

	for {
		msg, err := stream.Recv()
		if err == io.EOF { //cli closed
			return err
		}
		if err != nil {
			log.Error(err)
			return err
		}
		select {
		case rg <- RoleMsg{
			RoleID: guid,
			Msg:    msg,
		}:
		case <-ses.Die:
			return nil
		}
	}
}
