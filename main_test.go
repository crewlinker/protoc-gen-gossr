package main

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProtocGenGxlang(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "main")
}

var _ = Describe("e2e", func() {
	It("should run the generator", func(ctx context.Context) {
		errb := bytes.NewBuffer(nil)
		cmd := exec.CommandContext(ctx, "buf", "generate")
		cmd.Stderr = errb
		Expect(cmd.Run()).To(Succeed(), "failed to run: "+errb.String())
	})
})
