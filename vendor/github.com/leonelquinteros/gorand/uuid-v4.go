package gorand

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
)

// UUID as defined in RFC 4122
type UUID [16]byte

// UUIDv4 generates a version 4 (randomly generated) UUID
func UUIDv4() (uuid UUID, err error) {
	// Get 16 random bytes
	b, err := GetBytes(16)
	if err != nil {
		return
	}

	_, err = io.ReadFull(bytes.NewBuffer(b), uuid[:])
	if err != nil {
		return
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return
}

// UnmarshalUUID parses a string representation of a UUID and returns its []bytes value.
// It doesn't check for version or varian bits, so it can be used with invalid (non RFC 4122 compilant) values.
func UnmarshalUUID(s string) (uuid UUID, err error) {
	if len(s) != 36 {
		err = errors.New("Invalid UUID length")
		return
	}
	if s[8:9] != "-" || s[13:14] != "-" || s[18:19] != "-" || s[23:24] != "-" {
		err = errors.New("Invalid UUID format")
		return
	}

	b, err := hex.DecodeString(s[0:8] + s[9:13] + s[14:18] + s[19:23] + s[24:])
	if err != nil {
		return
	}

	_, err = io.ReadFull(bytes.NewBuffer(b), uuid[:])
	return
}

// MarshalUUID converts UUID into its canonical string representation.
// It doesn't check for version or varian bits, so it can be used with invalid (non RFC 4122 compilant) values.
func MarshalUUID(uuid UUID) (s string, err error) {
	var b [16]byte

	_, err = io.ReadFull(bytes.NewBuffer(uuid[:]), b[:])
	if err != nil {
		return
	}

	bin := hex.EncodeToString(b[:])
	s = bin[0:8] + "-" + bin[8:12] + "-" + bin[12:16] + "-" + bin[16:20] + "-" + bin[20:]

	return
}
