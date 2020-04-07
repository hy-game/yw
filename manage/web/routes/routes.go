package routes

import (
	log "com/log"
	"com/util"
	"manage/web/controllers"
	"net/http"
)

type HandleFnc func(http.ResponseWriter, *http.Request)

//处理异常的闭包封装函数
func logPanics(f HandleFnc) HandleFnc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Errorf("[%v] caught panic: %v", r.RemoteAddr, x)
				util.PrintStack()
			}
		}()
		f(w, r)
	}
}

func init() {
	//base
	http.HandleFunc("/", logPanics(controllers.WebMgr.SayHello))
	http.HandleFunc("/components/gm_role", logPanics(controllers.WebMgr.GmRole))
	http.HandleFunc("/components/gm_list", logPanics(controllers.WebMgr.GmList))
	http.HandleFunc("/components/gm_item", logPanics(controllers.WebMgr.GmItem))
	http.HandleFunc("/components/gm_system", logPanics(controllers.WebMgr.GmSystem))
	http.HandleFunc("/components/search_list", logPanics(controllers.WebMgr.SearchList))
	http.HandleFunc("/components/search_role", logPanics(controllers.WebMgr.SearchRole))
	http.HandleFunc("/components/search_rl", logPanics(controllers.WebMgr.SearchRL))

	http.HandleFunc("/cfgs/item_cfg", logPanics(controllers.WebMgr.ItemCfg))

	http.HandleFunc("/gm/role", logPanics(controllers.WebMgr.GmHRole))
	http.HandleFunc("/gm/alliance", logPanics(controllers.WebMgr.GmAlliance))
	http.HandleFunc("/gm/system", logPanics(controllers.WebMgr.GmHSystem))
	http.HandleFunc("/gm/item", logPanics(controllers.WebMgr.GmHItem))
	http.HandleFunc("/gm/search_role", logPanics(controllers.WebMgr.GmSRole))
	http.HandleFunc("/gm/search_rl", logPanics(controllers.WebMgr.GmSRL))
	//GM
	http.HandleFunc("/gm", logPanics(controllers.WebGmMgr.Handle))
}
