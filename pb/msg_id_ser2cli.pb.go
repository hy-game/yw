// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg_id_ser2cli.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Server向Client发送消息
type MsgIDS2C int32

const (
	MsgIDS2C_S2CNone       MsgIDS2C = 0
	MsgIDS2C_Gm2GtKickRole MsgIDS2C = 1
	MsgIDS2C_Ft2GtKickRole MsgIDS2C = 2
	// ------------以上是需要gate处理的消息------------------------
	MsgIDS2C_Gm2GtMsgIdMax      MsgIDS2C = 100
	MsgIDS2C_S2CInit            MsgIDS2C = 101
	MsgIDS2C_S2CHeartBeat       MsgIDS2C = 102
	MsgIDS2C_S2CAccInfoAck      MsgIDS2C = 103
	MsgIDS2C_S2CWaitMsgRetBegin MsgIDS2C = 200
	MsgIDS2C_S2CWaitMsgRetEnd   MsgIDS2C = 500
	MsgIDS2C_S2CLoginAck        MsgIDS2C = 501
	MsgIDS2C_S2CReConnAck       MsgIDS2C = 502
	MsgIDS2C_S2CReConnFtAck     MsgIDS2C = 503
	// -------------成长线--------------------
	MsgIDS2C_S2CHeroAck MsgIDS2C = 510
	MsgIDS2C_S2CItem    MsgIDS2C = 511
	MsgIDS2C_S2CMail    MsgIDS2C = 512
	// ----------------战斗--------------------
	MsgIDS2C_BattleYouLost         MsgIDS2C = 601
	MsgIDS2C_BattleCreateAck       MsgIDS2C = 602
	MsgIDS2C_BattleReconnectAck    MsgIDS2C = 603
	MsgIDS2C_BattleEnterAck        MsgIDS2C = 604
	MsgIDS2C_BattleLeaveAck        MsgIDS2C = 605
	MsgIDS2C_BattleStart           MsgIDS2C = 606
	MsgIDS2C_BattleFinish          MsgIDS2C = 607
	MsgIDS2C_BattleMonsterCreate   MsgIDS2C = 608
	MsgIDS2C_BattleTriggerEnter    MsgIDS2C = 609
	MsgIDS2C_BattleAreaScriptStart MsgIDS2C = 610
)

var MsgIDS2C_name = map[int32]string{
	0:   "S2CNone",
	1:   "Gm2GtKickRole",
	2:   "Ft2GtKickRole",
	100: "Gm2GtMsgIdMax",
	101: "S2CInit",
	102: "S2CHeartBeat",
	103: "S2CAccInfoAck",
	200: "S2CWaitMsgRetBegin",
	500: "S2CWaitMsgRetEnd",
	501: "S2CLoginAck",
	502: "S2CReConnAck",
	503: "S2CReConnFtAck",
	510: "S2CHeroAck",
	511: "S2CItem",
	512: "S2CMail",
	601: "BattleYouLost",
	602: "BattleCreateAck",
	603: "BattleReconnectAck",
	604: "BattleEnterAck",
	605: "BattleLeaveAck",
	606: "BattleStart",
	607: "BattleFinish",
	608: "BattleMonsterCreate",
	609: "BattleTriggerEnter",
	610: "BattleAreaScriptStart",
}
var MsgIDS2C_value = map[string]int32{
	"S2CNone":               0,
	"Gm2GtKickRole":         1,
	"Ft2GtKickRole":         2,
	"Gm2GtMsgIdMax":         100,
	"S2CInit":               101,
	"S2CHeartBeat":          102,
	"S2CAccInfoAck":         103,
	"S2CWaitMsgRetBegin":    200,
	"S2CWaitMsgRetEnd":      500,
	"S2CLoginAck":           501,
	"S2CReConnAck":          502,
	"S2CReConnFtAck":        503,
	"S2CHeroAck":            510,
	"S2CItem":               511,
	"S2CMail":               512,
	"BattleYouLost":         601,
	"BattleCreateAck":       602,
	"BattleReconnectAck":    603,
	"BattleEnterAck":        604,
	"BattleLeaveAck":        605,
	"BattleStart":           606,
	"BattleFinish":          607,
	"BattleMonsterCreate":   608,
	"BattleTriggerEnter":    609,
	"BattleAreaScriptStart": 610,
}

func (x MsgIDS2C) String() string {
	return proto.EnumName(MsgIDS2C_name, int32(x))
}
func (MsgIDS2C) EnumDescriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

func init() {
	proto.RegisterEnum("pb.MsgIDS2C", MsgIDS2C_name, MsgIDS2C_value)
}

