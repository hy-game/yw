package setup

var Setup *ServerCfg

type ServerCfg struct {
	HttpPort  uint16
	NSQ       string
	NSQLookup []string
}