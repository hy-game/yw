package mq

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
)

func WritePkg(msgId uint16, msg proto.Message, serId uint16) ([]byte, error) {
	var b []byte
	var err error

	if msg != nil {
		b, err = proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
	}

	msgLen := len(b) + 4
	data := make([]byte, msgLen)
	binary.BigEndian.PutUint16(data[0:2], msgId)
	binary.BigEndian.PutUint16(data[2:4], serId)
	copy(data[4:], b)
	return data, nil
}

func ParserPkg(data []byte) (msgId uint16, msg []byte, serId uint16) {
	msgId = binary.BigEndian.Uint16(data[0:2])
	serId = binary.BigEndian.Uint16(data[2:4])
	msg = data[4:]
	return
}
