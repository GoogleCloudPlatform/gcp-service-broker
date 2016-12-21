package db_service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDbService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DbService Suite")
}

