package setup

var Setup *ServerCfg

type NsqSetup struct {
	NSQAddr       string
	NSQLookup []string
}

type ServerCfg struct {
	TcpPort int
	HttpPort  int
	Nsq 	[]NsqSetup
}

