package gorand

import (
	"testing"
)

func TestGetHex(t *testing.T) {
	str, err := GetHex(8)
	if err != nil {
		t.Error(err.Error())
	}

	if len(str) != 16 {
		t.Error("Length of string isn't 16")
	}
}

func BenchmarkGetHex1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetHex(1)
	}
}

func BenchmarkGetHex10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetHex(10)
	}
}

func BenchmarkGetHex100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetHex(100)
	}
}

func BenchmarkGetHex1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetHex(1000)
	}
}
