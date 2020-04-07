package mail

import (
	"com/log"
	"database/sql"
	"game/gmdb"
	"game/types"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"pb"
)

type Data struct {
	mail map[string]*pb.MsgMail
}

//NewData	一些初始化操作
func NewData(r *types.Role) *Data {
	return &Data{mail:make(map[string]*pb.MsgMail)}
}

func getData(r *types.Role)*Data{
	c := r.GetComp(types.TCMail)
	if c == nil {
		return nil
	}
	return c.(*Data)
}

func Load(r *types.Role){

}

func OnProto(msgRecv *pb.MsgMailProto, r *types.Role){
	switch msgRecv.Op {
	//收附件
	//删除
	}
}

//Add	添加一封邮件
func Add(m *pb.MsgMail, r *types.Role){
	c := getData(r)
	if c == nil {
		return
	}
	c.mail[m.UUID] = m

	msg := &pb.MsgMailProto{Op:pb.MsgMailProto_Add,
		OneMail:m,
		}
	r.Send(pb.MsgIDS2C_S2CMail, msg)
}

//SendMail 给指定玩家发送邮件
func SendMail(id uint32, roleGuid uint32, params []string, prize []*pb.CPriceItem) bool{
	uid, err := uuid.NewV4()
	if err != nil {
		log.Warnf("new uuid err:%v when SendMail", err)
		return false
	}
	data:= &pb.MsgMailData{
		CfgID:     id,
		Params: params,
		Prize:  prize,
		State:  pb.EMailState_EMSInit,
	}

	b, err := proto.Marshal(data)
	if err != nil {
		log.Warnf("marshal err:%v when SendMail", err)
		return  false
	}
	t := &dbSaveMail{
		UID:uid.String(),
		RoleGuid:roleGuid,
		Data:b,
		msg: data,
	}
	gmdb.DBMisc.Write(t)
	return true
}


type dbSaveMail struct {
	UID string	`gorm:"primary_key"`
	RoleGuid uint32	`gorm:"index:idx_role"`
	Data []byte	`gorm:"type:blob(104857600)"`
	msg *pb.MsgMailData
}

func (t *dbSaveMail)Run(conn *sql.DB){
	_, err := conn.Exec("INSERT INTO db_mail VALUES (?, ?, ?)", t.UID, t.RoleGuid, t.Data)
	if err != nil {
		log.Warnf("send mail err:%v", err)
		return
	}
	types.PostEvt(t.RoleGuid, types.Evt{
		Type: types.RecvMail,
		Data: &pb.MsgMail{
			UUID:     t.UID,
			RoleGuid: t.RoleGuid,
			Data:     t.msg,
		},
	})
}

type loadMail struct {
	RoleGuid uint32
}

func (t *loadMail)Run(conn *sql.DB){

}