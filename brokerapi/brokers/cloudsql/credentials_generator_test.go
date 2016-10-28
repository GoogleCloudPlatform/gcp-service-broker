package cloudsql_test

import (
	. "gcp-service-broker/brokerapi/brokers/cloudsql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("credentials generation", func() {
	It("does not generate a username for empty instanceID/bindingID", func() {
		_, err := GenerateUsername("", "")
		Expect(err).To(HaveOccurred())
	})
	It("generates a username", func() {
		u, err := GenerateUsername("foo", "bar")
		Expect(err).ToNot(HaveOccurred())
		Expect(len(u)).To(BeNumerically(">", 1))
	})
	It("truncates very long instanceID/bindingIDs", func() {
		longStr := "foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoo"
		u, err := GenerateUsername(longStr, longStr)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(u)).To(BeNumerically("<", len(longStr)))
	})
	It("generates a password", func() {
		p, err := GeneratePassword()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(p)).To(BeNumerically(">", 1))
	})
	It("generates unique passwords", func() {
		generated := map[string]bool{}
		for i := 0; i < 10; i++ {
			p, err := GeneratePassword()
			Expect(err).ToNot(HaveOccurred())
			Expect(generated[p]).To(BeFalse())
			generated[p] = true
		}
	})
})
