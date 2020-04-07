package main

import (
	"com/log"
	"gate/app"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.Init("./log/", "gate.log", logrus.TraceLevel)

	app := app.Parse()
	if app != nil {
		app.Run(os.Args)
	}
}
