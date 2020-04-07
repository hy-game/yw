package network

import (
	"com/log"
	"com/mq"
	"com/util"
	"game/setup"
	"google.golang.org/grpc"
	"net"
	"os"
	"pb"
	"share"
	"strconv"
)

//Init 初始化网络，并开始接收服务器间消息
func Init() { //没用默认的init，因为需要配置
	sesSer = mq.NewNsqSession(setup.Setup.NSQ, setup.Setup.NSQLookup)
}

func StartSerNet() {
	sesSer.AddConsumer(share.GameTopic+strconv.Itoa(int(setup.Setup.ID)), "gm")
	sesSer.AddConsumer(share.BroadCastTopic, "gm"+strconv.Itoa(int(setup.Setup.ID)))
	sesSer.AddConsumer(share.GameAllTopic, "gm"+strconv.Itoa(int(setup.Setup.ID)))
}

//StartServe 开始服务监听gate转发的客户端消息
func StartServe() error {
	addr := ":" + strconv.Itoa(setup.Setup.TcpPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	log.Infof("start success. tcp on:%v ", listener.Addr())

	s := grpc.NewServer()
	ins := &Server{}
	pb.RegisterSrvServiceServer(s, ins)

	go func(){
		util.WaitExit()
		s.GracefulStop()
	}()

	return s.Serve(listener)
}
