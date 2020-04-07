package controllers

import (
	"bytes"
	"com/util"
	"fmt"
	"github.com/golang/protobuf/proto"
	"html/template"
	"log"
	"manage/logic/configs"
	"net/http"
	"pb"
	"share"
	"strconv"
	"strings"
	"time"
)

type Web struct {
	// 结构体内部可以记录日志
	// 或者用户验证
	Param interface{}
}

var WebMgr Web

//默认
func (*Web) SayHello(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/view/main.html", "web/view/header.html")
	if err != nil {
		log.Fatalln(err)
	}
	t.Execute(w, nil)
}

func (*Web) GmRole(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/gm_role.html", "web/view/header.html")
	t.Execute(w, nil)
}

func (*Web) GmItem(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/gm_item.html", "web/view/header.html")
	t.Execute(w, nil)
}

func (*Web) GmSystem(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/gm_system.html", "web/view/header.html")
	t.Execute(w, nil)
}

func (*Web) GmList(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/gm_list.html")
	t.Execute(w, nil)
}

func (*Web) SearchList(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/search_list.html")
	t.Execute(w, nil)
}

func (*Web) SearchRole(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/search_role.html", "web/view/header.html")
	t.Execute(w, nil)
}

func (*Web) SearchRL(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/view/components/search_rl.html", "web/view/header.html")
	t.Execute(w, nil)
}

func (*Web) GmHRole(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprint(w, "发送完成")
	cmd := r.FormValue("command")
	pam := r.FormValue("gmParam")
	ro := r.FormValue("gmPlayer")
	pl, err := strconv.Atoi(ro)
	if err != nil {
		return
	}

	task := &pb.MsgGMTask{
		MType:     pb.MsgGMTask_RoleCmd,
		MCmd:      cmd,
		MPlayerId: int32(pl),
		Params:    pam,
		Server:    share.GameTopic,
	}
	sendGMTask(task)
}

func (*Web) GmAlliance(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprint(w, "发送完成")
	cmd := r.FormValue("cmdally")
	pam := r.FormValue("prmally")
	ro := r.FormValue("gmPlayer")
	pl, err := strconv.Atoi(ro)
	if err != nil {
		return
	}

	task := &pb.MsgGMTask{
		MType:     pb.MsgGMTask_AllyCmd,
		MCmd:      cmd,
		MPlayerId: int32(pl),
		Params:    pam,
		Server:    share.CenterTopic,
	}
	sendGMTask(task)
}

func (*Web) ItemCfg(w http.ResponseWriter, r *http.Request) {
	svr := r.FormValue("chan")
	cfg_key := r.FormValue("cfg_key")

	cfgs := configs.GetCfgs(svr)
	buff := bytes.NewBuffer(nil)
	if cfg_key == "GoodsList" {
		buff.WriteString(`<select id="gmParam" name="gmParam" class="chosen_sel form-control">`)
		for _, cfg := range cfgs.TableCfgs {
			if cfg.Filename == cfg_key {
				for _, v := range cfg.Lines {
					strs := strings.FieldsFunc(v, util.SplitRule)
					if len(strs) > 3 {
						buff.WriteString(fmt.Sprintf(`<option value="%v">%v %v</option>`, strs[1], strs[1], strs[2]))
					}
				}
				break
			}
		}
		buff.WriteString(`</select>`)
	} else if cfg_key == "HeroAttri" {
		buff.WriteString(`<select id="gmParam" name="gmParam" class="chosen_sel form-control">`)
		for _, cfg := range cfgs.TableCfgs {
			if cfg.Filename == cfg_key {
				for _, v := range cfg.Lines {
					strs := strings.FieldsFunc(v, util.SplitRule)
					if len(strs) > 4 {
						buff.WriteString(fmt.Sprintf(`<option value="%v">%v %v</option>`, strs[1], strs[1], strs[3]))
					}
				}
				break
			}
		}
		buff.WriteString(`</select>`)
	} else {
		buff.WriteString(`<input type="text" class="form-control" id="gmParam" name="gmParam" placeholder="装备配置id" />`)
	}

	fmt.Fprint(w, buff.String())
}

