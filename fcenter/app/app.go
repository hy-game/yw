package app

import (
	"com/util"
	"fcenter/logic"
	"fcenter/network"
	cli "gopkg.in/urfave/cli.v2"
	"com/log"
)

// Parser :
func Parser() *cli.App {
	app := &cli.App{Name: "fight", Version: "1.0", Before: appInit, Action: appAction, After: appUnInit}
	return app
}

func appInit(c *cli.Context) error {
	err := initSetup()
	if err != nil {
		return nil
	}
	network.NsqInit(gAppSetup.NSQ, gAppSetup.NSQLookup)
	logic.Init()
	return nil
}

func appAction(c *cli.Context) error {
	network.NsqListen()

	log.Infof("start success. httpPort:%d")
	util.WaitExit()
	return nil
}

func appUnInit(c *cli.Context) error {
	return nil
}
