package mail

import (
	"center/network"
	"com/database/db"
	"com/log"
	"database/sql"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"pb"
)

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
	t := &DBMail{
		UID:uid.String(),
		RoleGuid:roleGuid,
		Data:b,
		msg: data,
	}
	db.Write(t)
	return true
}

type DBMail struct {
	UID string	`gorm:"primary_key"`
	RoleGuid uint32	`gorm:"index:idx_role"`
	Data []byte	`gorm:"type:blob(104857600)"`
	msg *pb.MsgMailData
}

func (t *DBMail)Run(conn *sql.DB){
	_, err := conn.Exec("INSERT INTO db_mail VALUES (?, ?, ?)", t.UID, t.RoleGuid, t.Data)
	if err != nil {
		log.Warnf("send mail err:%v", err)
		return
	}

	network.SendToGm(network.GetGameIDByRole(t.RoleGuid), pb.MsgIDS2S_Ct2GmSendMail, &pb.MsgMail{
		UUID:     t.UID,
		RoleGuid: t.RoleGuid,
		Data:     t.msg,
	})
}