func (*Web) GmHSystem(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprint(w, "发送完成")
	cmd := r.FormValue("command")
	pam := r.FormValue("gmParam")

	switch cmd {
	case "rcfgs":
		fallthrough
	case "ryyact":
		fallthrough
	case "rrcfgs":
		fallthrough
	case "rryyact":
		{
			task := &pb.MsgGMTask{
				MType:  pb.MsgGMTask_ReloadCmd,
				MCmd:   cmd,
				Params: pam,
				Server: share.ManageTopic,
			}
			sendGMTask(task)
		}
	case "s":
		fallthrough
	case "q":
		{
			task := &pb.MsgGMTask{
				MType:  pb.MsgGMTask_SystemCmd,
				MCmd:   cmd,
				Params: pam,
				Server: share.ManageTopic,
			}
			sendGMTask(task)
		}
	default:
	}
}

func (*Web) GmHItem(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprint(w, "发送完成")
	cmd := r.FormValue("command")
	mtype := r.FormValue("mtype")
	pam := r.FormValue("gmParam")
	count := r.FormValue("mcount")
	player := r.FormValue("gmPlayer")
	item_typ, err := strconv.Atoi(mtype)
	if err != nil {
		return
	}
	item_ori, err := strconv.Atoi(pam)
	if err != nil {
		return
	}
	item_cnt, err := strconv.Atoi(count)
	if err != nil {
		return
	}
	pl, err := strconv.Atoi(player)
	if err != nil {
		return
	}
	item := &pb.CPriceItem{
		Mtype:   pb.EItemType(item_typ),
		Oriname: int32(item_ori),
		Count:   int32(item_cnt),
	}

	task := &pb.MsgGMTask{
		MType:     pb.MsgGMTask_RoleCmd,
		MCmd:      cmd,
		Params:    proto.MarshalTextString(item),
		MPlayerId: int32(pl),
		Server:    share.GameTopic,
	}
	sendGMTask(task)
}

func (this *Web) GmSRole(w http.ResponseWriter, r *http.Request) {
	skey := r.FormValue("search_key")
	pid, err := strconv.Atoi(skey)
	if err != nil {
		fmt.Fprint(w, "查询参数为空")
		return
	}
	task := &pb.MsgGMTask{
		MType:     pb.MsgGMTask_SearchCmd,
		MCmd:      "role",
		MPlayerId: int32(pid),
		Server:    share.GameTopic,
	}
	sendGMTask(task)

	ret := make(chan *pb.MsgPlayerData, 1)
	tout := time.After(10 * time.Second)
	this.Param = ret
	select {
	case data := <-ret:
		FullRoleSearch(w, data)
	case <-tout:
		close(ret)
		fmt.Fprint(w, "查询超时")
	}
}

func (this *Web) GmSRL(w http.ResponseWriter, r *http.Request) {
	skey := r.FormValue("search_key")
	pid, err := strconv.Atoi(skey)
	if err != nil {
		fmt.Fprint(w, "查询参数为空")
		return
	}

	task := &pb.MsgGMTask{
		MType:     pb.MsgGMTask_SearchCmd,
		MCmd:      "rl",
		MPlayerId: int32(pid),
		Server:    share.CenterTopic,
	}
	sendGMTask(task)

	ret := make(chan *pb.MsgRankListPack, 1)
	tout := time.After(10 * time.Second)
	this.Param = ret
	select {
	case data := <-ret:
		FullRLSearch(w, data)
	case <-tout:
		close(ret)
		fmt.Fprint(w, "查询超时")
	}
}

func FullRoleSearch(w http.ResponseWriter, data *pb.MsgPlayerData) {
	t, _ := template.ParseFiles("web/view/components/search_role_ret.html")
	t.Execute(w, data)
}

func FullRLSearch(w http.ResponseWriter, data *pb.MsgRankListPack) {
	t, _ := template.ParseFiles("web/view/components/search_rl_ret.html")
	t.Execute(w, data)
}
