package app

import (
	"balance/logic"
	"balance/network"
	"balance/setup"
	"com/database/orm"
	"com/log"
	"com/util"
	"gate/gnet"
	"gopkg.in/urfave/cli.v2"
	"net/http"
	_ "net/http/pprof"
)

func Init(c *cli.Context) error {
	log.Info("server init")
	setup.Setup = &setup.ServerCfg{}
	err := util.ReadJson(setup.Setup, "./setup.json")
	if err != nil {
		log.Panicf("read setup err:%v", err)
		return nil
	}

//	initDb()
	network.InitSerNet()
	logic.Init()
	network.StartSerNet()

	return nil
}

func initDb() {
	cfg := &orm.DbCfg{
		Driver:  "mysql",
		Usr:     "root",
		Pswd:    "root",
		Host:    "127.0.0.1",
		Port:    3306,
		Db:      "center",
		Charset: "utf8",
	}
	orm.Init(cfg, 1000, 2, logic.Migrate)
	orm.Close()

}

func Action(c *cli.Context) error {
	go func() {
		err := http.ListenAndServe(":"+util.ToString(setup.Setup.HttpPort), nil)
		if err != nil {
			log.Panicf("start http serve err:%v", err)
		}
	}()

	cfg := &gnet.Config{}
	err := util.ReadJson(cfg, "./config.json")
	if err != nil {
		log.Panicf("read config err:%v", err)
		return err
	}
	go gnet.StartTCPServer("0.0.0.0:"+util.ToString(setup.Setup.TcpPort), cfg)

	log.Infof("start success tcp on:%d, http on:%d", setup.Setup.TcpPort, setup.Setup.HttpPort)

	util.WaitExit()
	gnet.Close()
	log.Info("server closed")
	return nil
}

func UnInit(c *cli.Context) error {
	log.Info("closing...")
	log.Info("server exit")
	return nil
}

func Parse() *cli.App {
	app := &cli.App{
		Name:    "account server",
		Usage:   "login check",
		Version: "1.0",

		Before: Init,
		Action: Action,
		After:UnInit,
	}
	return app
}
