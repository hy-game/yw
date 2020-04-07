package enet

import (
	"github.com/golang/protobuf/proto"
	"math"
)

type NetDlgt struct {
	Handler *Route
	Cfg     *Config
}

//OnCreateListenSession  创建session需要的代理，非线程安全
func (n *NetDlgt) OnCreateSession(c *Connection) {
	c.Ses = &Session{Conn: c}
	c.Handler = n.Handler
	c.Pkg = &Pkg{}
}

func (n *NetDlgt) GetCfg() *Config {
	return n.Cfg
}

func (n *NetDlgt) MsgRegistry(msgID uint16, cf func() proto.Message, df func(msg proto.Message, s ISession)) {
	n.Handler.Register(msgID, cf, df)
}

var DefaultNetListen = &NetDlgt{
	Handler: NewRoute(math.MaxInt16),
	Cfg: &Config{
		ReadDeadline:        0,
		OutChanSize:         128,
		RpmLimit:            0,
		EvtChanSize:         128,
		ReadSocketBuffSize:  32767,
		WriteSocketBuffSize: 32767,
		RecvPackegLimit:     32767,
		SendPacketLimit:     32767,
	},
}
var DefaultNetConn = &NetDlgt{
	Handler: NewRoute(math.MaxInt16),
	Cfg: &Config{
		ReadDeadline: 0,
		OutChanSize:  128,
		RpmLimit:     0,

		ReadSocketBuffSize:  32767,
		WriteSocketBuffSize: 32767,
		RecvPackegLimit:     32767,
		SendPacketLimit:     32767,
	},
}
