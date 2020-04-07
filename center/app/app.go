package app

import (
	"center/configs"
	"center/logic"
	"center/logic/ranklist"
	"center/setup"
	"com/log"
	"com/util"
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
		return err
	}

	logic.Init()
	configs.Init()
	initDbRole()
	initDbMisc()
	ranklist.Init()
	logic.InitAfterDB()

	return nil
}

func Action(c *cli.Context) error {
	go func() {
		err := http.ListenAndServe(":"+util.ToString(setup.Setup.HttpPort), nil)
		if err != nil {
			log.Panicf("start http serve err:%v", err)
		}
	}()

	log.Infof("start success. httpPort:%d", setup.Setup.HttpPort)
	util.WaitExit()

	return nil
}

func UnInit(c *cli.Context) error {
	log.Info("closing...")
	log.Info("server exit")
	return nil
}

func Parse() *cli.App {
	app := &cli.App{
		Name:    "center server",
		Usage:   "login check",
		Version: "1.0",

		Before: Init,
		Action: Action,
		After:UnInit,
	}
	return app
}
