package setup

var Setup ServerCfg

type ServerCfg struct {
	HttpPort  int
	NSQ       string
	NSQLookup []string
}
