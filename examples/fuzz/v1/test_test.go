package fuzzv1_test

import (
	"regexp"
	"testing"
	"time"

	fuzzv1 "github.com/crewlinker/protoc-gen-gossr/examples/fuzz/v1"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestV1(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "examples/fuzz/v1")
}

var _ = Describe("fuzz complex message", func() {
	DescribeTable("generating protobuf messages", func(seed int64) {
		fzr := fuzz.NewWithSeed(seed).
			NilChance(0.3).
			MaxDepth(5).
			SkipFieldsWithPattern(regexp.MustCompile(`OneofField`))

		for i := int64(0); i < 100; i++ {
			var target fuzzv1.TestAllTypes
			fzr.Fuzz(&target)
		}
	},
		Entry("current time", int64(time.Now().Nanosecond())))
})

var _ = Describe("test valid xml", func() {
	It("should be valid xml", func() {
	})
})
