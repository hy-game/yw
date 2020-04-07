package main

import (
	"account/app"
	"com/log"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.Init("./log/", "server.log", logrus.TraceLevel)

	app := app.Parse()
	if app != nil {
		app.Run(os.Args)
	}
}
