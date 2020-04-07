package app

import (
	"com/log"
	"os"
	"strconv"
	"time"
)

func WriteIPC() {
	ipc, err := os.OpenFile("pid.txt", os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		log.Warnf("open file err:%v", err)
		return
	}
	pid := strconv.Itoa(os.Getpid())
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.C:
			ipc.Seek(0, 0)
			ipc.Write([]byte(pid))
			ipc.Sync()
		}
	}
}
