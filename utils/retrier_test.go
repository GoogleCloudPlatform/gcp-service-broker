package utils

import (
	"errors"
	"fmt"
	"time"
)

func ExampleRetry_full() {
	i := 0

	err := Retry(3, 1*time.Second, func() error {
		i++
		if i != 3 {
			fmt.Printf("failing: %d\n", i)
			return errors.New("error")
		}

		fmt.Printf("succeeding: %d\n", i)
		return nil
	})

	fmt.Println(err)

	// Output: failing: 1
	// failing: 2
	// succeeding: 3
	// <nil>
}

func ExampleRetry_fail() {
	err := Retry(15, 1*time.Millisecond, func() error {
		return errors.New("error")
	})

	fmt.Println(err)

	// Output: error
}
