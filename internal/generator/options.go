package generator

import (
	gossrv1 "github.com/crewlinker/protoc-gen-gossr/proto/gossr/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// messageOptions returns our plugin specific options for a field. If the field has no options
// it returns nil.
func messageOptions(m *protogen.Message) *gossrv1.MessageOptions {
	opts, hasOpts := m.Desc.Options().(*descriptorpb.MessageOptions)
	if !hasOpts {
		return nil
	}

	ext, hasOpts := proto.GetExtension(opts, gossrv1.E_Msg).(*gossrv1.MessageOptions)
	if !hasOpts {
		return nil
	}

	if ext == nil {
		return nil
	}

	return ext
}

// returns message templates for a message.
func messageTemplates(m *protogen.Message) []string {
	if fopts := messageOptions(m); fopts != nil && fopts.Template != nil {
		return fopts.Template
	}

	return nil
}
