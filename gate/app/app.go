package app

import (
	"com/log"
	"com/util"
	"gate/gnet"
	"gate/logic"
	"gate/service"
	"gate/setup"
	"gopkg.in/urfave/cli.v2"
	"net/http"
	_ "net/http/pprof"
	"pb"
	"strconv"
	"time"
)

var (
	ctrl = make(chan struct{})
)

func Parse() *cli.App {
	app := &cli.App{
		Name:    "gate",
		Usage:   "a gateway for games",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "services",
				Value: cli.NewStringSlice("gate"),
				Usage: "need service name",
			},
		},
		Before: Init,
		Action: Action,
		After:  UnInit,
	}
	return app
}

func Init(c *cli.Context) error {
	log.Info("server init")
	setup.Setup = &setup.ServerCfg{}
	err := util.ReadJson(setup.Setup, "./setup.json")
	if err != nil {
		log.Panicf("read setup err:%v", err)
		return err
	}

	service.Init("game", "fight")
	logic.Init()

	return nil
}

func UnInit(c *cli.Context) error {
	service.Clear()

	return nil
}

func Action(c *cli.Context) error {
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(setup.Setup.HttpPort), nil)
		if err != nil {
			log.Panicf("start http serve err:%v", err)
		}
	}()

	go loop()

	cfg := &gnet.Config{}
	err := util.ReadJson(cfg, "./config.json")
	if err != nil {
		log.Panicf("read config err:%v", err)
		return err
	}
	go gnet.StartTCPServer("0.0.0.0:"+strconv.Itoa(setup.Setup.ListenPort), cfg)

	log.Infof("start success tcp on:%d, http on:%d", setup.Setup.ListenPort, setup.Setup.HttpPort)

	util.WaitExit()
	close(ctrl)
	gnet.Close()
	log.Info("server closed")
	return nil
}

func heartBeat(tcpAddr string) {
	logic.SendToBa(pb.MsgIDS2S_GtHeartBeat, &pb.MsgGtHeartBeat{
		EndPoint: tcpAddr + ":" + strconv.Itoa(setup.Setup.ListenPort),
		RoleCnt:  gnet.GetCliSessCnt(),
	})
}

func loop() {
	tSec5 := time.NewTicker(time.Second * 5)

	heartBeat(setup.Setup.NetAddr)	//外网地址

	for {
		select {
		case <-tSec5.C:
			heartBeat(setup.Setup.NetAddr)
		}
	}
}
