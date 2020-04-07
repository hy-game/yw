package main

import (
	"com/log"
	"fight/app"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.Init("./log/", "fight.log", logrus.TraceLevel)

	app := app.Parser()
	app.Run(os.Args)
}
