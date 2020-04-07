package app

import (
	"com/log"
	"com/util"
)

var gAppSetup = &appSetup{}

type appSetup struct {
	NSQ       string
	NSQLookup []string
}

func initSetup() error {
	err := util.ReadJson(gAppSetup, "setup.json")
	if err != nil {
		log.Panicf("read setup err:%v", err)
		return err
	}
	return nil
}
