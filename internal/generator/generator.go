// Package generator implements our code generator.
package generator

import (
	"errors"
	"fmt"
	"io"
	"strings"

	. "github.com/dave/jennifer/jen" //nolint
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

// ErrTmplAbsolutePathNotAllowed is returned when a user specified a absolute path.
var ErrTmplAbsolutePathNotAllowed = errors.New("absolute path for template not supported")

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
func (tg *Target) GenerateRendering(wrt io.Writer) error {
	file := NewFile(string(tg.src.GoPackageName))
	file.HeaderComment("Code generated by protoc-gen-gossr. DO NOT EDIT.")

	for _, msg := range tg.src.Messages {
		tg.generateRenderingTmplNames(file, msg)
		tg.generateRenderingRegister(file, msg)
		tg.generateRenderingRender(file, msg)
	}

	if err := file.Render(wrt); err != nil {
		return fmt.Errorf("failed to render file: %w", err)
	}

	return nil
}

const (
	// protoreflect package path.
	protoreflect = "google.golang.org/protobuf/reflect/protoreflect"
)

// re-usable identifier for template names.
func tmplNamesIdent(msg *protogen.Message) string {
	return strcase.ToLowerCamel(msg.GoIdent.GoName) + "TmplNames"
}

// re-usable identifier for template files.
func tmplFilesIdent(msg *protogen.Message) string {
	return strcase.ToLowerCamel(msg.GoIdent.GoName) + "TmplFiles"
}

// generateRenderingTmplNames generates the template names and go:embed statement.
func (tg *Target) generateRenderingTmplNames(file *File, msg *protogen.Message) {
	names := messageTemplates(msg)
	lits := make([]Code, 0, len(names))

	for _, name := range names {
		lits = append(lits, Lit(name))
	}

	file.Func().Id(tmplNamesIdent(msg)).Params().Params(Index().String()).Block(
		Return(Index().String().Values(lits...)),
	)

	file.Commentf("//go:embed %s", strings.Join(names, " "))
	file.Var().Id(tmplFilesIdent(msg)).Qual("embed", "FS")
}

// generateRenderingRegister generates the registeration for embedded templates.
func (tg *Target) generateRenderingRegister(file *File, msg *protogen.Message) {
	// view interface for parameter
	view := Interface(
		Id("Parse").
			Params(Qual("io/fs", "FS"), String(), Op("...").String()).
			Params(Op("*").Qual("html/template", "Template"), Error()),
		Id("RegisterEmbedded").
			Params(Qual(protoreflect, "MessageDescriptor"), Op("*").Qual("html/template", "Template")).
			Params(Error()),
	)

	// generate the function sig and body
	file.Commentf("Register" + strcase.ToCamel(msg.GoIdent.GoName) + "Template registers the embedded template onto " +
		" the view.")
	file.Func().Id("Register"+strcase.ToCamel(msg.GoIdent.GoName)+"Template").
		Params(Id("view").Add(view)).
		Params(Error()).
		Block(
			List(Id("tmpl"), Err()).Op(":=").
				Id("view").Dot("Parse").
				Call(
					Id(tmplFilesIdent(msg)),
					Id(tmplNamesIdent(msg)).Call().Index(Lit(0)),
					Id(tmplNamesIdent(msg)).Call().Index(Lit(1).Op(":")).Op("..."),
				),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to parse: %w"), Id("err"))),
			),

			Id("err").Op("=").Id("view").Dot("RegisterEmbedded").Call(
				Params(Op("&").Id(msg.GoIdent.GoName).Values()).Dot("ProtoReflect").Call().Dot("Descriptor").Call(),
				Id("tmpl"),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to register: %w"), Id("err"))),
			),

			Return(Nil()),
		)
}

// generateRenderingRender generates the message's render method.
func (tg *Target) generateRenderingRender(file *File, msg *protogen.Message) {
	// view interface for parameter
	view := Interface(
		Id("LiveDir").Params().Params(String()),

		Id("Parse").
			Params(Qual("io/fs", "FS"), String(), Op("...").String()).
			Params(Op("*").Qual("html/template", "Template"), Error()),
		Id("Embedded").
			Params(Qual(protoreflect, "MessageDescriptor")).
			Params(Op("*").Qual("html/template", "Template"), Error()),
	)

	file.Commentf("Render renders the message using a template.")
	file.Func().Params(Id("x").Op("*").Id(msg.GoIdent.GoName)).
		Id("Render").
		Params(Id("wrt").Qual("io", "Writer"), Id("view").Add(view)).
		Params(Error()).
		Block(
			List(Id("tmpl"), Err()).Op(":=").Id("view").Dot("Embedded").Call(
				Id("x").Dot("ProtoReflect").Call().Dot("Descriptor").Call(),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to get embedded template: %w"), Id("err"))),
			),

			If(Id("liveDir").Op(":=").Id("view").Dot("LiveDir").Call(), Id("liveDir").Op("!=").Lit("")).Block(
				List(Id("tmpl"), Err()).Op("=").Id("view").Dot("Parse").Call(
					Qual("os", "DirFS").Call(Id("liveDir")),
					Id(tmplNamesIdent(msg)).Call().Index(Lit(0)),
					Id(tmplNamesIdent(msg)).Call().Index(Lit(1).Op(":")).Op("..."),
				),
				If(Err().Op("!=").Nil()).Block(
					Return(Qual("fmt", "Errorf").Call(Lit("failed to parse template: %w"), Id("err"))),
				),
			),

			Err().Op("=").Id("tmpl").Dot("Execute").Call(Id("wrt"), Id("x")),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to get execute template: %w"), Id("err"))),
			),

			Return(Nil()),
		)
}
