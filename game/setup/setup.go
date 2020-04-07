package setup

var Setup ServerCfg

//ServerCfg game的配置
type ServerCfg struct {
	ID          uint32
	TcpPort     int
	TcpListenIp string
	HttpPort    int
	NSQ         string
	NSQLookup   []string
}