func init() { proto.RegisterFile("msg_id_ser2cli.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	// 385 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xcb, 0x6e, 0x53, 0x31,
	0x10, 0x86, 0x69, 0x8e, 0x05, 0xc8, 0x6d, 0xe9, 0x74, 0xda, 0x0a, 0xc4, 0x23, 0xb0, 0x60, 0x61,
	0x9e, 0x20, 0x31, 0x4d, 0x1b, 0x91, 0xb0, 0x88, 0x91, 0x10, 0xab, 0xca, 0x75, 0xa6, 0x07, 0xab,
	0x27, 0x76, 0xe4, 0x33, 0x20, 0x96, 0x3c, 0x1e, 0xaf, 0xc0, 0xfd, 0xf2, 0x0c, 0x5c, 0x24, 0x16,
	0x80, 0x6c, 0x43, 0x14, 0x96, 0xfe, 0xc6, 0xfe, 0xfe, 0x7f, 0x24, 0xcb, 0xc3, 0x65, 0xdf, 0x9e,
	0xf9, 0xc5, 0x59, 0x4f, 0x49, 0xb9, 0xce, 0xdf, 0x5d, 0xa5, 0xc8, 0x11, 0x07, 0xab, 0xf3, 0x3b,
	0x3f, 0x1b, 0x79, 0x7d, 0xd6, 0xb7, 0x93, 0xfb, 0x46, 0x69, 0xdc, 0x96, 0xd7, 0x8c, 0xd2, 0x0f,
	0x63, 0x20, 0xb8, 0x82, 0xfb, 0x72, 0xf7, 0x64, 0xa9, 0x4e, 0xf8, 0x81, 0x77, 0x97, 0xf3, 0xd8,
	0x11, 0x6c, 0x65, 0x34, 0xe6, 0x4d, 0x34, 0x58, 0xdf, 0xca, 0x8e, 0xc5, 0xcc, 0xbe, 0x80, 0xc5,
	0x5f, 0xcb, 0x24, 0x78, 0x06, 0x42, 0x90, 0x3b, 0x46, 0xe9, 0x53, 0xb2, 0x89, 0x47, 0x64, 0x19,
	0x2e, 0xf2, 0x0b, 0xa3, 0xf4, 0xd0, 0xb9, 0x49, 0xb8, 0x88, 0x43, 0x77, 0x09, 0x2d, 0xde, 0x94,
	0x68, 0x94, 0x7e, 0x6c, 0x7d, 0xd6, 0xcc, 0x89, 0x47, 0xd4, 0xfa, 0x00, 0xaf, 0xb6, 0xf0, 0x48,
	0xc2, 0x7f, 0x83, 0xe3, 0xb0, 0x80, 0xaf, 0x0d, 0x82, 0xdc, 0x36, 0x4a, 0x4f, 0x63, 0xeb, 0x43,
	0x16, 0x7c, 0x6b, 0x70, 0xbf, 0xc4, 0xcc, 0x49, 0xc7, 0x50, 0xd0, 0xf7, 0x06, 0x0f, 0xe4, 0x8d,
	0x35, 0x1a, 0x73, 0x86, 0x3f, 0x1a, 0xdc, 0x93, 0xb2, 0xd4, 0x49, 0x25, 0xf9, 0x57, 0x83, 0x3b,
	0xb5, 0x2c, 0xd3, 0x12, 0x7e, 0xff, 0x3b, 0xcd, 0xac, 0xef, 0xe0, 0xa5, 0x40, 0x94, 0xbb, 0x23,
	0xcb, 0xdc, 0xd1, 0x93, 0xf8, 0x6c, 0x1a, 0x7b, 0x86, 0xd7, 0x02, 0x0f, 0xe5, 0x5e, 0x65, 0x3a,
	0x91, 0x65, 0xca, 0x96, 0x37, 0x22, 0x2f, 0x50, 0xe9, 0x9c, 0x5c, 0x0c, 0x81, 0x5c, 0xc9, 0x7b,
	0x2b, 0x72, 0x89, 0x3a, 0x38, 0x0e, 0x4c, 0x29, 0xc3, 0x77, 0x1b, 0x70, 0x4a, 0xf6, 0x79, 0x51,
	0xbc, 0x17, 0x79, 0xa7, 0x0a, 0x0d, 0xdb, 0xc4, 0xf0, 0x41, 0xe4, 0x9d, 0x2a, 0x19, 0xfb, 0xe0,
	0xfb, 0xa7, 0xf0, 0x51, 0xe0, 0x2d, 0x79, 0x50, 0xd1, 0x2c, 0x86, 0x9e, 0x29, 0xd5, 0x12, 0xf0,
	0x69, 0xa3, 0xc1, 0xa3, 0xe4, 0xdb, 0x96, 0x52, 0xc9, 0x83, 0xcf, 0x02, 0x6f, 0xcb, 0xa3, 0x3a,
	0x18, 0x26, 0xb2, 0xc6, 0x25, 0xbf, 0xe2, 0x9a, 0xf0, 0x45, 0x8c, 0x06, 0xa7, 0xcd, 0xf9, 0xd5,
	0xf2, 0x17, 0xee, 0xfd, 0x09, 0x00, 0x00, 0xff, 0xff, 0x66, 0x72, 0xbc, 0x15, 0x23, 0x02, 0x00,
	0x00,
}