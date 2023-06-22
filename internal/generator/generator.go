// Package generator implements our code generator.
package generator

import (
	"fmt"
	"io"

	. "github.com/dave/jennifer/jen" //nolint
	"google.golang.org/protobuf/compiler/protogen"
)

// Generator implements the generator logic.
type Generator struct{}

// NewGenerator initializes a new generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// NewTarget inits a target.
func (gen *Generator) NewTarget(src *protogen.File) *Target {
	return &Target{gen: gen, src: src}
}

// Target represents a single proto file to generate for.
type Target struct {
	gen *Generator
	src *protogen.File
}

// HasRenderableMessages returns true if the target has at least 1 message that
// requires generating rendering logic.
func (tg *Target) HasRenderableMessages() bool {
	for _, msg := range tg.src.Messages {
		if len(messageTemplates(msg)) > 0 {
			return true
		}
	}

	return false
}

// GenerateRendering generates the rendering code for the target.
func (tg *Target) GenerateRendering(w io.Writer) error {
	f := NewFile(string(tg.src.GoPackageName))
	f.HeaderComment("Code generated by protoc-gen-gossr. DO NOT EDIT.")

	if err := f.Render(w); err != nil {
		return fmt.Errorf("failed to render file: %w", err)
	}

	return nil
}