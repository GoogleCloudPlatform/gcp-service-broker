package cloudsql_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCloudsql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CloudSQL Suite")
}
