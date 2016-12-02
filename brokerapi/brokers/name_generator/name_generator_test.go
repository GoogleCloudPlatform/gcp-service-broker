package name_generator_test

import (
	. "gcp-service-broker/brokerapi/brokers/name_generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testUniqueness(generator func() string) {
	generated := map[string]bool{}
	for i := 0; i < 10; i++ {
		name := generator()
		Expect(generated[name]).To(BeFalse())
		generated[name] = true
	}
}

var _ = Describe("NameGenerator", func() {
	Describe("BasicNameGenerator", func() {
		var (
			generator BasicNameGenerator
		)
		It("Generates a name", func() {
			Expect(generator.InstanceName()).To(Not(BeEmpty()))
		})
		It("Generates unique names", func() {
			testUniqueness(func() string { return generator.InstanceName() })
		})
	})
	Describe("SqlNameGenerator", func() {
		var (
			generator SqlNameGenerator
		)
		It("Generates a name", func() {
			Expect(generator.InstanceName()).To(Not(BeEmpty()))
			Expect(generator.DatabaseName()).To(Not(BeEmpty()))
		})
		It("Generates unique names", func() {
			testUniqueness(func() string { return generator.InstanceName() })
			testUniqueness(func() string { return generator.DatabaseName() })
		})
		Describe("credentials generation", func() {
			It("does not generate a username for empty instanceID/bindingID", func() {
				_, err := generator.GenerateUsername("", "")
				Expect(err).To(HaveOccurred())
			})
			It("generates a username", func() {
				u, err := generator.GenerateUsername("foo", "bar")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(u)).To(BeNumerically(">", 1))
			})
			It("truncates very long instanceID/bindingIDs", func() {
				longStr := "foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoo"
				u, err := generator.GenerateUsername(longStr, longStr)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(u)).To(BeNumerically("<", len(longStr)))
			})
			It("generates a password", func() {
				p, err := generator.GeneratePassword()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(p)).To(BeNumerically(">", 1))
			})
			It("generates unique passwords", func() {
				testUniqueness(func() string {
					val, err := generator.GeneratePassword()
					Expect(err).ToNot(HaveOccurred())
					return val
				})
			})
		})
	})
})
