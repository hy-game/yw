package logic

import (
	"center/logic/ranklist"
	"com/util"
	"pb"
	"share"
)

func handleGMOrder(task *pb.MsgGMTask, serId uint16) {
	switch task.MType {
	case pb.MsgGMTask_SearchCmd:
		add_test()
		handleSearchCmd(task)
	default:
	}
}

//添加测试数据
func add_test() {
	for i := 0; i < 10; i++ {
		task := &pb.MsgRanklistHandle{
			Oper: pb.MsgRanklistHandle_Update,
			Type: pb.ERankListType_ERLT_PlayerLv,
			Data: &pb.MsgRankListData{
				Id:       uint32(1001 + i),
				Score:    uint64(i + 1),
				Name:     util.ToString(i + 1),
				Level:    uint32(i + 1),
				Protrait: 1001,
				AllyName: "工会",
			},
		}
		ranklist.PushTask(task)
	}

	task := &pb.MsgRanklistHandle{
		Oper: pb.MsgRanklistHandle_Update,
		Type: pb.ERankListType_ERLT_PlayerLv,
		Data: &pb.MsgRankListData{
			Id:       1012,
			Score:    10,
			Name:     "11",
			Level:    10,
			Protrait: 1001,
			AllyName: "工会",
		},
		Force: true,
	}
	ranklist.PushTask(task)
}

func handleSearchCmd(task *pb.MsgGMTask) {
	switch task.MCmd {
	case "rl": //排行榜
		{
			work := &pb.MsgRanklistHandle{
				Oper:   pb.MsgRanklistHandle_Pack,
				Type:   pb.ERankListType(task.MPlayerId),
				Msgid:  int32(pb.MsgIDS2S_Ct2MaRLData),
				Server: share.ManageTopic,
			}
			ranklist.PushTask(work)
		}
	}
}
