package util

import (
	"fmt"
	"testing"
)

func TestIDPool(t *testing.T) {
	p := NewIDPool(30)

	ids := make([]uint16, 20)
	for i := 0; i < 20; i++ {
		ids[i] = p.GetID()
	}

	fmt.Print("1\n")
	fmt.Print(p.ids)
	fmt.Print(ids)

	id := p.GetID()

	for i := 0; i < 20; i++ {
		p.RelaseID(ids[i])
	}

	p.RelaseID(id)

	fmt.Print("2\n")
	fmt.Print(p.ids)
	fmt.Print(ids)

	ids2 := make([]uint16, 10)
	for i := 0; i < 10; i++ {
		ids2[i] = p.GetID()
	}
	for i := 0; i < 10; i++ {
		p.RelaseID(ids2[i])
	}

	fmt.Print("3\n")
	fmt.Print(p.ids)
	fmt.Print(ids)
	fmt.Print("ok")
}
