package main

import (
	"com/log"
	"fcenter/app"
	_ "net/http/pprof"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	log.Init("./log/", "fcenter.log", logrus.TraceLevel)

	app := app.Parser()
	app.Run(os.Args)
}
