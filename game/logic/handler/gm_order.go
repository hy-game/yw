package handler

import (
	"game/logic/comps/hero"
	"game/logic/comps/item"
	"game/logic/role"
	"game/types"
	"github.com/golang/protobuf/proto"
	"math"
	"pb"
	"strconv"
)

func handleGMOrder(task *pb.MsgGMTask, s *types.Session, r *types.Role) {
	switch task.MType {
	case pb.MsgGMTask_RoleCmd:
		handleRoleCmd(task, s, r)
	case pb.MsgGMTask_SearchCmd:
		handleSearchCmd(task, s, r)
	default:
	}
}

func handleRoleCmd(task *pb.MsgGMTask, s *types.Session, r *types.Role) {
	switch task.MCmd {
	case "setlevel":
		{
			lv, err := strconv.Atoi(task.Params)
			if err != nil {
				return
			}
			if lv <= 0 || lv > math.MaxUint8 {
				return
			}
			r.Data.Level = uint32(lv)
		}
	case "additem":
		{
			gMAddGoods(task, r)
		}
	}
}

func handleSearchCmd(task *pb.MsgGMTask, s *types.Session, r *types.Role) {
	switch task.MCmd {
	case "role":
		role.PackToManage(r)
	default:
	}
}

func gMAddGoods(task *pb.MsgGMTask, r *types.Role) {
	good := &pb.CPriceItem{}
	proto.UnmarshalText(task.Params, good)
	switch good.Mtype {
	case pb.EItemType_Hero:
		{
			for i := int32(0); i < good.Count; i++ {
				hero.Add(uint32(good.Oriname), r)
			}
		}
	case pb.EItemType_Good:
		{
			item.Prize(good, pb.ESource_System, r)
		}
	default:
	}
}
