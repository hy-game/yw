package util

import (
	"com/log"
	"os"
	"os/signal"
	"syscall"
)

func WaitExit() {
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range exitChan {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Info("Signal: %v server closing ...", s)
			return
		}
	}
}
