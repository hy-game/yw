package gnet

import (
	"crypto/rc4"
	"encoding/binary"
	"errors"
)

//---------------------------------------------------
func NewReader(data []byte) (*PkgReader, error) {
	if len(data) < 6 {
		return nil, errors.New("packet head < 6")
	}
	return &PkgReader{data: data}, nil
}

//msgid:2/seqNum:4/data
type PkgReader struct {
	data []byte
}

func (p *PkgReader) GetMsgID() uint16 {
	return binary.BigEndian.Uint16(p.data[0:2])
}

func (p *PkgReader) GetSeqNum() uint32 {
	return binary.BigEndian.Uint32(p.data[2:6])
}

func (p *PkgReader) GetData() []byte {
	return p.data[6:]
}

//---------------------------------------------------
type PkgWriter struct {
	msgId   uint16
	data    []byte
	crypto  bool
	EnCoder *rc4.Cipher
}

func (p *PkgWriter) Write(retCache []byte) int {
	msgLen := len(p.data) + 2

	binary.BigEndian.PutUint32(retCache[0:4], uint32(msgLen))
	binary.BigEndian.PutUint16(retCache[4:6], p.msgId)
	copy(retCache[6:], p.data)

	endPos := msgLen + 4
	//if p.crypto {
	//	p.EnCoder.XORKeyStream(retCache[4:endPos], retCache[4:endPos])
	//}

	return endPos
}
