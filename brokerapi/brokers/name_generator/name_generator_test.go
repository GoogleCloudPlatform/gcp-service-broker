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
	})
})
