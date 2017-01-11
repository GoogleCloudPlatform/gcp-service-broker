package gorand

import (
	"testing"
)

func TestGetBytes(t *testing.T) {
	b, err := GetBytes(10)
	if err != nil {
		t.Error(err.Error())
	}

	if len(b) != 10 {
		t.Error("Length of bytes slice isn't 10")
	}
}

func BenchmarkGetBytes1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetBytes(1)
	}
}

func BenchmarkGetBytes10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetBytes(10)
	}
}

func BenchmarkGetBytes100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetBytes(100)
	}
}

func BenchmarkGetBytes1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetBytes(1000)
	}
}
