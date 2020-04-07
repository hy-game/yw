package app

import (
	"com/log"
	"com/util"
	"fight/logic"
	"fight/network"
	"fight/setup"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"pb"
	"strconv"
	"time"
)

//Parser 解析参数
func Parser() *cli.App {
	app := &cli.App{
		Name:    "fight",
		Version: "1.0",
		Before:  Init,
		Action:  Action,
		After:   UnInit,
	}
	return app
}

//Init 所有初始化操作
func Init(c *cli.Context) error {
	setup.Init()
	logic.Init()

	return nil
}

//Action 服务总入口
func Action(c *cli.Context) error {
	go startHTTPServe(":" + strconv.Itoa(setup.Setup.HttpPort))
	go run()
	err := startServe()
	if err == nil {
		log.Info("server start success")
	}
	return err
}

func run() {
	tSec5 := time.NewTicker(time.Second * 5)

	tcpAddr := getLocalIp()
	heartBeat(tcpAddr)

	for {
		select {
		case <-tSec5.C:
			heartBeat(tcpAddr)
		}
	}
}

func heartBeat(tcpAddr string) {
	network.Broadcast(pb.MsgIDS2S_FtHeartBeat, &pb.MsgFtHeartBeat{
		EndPoint: tcpAddr + ":" + strconv.Itoa(setup.Setup.TcpPort),
	})
}

func getLocalIp() string {
	if setup.Setup.TcpListenIp == "nil" {
		ipNets := util.GetComputerIp()
		if len(ipNets) != 1 {
			log.Panic("本机发现多网卡，需在setup中配置TcpListenIp")
			return ""
		} else {
			return ipNets[0]
		}
	}
	return setup.Setup.TcpListenIp
}

//UnInit 退出时清理操作
func UnInit(c *cli.Context) error {
	log.Info("server shutdown")
	return nil
}

//----------------初始化相关-----------------------

//---------------------服务----------------------------
func startHTTPServe(addr string) {
	log.Infof("listen http on:%d", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Errorf("start http error:%v", err)
		os.Exit(-1)
	}
}

func startServe() error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(setup.Setup.TcpPort))
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	log.Infof("tcp on %v", listener.Addr())

	s := grpc.NewServer()
	ins := &network.Server{}
	pb.RegisterSrvServiceServer(s, ins)

	go func(){
		util.WaitExit()
		s.GracefulStop()
	}()

	return s.Serve(listener)
}
