package lagerflags_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"code.cloudfoundry.org/lager/lagerflags"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeFormat", func() {
	const InvalidFormat = lagerflags.TimeFormat(123456)

	It("MarshalJSON", func() {
		b, err := json.Marshal(lagerflags.FormatUnixEpoch)
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(MatchJSON(`"unix-epoch"`))

		b, err = json.Marshal(lagerflags.FormatRFC3339)
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(MatchJSON(`"rfc3339"`))

		_, err = json.Marshal(InvalidFormat)
		Expect(err).To(HaveOccurred())
	})

	It("UnmarshalJSON", func() {
		var testCases = []struct {
			Format lagerflags.TimeFormat
			Data   string
			Valid  bool
		}{
			{
				Format: lagerflags.FormatUnixEpoch,
				Data:   `"unix-epoch"`,
				Valid:  true,
			},
			{
				Format: lagerflags.FormatRFC3339,
				Data:   `"rfc3339"`,
				Valid:  true,
			},
			// integer values
			{
				Format: lagerflags.FormatUnixEpoch,
				Data:   "0",
				Valid:  true,
			},
			{
				Format: lagerflags.FormatRFC3339,
				Data:   "1",
				Valid:  true,
			},
			// invalid
			{
				Format: InvalidFormat,
				Data:   "",
				Valid:  false,
			},
			{
				Format: lagerflags.FormatRFC3339,
				Data:   `"RFC3339"`,
				Valid:  false,
			},
		}
		for _, test := range testCases {
			var tf lagerflags.TimeFormat
			err := json.Unmarshal([]byte(test.Data), &tf)
			if !test.Valid {
				Expect(err).To(HaveOccurred())
				continue
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(tf).To(Equal(test.Format))
		}
	})

	Context("TimeFormat FlagSet", func() {
		var flagSet *flag.FlagSet
		var timeFormat lagerflags.TimeFormat

		BeforeEach(func() {
			timeFormat = InvalidFormat
			flagSet = flag.NewFlagSet("test", flag.ContinueOnError)
			flagSet.Usage = func() {}
			flagSet.SetOutput(nopWriter{})
			flagSet.Var(
				&timeFormat,
				"timeFormat",
				`Format for timestamp in component logs. Valid values are "unix-epoch" and "rfc3339".`,
			)
		})

		testValidTimeFormatFlag := func(expected lagerflags.TimeFormat, argument string) {
			Expect(flagSet.Parse([]string{"-timeFormat", argument})).To(Succeed())
			Expect(timeFormat).To(Equal(expected),
				fmt.Sprintf("Valid TimeFormat flag (expect: %q): %q", expected, argument))
		}

		testInvalidTimeFormatFlag := func(argument string) {
			Expect(flagSet.Parse([]string{"-timeFormat", argument})).ToNot(Succeed(),
				fmt.Sprintf("Invalid TimeFormat flag: %q", argument))
		}

		It("parses valid flags", func() {
			testValidTimeFormatFlag(lagerflags.FormatUnixEpoch, "unix-epoch")
			testValidTimeFormatFlag(lagerflags.FormatUnixEpoch, "0")
			testValidTimeFormatFlag(lagerflags.FormatRFC3339, "rfc3339")
			testValidTimeFormatFlag(lagerflags.FormatRFC3339, "1")
		})

		It("errors when the flag is invalid", func() {
			testInvalidTimeFormatFlag("UNIX-EPOCH")
			testInvalidTimeFormatFlag("RFC3339")
			testInvalidTimeFormatFlag("")
			testInvalidTimeFormatFlag(strconv.Itoa(int(InvalidFormat)))
		})
	})
})

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }
