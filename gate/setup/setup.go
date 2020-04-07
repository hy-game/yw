package setup

var Setup *ServerCfg

type ServerCfg struct {
	Id         uint32
	NetAddr		string
	ListenPort int
	HttpPort   int
	NSQ        string
	NSQLookup  []string
}
