package controllers

import (
	"com/log"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"manage/logic/configs"
	"manage/network"
	"net/http"
	"pb"
	"share"
)

//GM命令在此处理
type WebGM struct {
	Web
}

var WebGmMgr WebGM

//处理GM命令的handler
func (*WebGM) Handle(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprint(w, "发送成功!")
	data, err := ioutil.ReadAll(r.Body)
	//data := r.FormValue("data")
	if err != nil {
		log.Errorf("handle gm order error: %v", err)
	}
	task := &pb.MsgGMTask{}
	err = proto.Unmarshal(data, task)
	if err != nil {
		log.Errorf("handle gm order data error: %v", err)
	}
	sendGMTask(task)
}

func sendGMTask(task *pb.MsgGMTask) {
	switch task.Server {
	case share.GameTopic:
		{
			if task.Serverid == 0 {
				network.SendToGameAll(pb.MsgIDS2S_Ma2GmGMOrder, task)
			} else {
				network.SendToGame(uint16(task.Serverid), pb.MsgIDS2S_Ma2GmGMOrder, task)
			}
		}
	case share.CenterTopic:
		{
			network.SendToCenter(pb.MsgIDS2S_Ma2CtGMOrder, task)
		}
	case share.ManageTopic:
		{
			switch task.MType {
			case pb.MsgGMTask_ReloadCmd:
				if task.MCmd == "rcfgs" {
					configs.BroadCastToServer(share.GameTopic, 0, false)
				} else if task.MCmd == "ryyact" {
					configs.BroadCastToServer(share.CenterTopic, 0, true)
					configs.BroadCastToServer(share.GameTopic, 0, true)
				} else if task.MCmd == "rrcfgs" {
					configs.Reload(false)
				} else if task.MCmd == "rryyact" {
					configs.Reload(true)
				}
			case pb.MsgGMTask_SystemCmd:
				{
					switch task.MCmd {
					case "s":
					case "q":
					default:
					}
				}
			}
		}
	default:
		log.Errorf("handle gm order server [%v] error", task.Server)
	}
}
