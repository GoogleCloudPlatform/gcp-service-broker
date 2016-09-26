package brokers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBrokers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brokers Suite")
}
