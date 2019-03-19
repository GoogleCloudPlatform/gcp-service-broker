package utils

import "math/rand"
import "time"

// Retry retries the call repeatedly with a given delay and jitter up to 200ms.
func Retry(times int, delay time.Duration, callback func() error) error {
	var err error

	for i := 0; i < times; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		jitter := time.Duration(rand.Int63n(200)) * time.Millisecond
		time.Sleep(delay + jitter)
	}

	return err
}
