package gorand

import "testing"

func TestUUID(t *testing.T) {
	uuid, err := UUIDv4()
	if err != nil {
		t.Error(err.Error())
	}

	m, err := MarshalUUID(uuid)
	if err != nil {
		t.Error(err.Error())
	}

	u, err := UnmarshalUUID(m)
	if err != nil {
		t.Error(err.Error())
	}

	if u != uuid {
		t.Errorf("%v != %v after Unmarshal and Marshal", u, uuid)
	}
}

func TestUnmarshalUUIDFail(t *testing.T) {
	_, err := UnmarshalUUID("1234567890")
	if err == nil {
		t.Fatal("UnmarshalUUID('1234567890') should fail.")
	}
	_, err = UnmarshalUUID("123456789012345678901234567890123456")
	if err == nil {
		t.Fatal("UnmarshalUUID('123456789012345678901234567890123456') should fail.")
	}
}

func BenchmarkUUID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		UUIDv4()
	}
}
