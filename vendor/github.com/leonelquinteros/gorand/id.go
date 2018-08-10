package gorand

import (
	"bytes"
	"io"
	"time"
)

var localID [9]byte

// Initializes the value for the local process run identifier
func init() {
	buf, err := GetBytes(9)
	if err != nil {
		localID = [9]byte{'D', 'e', 'f', 'a', 'u', 'l', 't', 'I', 'D'}
	} else {
		_, err = io.ReadFull(bytes.NewBuffer(buf), localID[:])
		if err != nil {
			localID = [9]byte{'D', 'e', 'f', 'a', 'u', 'l', 't', 'I', 'D'}
		}
	}
}

// ID generates a [64]byte random value, using time and local identifier into it.
//
// First (most-significative) 15 bytes: time value
// Next 9 bytes: Local process randomly-generated identifier
// Next 40 bytes: Random value
func ID() (id [64]byte, err error) {
	var buf []byte

	// Time part
	now, err := time.Now().MarshalBinary()
	if err != nil {
		return
	}
	buf = append(buf, now[:15]...)

	// Local ID part
	buf = append(buf, localID[:]...)

	// Random value
	r, err := GetBytes(40)
	if err != nil {
		return
	}
	buf = append(buf, r...)

	_, err = io.ReadFull(bytes.NewBuffer(buf), id[:])
	return
}
