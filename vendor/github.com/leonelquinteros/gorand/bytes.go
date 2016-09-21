/*
Package gorand defines a set of utility methods to work with random data.

The main goal is to curate a collection of functions for random data generation
in different formats to be used for different purposes.

It uses the "crypto/rand" package to provide the most secure random data,
unless indicated otherwise where performance can be the main focus of the method.

Most implementations are really trivial and anybody could write the same on their own packages,
but here we can centralize all of them and keep a unified way of retrieving random data.

Unified QA is another motivator to have and use this package.
*/
package gorand

import (
	"crypto/rand"
)

// GetBytes returns a fixed amount of random bytes.
// Specify the amount of bytes wanted by passing it on the n parameter.
// This function is the base for most of the methods on this package.
func GetBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
