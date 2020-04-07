package main

import (
	"com/log"
	"game/app"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.Init("./log/", "game.log", logrus.TraceLevel)

	app := app.Parser()
	app.Run(os.Args)
}
