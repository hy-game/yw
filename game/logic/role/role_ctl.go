package role

import (
	"com/log"
	"game/gmdb"
	"game/logic/comps/fb"
	"game/logic/comps/hero"
	"game/logic/comps/item"
	"game/logic/comps/mail"
	"game/network"
	"game/setup"
	"game/types"
	"github.com/golang/protobuf/proto"
	"pb"
)

//NewRole	新建一个角色
func NewRole(data *pb.MsgPlayerData) *types.Role {
	r := &types.Role{
		Guid:  data.Guid,
		Acc:   data.Account,
		Data:  data,
		Comps: make(map[types.TypeComp]types.IComp),
	}
	createComps(r)

	return r
}

func createComps(r *types.Role) {
	for i := types.TCHero; i < types.TCMax; i++ {
		switch i {
		case types.TCHero:
			r.Comps[i] = hero.NewData(r)
		case types.TCFB:
			r.Comps[i] = fb.NewData(r)
		case types.TCItem:
			r.Comps[i] = item.NewData(r)
		case types.TCMail:
			r.Comps[i] = mail.NewData(r)
		}
	}
}

func testData(r *types.Role) {
	item.AddJinBi(9999, pb.ESource_System, r)
	item.Add(4, 888, pb.ESource_System, r)
	hero.Add(101, r)
	hero.Add(103, r)
	hero.Add(104, r)
}
func PackToManage(r *types.Role) {
	network.SendToManage(pb.MsgIDS2S_Gm2MaRoleData, r.Data)
}

func makeMsgForCli(r *types.Role) *pb.MsgPlayerData {
	msg := *r.Data
	for _, v := range r.Comps {
		if iChgComp, ok := v.(types.ICompChgData); ok {
			iChgComp.ChangeForCli(&msg, r)
		}
	}

	return &msg
}

func save(r *types.Role) {
	b, err := proto.Marshal(r.Data)
	if err != nil {
		log.Errorf("marshal role data err:%v", err)
		return
	}
	//基础数据
	t := &DbSaveRole{Guid:  r.Guid,
		Acc:   r.Acc,
		Name:  r.Data.Name,
		Level: r.Data.Level,
		Data:  b,}
	gmdb.DBRole.Write(t)
}

//Online 角色上线处理
func Online(r *types.Role, ses *types.Session) {
	if r == nil || ses == nil {
		return
	}
	types.BindRoleSess(r, ses)

	msgSend := &pb.MsgLoginForCli{}
	msgSend.Player = makeMsgForCli(r)
	r.Send(pb.MsgIDS2C_S2CLoginAck, msgSend)

	testData(r)

	log.Infof("[%d]%s online", r.Guid, r.Acc)
	Mgr4OfflineOp.Load(r.Guid)

	network.SendToCt(pb.MsgIDS2S_Gm2CtLogin, &pb.MsgKeyValueU{
		Key:   r.Guid,
		Value: setup.Setup.ID,
	})
}

func ReConnSuccess(r *types.Role, ses *types.Session) {
	msgSend := &pb.MsgLoginForCli{}

	if r == nil {
		msgSend.Ret = pb.LoginCode_LCCanNotReConn
		ses.Send(pb.MsgIDS2C_S2CReConnAck, msgSend)
		ses.Close()
		return
	}

	types.BindRoleSess(r, ses)

	msgSend.Ret = pb.LoginCode_LCSuccess
	msgSend.Player = makeMsgForCli(r)
	r.Send(pb.MsgIDS2C_S2CReConnAck, msgSend)

	network.SendToCt(pb.MsgIDS2S_Gm2CtLogin, &pb.MsgKeyValueU{
		Key:   r.Guid,
		Value: setup.Setup.ID,
	})
}

//OnOffline 角色下线的处理
func OnOffline(r *types.Role) {
	save(r)

	Mgr4Role.RoleOffline(r.Acc)

	network.SendToCt(pb.MsgIDS2S_Gm2CtOffline, &pb.MsgKeyValueU{
		Key:   r.Guid,
		Value: setup.Setup.ID,
	})

	types.UnBindRoleSess(r, r.Ses)
	log.Infof("[%d]%s offline", r.Guid, r.Acc)
}

//SecLoop	每秒循环调用
func SecLoop(r *types.Role) {
	for _, v := range r.Comps {
		if iSec, ok := v.(types.ICompSecLoop); ok {
			iSec.SecLoop(r)
		}
	}
}

//DataReset	每日数据重置
func DataReset(r *types.Role) {
	for _, v := range r.Comps {
		if iSec, ok := v.(types.ICompDataReset); ok {
			iSec.OnDataReset(r)
		}
	}
}

func CreateBattle(regionCfg *pb.MsgRegionList, r *types.Role, regionID uint32, msgFighter []uint32) {
	//define msg
	msg := &pb.MsgBattleStartData{BattleID: regionID}
	//fill msgRole
	msgRole := &pb.MsgRoleInFight{}
	msgRole.SesID = r.Ses.ID
	msgRole.Acc = r.Acc
	msgRole.Guid = r.Guid
	msgRole.ServerID = setup.Setup.ID
	msg.Roles = append(msg.Roles, msgRole)
	//build fighter
	FighterGuid := (uint32)(0)
	for _, v := range msgFighter {
		FighterGuid++
		fighter := hero.MakeFighter(v, r)
		if fighter == nil {
			continue
		}
		fighter.Guid = FighterGuid
		fighter.MaxHp = fighter.Attrs[uint32(pb.AttrType_InitHp)]
		fighter.Hp = fighter.MaxHp
		msgRole.Fighters = append(msgRole.Fighters, fighter)

		//todo delete 琼元说先填到这里
		fighter.Attrs[pb.AttrType_MoveSpeed] *= 100
		fighter.Speed = 800
	}
	if len(msgRole.Fighters) > 0 {
		msgRole.Fighters[0].Visible = true
	}
	//sen msg
	if regionCfg.CheckType == 0 { // do nothing
		r.Battle = msg //记录战斗简要信息
		msgSend := &pb.MsgBattleCreateAck{BtStart: msg}
		r.Send(pb.MsgIDS2C_BattleCreateAck, msgSend)
	} else {
		network.SendToFc(pb.MsgIDS2S_Gm2FcBattleCreateReq, msg)
	}
}

//CanEnterBattle	战斗是否可以进入
func CanEnterBattle(regionCfg *pb.MsgRegionList, r *types.Role) (pb.EFormationType, bool) {
	switch regionCfg.Type {
	case pb.EFBType_FBTNormal:
		return  fb.CanEnter(regionCfg, r)
	}
	return pb.EFormationType_EFBNormal, false
}

func FinishBattle(r *types.Role, data *pb.MsgBattleFinishData) {
	regionCfg := fb.CfgRegion(data.BattleID)
	if regionCfg == nil {
		return
	}
	switch regionCfg.Type {
	case pb.EFBType_FBTNormal:
		fb.OnFinish(r, data)
	}
	r.Send(pb.MsgIDS2C_BattleFinish, &pb.MsgBattleFinish{BtFinish: data})
}
