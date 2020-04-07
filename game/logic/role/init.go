package role

import "game/types"

//Init	初始化
func Init() {
	types.AddOfflineEvt(OnOffline)
}
