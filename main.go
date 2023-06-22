// Package main holds the main program's entrypoint
package main

import (
	"flag"
	"fmt"

	"github.com/crewlinker/protoc-gen-gossr/internal/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// run the plugin logic.
func run(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	// initialize our generator
	gen := generator.NewGenerator()

	for _, name := range plugin.Request.FileToGenerate {
		// get the file the request wants to generate for
		plugf := plugin.FilesByPath[name]
		if len(plugf.Messages) < 1 {
			continue // no messages in the target at all
		}

		// init a target on for the protobuf file
		tgt := gen.NewTarget(plugf)
		if !tgt.HasRenderableMessages() {
			continue // no messages for rendering in the taget.
		}

		// initialize a generated file to write to
		genf := plugin.NewGeneratedFile(fmt.Sprintf("%s.render.go", plugf.GeneratedFilenamePrefix), plugf.GoImportPath)

		// generate the actual rendering code
		if err := tgt.GenerateRendering(genf); err != nil {
			return fmt.Errorf("failed to generate rendering: %w", err)
		}
	}

	return nil
}

// programs entrypoint.
func main() {
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(run)
}
