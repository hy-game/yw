package util

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type IDPool struct {
	ids []uint16
	len uint16
	mtx sync.Mutex
}

//NewIDPool id pool
func NewIDPool(size uint16) *IDPool {
	m := &IDPool{
		ids: make([]uint16, size),
		len: size,
	}

	for i := uint16(0); i < size; i++ {
		m.ids[i] = size - i
	}

	return m
}

//GetID return 0 is failed
func (m *IDPool) GetID() uint16 {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.len == 0 {
		log.Error("GetID failed")
		return 0
	}
	m.len--
	return m.ids[m.len]
}

func (m *IDPool) RelaseID(id uint16) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if int(m.len) >= len(m.ids) {
		log.Error("Release ID failed")
		return
	}
	m.ids[m.len] = id
	m.len++
}
