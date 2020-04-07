package main

import (
	"com/log"
	"github.com/sirupsen/logrus"
	"manage/app"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.Init("./log/", "server.log", logrus.DebugLevel)

	app := app.Parse()
	if app != nil {
		app.Run(os.Args)
	}
}
