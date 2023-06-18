// Package main describes mage targets
package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/magefile/mage/sh"
)

// init performs some sanity checks before running anything.
func init() {
	mustBeInRoot()
}

// Generate generates protobuf Go files.
func Generate() error {
	if err := sh.Run("buf", "generate"); err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	return nil
}

// Check checks the codebase using static analysis.
func Check() error {
	if err := sh.Run("golangci-lint", "run"); err != nil {
		return fmt.Errorf("failed lint: %w", err)
	}

	return nil
}

// Test tests the whole project's unit tests.
func Test() error {
	if err := sh.Run(
		"go", "run", "-mod=readonly", "github.com/onsi/ginkgo/v2/ginkgo",
		"-p", "-randomize-all", "-repeat=5", "--fail-on-pending", "--race", "--trace",
		"--junit-report=test-report.xml",
		"./...",
	); err != nil {
		return fmt.Errorf("failed to run ginkgo: %w", err)
	}

	return nil
}

// error when wrong version format is used.
var errVersionFormat = fmt.Errorf("version must be in format vX,Y,Z")

// Release releases the current code under a new version.
func Release(version string) error {
	if !regexp.MustCompile(`^v([0-9]+).([0-9]+).([0-9]+)$`).Match([]byte(version)) {
		return errVersionFormat
	}

	if err := sh.Run("git", "tag", version); err != nil {
		return fmt.Errorf("failed to tag version: %w", err)
	}

	if err := sh.Run("git", "push", "origin", version); err != nil {
		return fmt.Errorf("failed to push version tag: %w", err)
	}

	return nil
}

// mustBeInRoot checks that the command is run in the project root.
func mustBeInRoot() {
	var err error
	if _, err = os.ReadFile("go.mod"); err != nil {
		panic("must be in project root, couldn't stat go.mod file: " + err.Error())
	}
}
