package util

import (
	"math/rand"
)

const (
	MAX_RATE = 10000
)

type randerData struct {
	rate int
	data interface{}
}

type Rander struct {
	data []*randerData
}

func NewRander() *Rander {
	r := &Rander{
		data: make([]*randerData, 0),
	}
	return r
}

func (r *Rander) Add(rate int, v interface{}) {
	d := &randerData{
		rate: rate,
		data: v,
	}
	r.data = append(r.data, d)
}

func (r *Rander) Get() interface{} {
	rate := rand.Intn(MAX_RATE)
	cr := 0
	var ret interface{}
	for _, v := range r.data {
		cr += v.rate
		ret = v.data
		if cr >= rate {
			return ret
		}
	}
	return ret
}
