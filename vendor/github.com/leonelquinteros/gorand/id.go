package gorand

import (
	"encoding/hex"
	"errors"
)

// ID is a wrapper for GetHex(64).
// 64 bytes should be enough to use as unique IDs to avoid collisions between generated values.
func ID() (string, error) {
	return GetHex(64)
}

// UUID generates a version 4 (randomly generated) UUID as defined in RFC 4122.
// It returns the canonical string representation "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx".
func UUID() (string, error) {
	// Get initial random bytes
	uuid, err := GetHex(16)
	if err != nil {
		return "", err
	}

	// Format UUID with version and variant bits
	return uuid[0:8] + "-" + uuid[8:12] + "-4" + uuid[13:16] + "-a" + uuid[17:20] + "-" + uuid[20:], nil
}

// UnmarshalUUID parses a string representation of a UUID and returns its []bytes value.
// It doesn't check for version or varian bits, so it can be used with invalid (non RFC 4122 compilant) values.
func UnmarshalUUID(uuid string) ([]byte, error) {
	if len(uuid) != 36 {
		return nil, errors.New("Invalid uuid length")
	}
	if uuid[8:9] != "-" || uuid[13:14] != "-" || uuid[18:19] != "-" || uuid[23:24] != "-" {
		return nil, errors.New("Invalid uuid format")
	}

	return hex.DecodeString(uuid[0:8] + uuid[9:13] + uuid[14:18] + uuid[19:23] + uuid[24:])
}

// MarshalUUID converts a []byte UUID into its canonical string representation.
// It doesn't check for version or varian bits, so it can be used with invalid (non RFC 4122 compilant) values.
func MarshalUUID(b []byte) (string, error) {
	if len(b) != 16 {
		return "", errors.New("Invalid bytes length")
	}

	uuid := hex.EncodeToString(b)

	return uuid[0:8] + "-" + uuid[8:12] + "-" + uuid[12:16] + "-" + uuid[16:20] + "-" + uuid[20:], nil
}
