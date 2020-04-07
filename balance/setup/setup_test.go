package setup

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Setup(t *testing.T){
	s := &ServerCfg{
		TcpPort:  2001,
		HttpPort: 8082,
	}
	for i := 0; i != 2; i++ {
		nsq := NsqSetup{
			NSQAddr:   "127.0.0.1:4150",
		}
		nsq.NSQLookup = append(nsq.NSQLookup, "127.0.0.1:4161")
		nsq.NSQLookup = append(nsq.NSQLookup, "192.168.0.101:4161")
		s.Nsq = append(s.Nsq, nsq)
	}

	b, err := json.Marshal(s)
	if err != nil{
		t.Fatalf("marshal err:%v", err)
	}
	err = ioutil.WriteFile("./setup.json", b, os.ModePerm)
	if err != nil {
		t.Fatalf("write file err:%v", err)
	}
}