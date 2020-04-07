package setup

import (
	"com/log"
	"com/util"
)

var Setup *ServerCfg

type ServerCfg struct {
	ID          uint32
	TcpPort     int
	TcpListenIp string
	HttpPort    int
	NSQ         string
	NSQLookup   []string
}

var tcpEndPoint string

//GetTcpEndPoint 得到监听的实际地址
func GetTcpEndPoint() string {
	return tcpEndPoint
}

func Init() error {
	err := util.ReadJson(&Setup, "setup.json")
	if err != nil {
		log.Panicf("read setup err:%v", err)
		return err
	}
	if Setup.TcpListenIp == "nil" {
		ipNets := util.GetComputerIp()
		if len(ipNets) != 1 {
			log.Panic("本机发现多网卡，需在setup中配置TcpListenIp")
			return nil
		} else {
			tcpEndPoint = ipNets[0]
		}
	}
	tcpEndPoint = Setup.TcpListenIp
	return nil
}
