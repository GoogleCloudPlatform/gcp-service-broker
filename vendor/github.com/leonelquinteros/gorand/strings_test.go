package gorand

import (
	"testing"
)

func TestGetRandomChars(t *testing.T) {
	str, err := GetRandomChars("abc", 10)
	if err != nil {
		t.Error(err.Error())
	}

	if len(str) != 10 {
		t.Error("Length of string isn't 10")
	}
}

func TestGetAlphaNumString(t *testing.T) {
	str, err := GetAlphaNumString(10)
	if err != nil {
		t.Error(err.Error())
	}

	if len(str) != 10 {
		t.Error("Length of string isn't 10")
	}
}

func TestGetAlphaString(t *testing.T) {
	str, err := GetAlphaString(10)
	if err != nil {
		t.Error(err.Error())
	}

	if len(str) != 10 {
		t.Error("Length of string isn't 10")
	}
}

func TestGetNumString(t *testing.T) {
	str, err := GetNumString(10)
	if err != nil {
		t.Error(err.Error())
	}

	if len(str) != 10 {
		t.Error("Length of string isn't 10")
	}
}

func BenchmarkAlphaNumString1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetAlphaNumString(1)
	}
}

func BenchmarkAlphaNumString10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetAlphaNumString(10)
	}
}

func BenchmarkAlphaNumString100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetAlphaNumString(100)
	}
}

func BenchmarkAlphaNumString1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetAlphaNumString(1000)
	}
}
