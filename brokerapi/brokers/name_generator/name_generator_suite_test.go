package name_generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNameGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NameGenerator Suite")
}
