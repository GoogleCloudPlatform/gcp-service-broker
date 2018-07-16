package gorand

import (
	"encoding/hex"
	"testing"
)

func TestID(t *testing.T) {
	id, err := ID()
	if err != nil {
		t.Error(err.Error())
	}

	if len(id) != 64 {
		t.Error("Length of ID isn't 64 bytes")
	}
}

func TestIDOrder(t *testing.T) {
	ids := make([][64]byte, 100000)
	for i := 0; i < 100000; i++ {
		id, err := ID()
		if err != nil {
			t.Fatal(err)
		}
		ids[i] = id
	}

	for i, v := range ids {
		if i == 0 {
			continue
		}
		if hex.EncodeToString(v[:]) <= hex.EncodeToString(ids[i-1][:]) {
			t.Fatalf("%v is lesser or equal than %v", hex.EncodeToString(v[:]), hex.EncodeToString(ids[i-1][:]))
		}
	}
}

func BenchmarkID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ID()
	}
}
