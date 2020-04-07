package app

import (
	"com/database/db"
	"com/log"
	"com/util"
	"game/gmdb"
	"game/logic"
	"game/network"
	"game/setup"
	"game/types"
	"gopkg.in/urfave/cli.v2"
	"net/http"
	_ "net/http/pprof"
	"pb"
	"strconv"
	"time"
)

//Parser 解析参数
func Parser() *cli.App {
	app := &cli.App{
		Name:    "game",
		Version: "1.0",
		Before:  Init,
		Action:  Action,
		After:   UnInit,
	}
	return app
}

//Init 所有初始化操作
func Init(c *cli.Context) error {
	err := util.ReadJson(&setup.Setup, "setup.json")
	if err != nil {
		log.Panicf("read setup err:%v", err)
		return err
	}

	logic.Init()
	gmdb.Init()
	logic.InitAfterDBInit()

	return nil
}

//Action 服务总入口
func Action(c *cli.Context) error {
	go startHTTPServe()
	go loop()

	return network.StartServe()
}

func loop() {
	tSec5 := time.NewTicker(time.Second * 5)
	tSec1 := time.NewTicker(time.Second)

	tcpAddr := getLocalIp()
	heartBeat(tcpAddr)

	for {
		select {
		case <-tSec5.C:
			heartBeat(tcpAddr)
		case <-tSec1.C:
			secLoop()
		}
	}
}

func secLoop() {
	types.PostEvtToAllOnline(types.Evt{
		Data: nil,
		Type: types.SecLoop,
	})
}

func heartBeat(tcpAddr string) {
	network.Broadcast(pb.MsgIDS2S_GmHeartBeat, &pb.MsgGameHeartBeat{
		EndPoint: tcpAddr + ":" + strconv.Itoa(setup.Setup.TcpPort),
		RoleCnt:  types.GetRoleSesCnt(),
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
	db.Close()
	return nil
}

//----------------初始化相关-----------------------
//---------------------服务----------------------------
func startHTTPServe() {
	log.Infof("listen on http:%d", setup.Setup.HttpPort)
	addr := ":" + strconv.Itoa(setup.Setup.HttpPort)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panicf("start http error:%v", err)
	}
}
