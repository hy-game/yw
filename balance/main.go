package main

import (
	"balance/app"
	"com/log"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.Init("./log/", "center.log", logrus.TraceLevel)

	app := app.Parse()
	if app != nil {
		app.Run(os.Args)
	}
}
