package enet

import (
	"encoding/binary"
	"errors"
)

//---------------------------------------------------
//msgid:2|data
type PkgReader struct {
	Data []byte
}

func (p *PkgReader) GetMsgID() uint16 {
	return binary.BigEndian.Uint16(p.Data[0:2])
}

func (p *PkgReader) GetSeqNum() uint32 {
	return 0
}

func (p *PkgReader) GetData() []byte {
	return p.Data[2:]
}

//-----------------------------------------------
//len:4|msgID:2|data
type PkgWriter struct {
	MsgId uint16
	Data  []byte
}

func (p *PkgWriter) Write(retCache []byte) (totalLen int) {
	msgLen := len(p.Data) + 2
	totalLen = msgLen + 4 //total len
	binary.BigEndian.PutUint32(retCache[0:4], uint32(msgLen))
	binary.BigEndian.PutUint16(retCache[4:6], p.MsgId)
	copy(retCache[6:], p.Data)
	return
}

//-----------------------------------------------
type Pkg struct {
}

func (p *Pkg) NewReader(data []byte, s ISession) (IPkgReader, error) {
	if len(data) < 2 {
		return nil, errors.New("read packet head error")
	}
	return &PkgReader{Data: data}, nil
}

func (p *Pkg) NewWriter(msgID uint16, data []byte, s ISession) IPkgWriter {
	r := &PkgWriter{
		MsgId: msgID,
		Data:  data,
	}
	return r
}

//----------------------------------------------------
