package network

import (
	"fight/types"
	"pb"
	"sync"
)

var MgrForRegion = newRegionMgr()

type RegionMgr struct {
	Regions map[uint32]*RgnForm
	mtx     sync.Mutex
}

func newRegionMgr() *RegionMgr {
	m := &RegionMgr{
		Regions: make(map[uint32]*RgnForm),
	}
	return m
}

//Create	创建一个场景框架
func (m *RegionMgr) Create(msg *pb.MsgBattleStartData) *RgnForm {
	rf := NewRgnForm(msg.BattleGuid, msg)
	//todo 根据msg中配置id找到场景类型
	//switch  {
	//case pb.RegionType_RT_Normal:
	//	rf.region = &Region{}
	//}
	rf.Battle = &Region{}

	return rf
}

//Add	添加这个场景框架到管理器
func (m *RegionMgr) Add(rg *RgnForm) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.Regions[rg.Id] = rg
}

//Del	删除
func (m *RegionMgr) Del(rgID uint32) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.Regions, rgID)
}

func (m *RegionMgr) roleEnter(rgid uint32, roleID uint32, ses *types.Session, isReConn bool) (mq chan RoleMsg, ok bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	rg, ok := m.Regions[rgid]
	if !ok {
		if isReConn {
			ses.Send(pb.MsgIDS2C_S2CReConnFtAck, &pb.MsgReConnFtAck{
				Ret: pb.LoginCode_LCCanNotReConn,
			})
		}
		return nil, false
	}

	rg.bindRoleSession(roleID, ses)

	return rg.MsgQ, true
}

func (m *RegionMgr) roleLeave(rgid uint32, roleID uint32) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	rg := m.Regions[rgid]
	if rg != nil {
		rg.bindRoleSession(roleID, nil)
	}
	return
}
