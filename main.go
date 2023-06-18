// Package main holds the main program's entrypoint
package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// run the plugin logic.
func run(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	return nil
}

// programs entrypoint.
func main() {
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(run)
}
