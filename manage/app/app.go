package app

import (
	"com/log"
	"com/util"
	"gopkg.in/urfave/cli.v2"
	"manage/logic"
	"manage/logic/configs"
	"manage/setup"
	_ "manage/web/routes"
	"net/http"
)

func Init(c *cli.Context) error {
	log.Info("server init")

	err := util.ReadJson(&setup.Setup, "./setup.json")
	if err != nil {
		log.Panicf("load setup err:%v", err)
		return err
	}

	logic.Init()
	configs.Reload(false)

	return nil
}

func Action(c *cli.Context) error {
	go http.ListenAndServe(":"+util.ToString(setup.Setup.HttpPort), nil)

	go WriteIPC()
	log.Infof("start success, listen http on:%d", setup.Setup.HttpPort)
	util.WaitExit()

	return nil
}

func UnInit(c *cli.Context) error{
	log.Info("closing...")
	log.Info("server exit")
	return nil
}

func Parse() *cli.App {
	app := &cli.App{
		Name:    "gate",
		Usage:   "a gateway for games",
		Version: "1.0",
		Before:  Init,
		Action:  Action,
		After:UnInit,
	}
	return app
}
