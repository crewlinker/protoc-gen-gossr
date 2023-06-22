package gossr_test

import (
	"bytes"
	"embed"
	"html/template"
	"testing"

	"github.com/crewlinker/protoc-gen-gossr/gossr"
	blogv1 "github.com/crewlinker/protoc-gen-gossr/proto/examples/blog/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGossr(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "gossr")
}

//go:embed testdata
var testdata embed.FS

var _ = Describe("view", func() {
	var view1 *gossr.View
	BeforeEach(func() {
		view1 = gossr.New("", template.FuncMap{})
	})

	It("should allow parsing simple partial", func() {
		tmpl1, err := view1.Parse(testdata, "testdata/partial1.html")
		Expect(err).ToNot(HaveOccurred())
		buf1 := bytes.NewBuffer(nil)

		Expect(tmpl1.Execute(buf1, struct{}{})).To(Succeed())
		Expect(buf1.String()).To(Equal("Hi from partial"))
	})

	It("should allow parsing a extend chain", func() {
		tmpl1, err := view1.Parse(testdata, "testdata/layout1.html", "testdata/layout2.html", "testdata/page1.html")
		Expect(err).ToNot(HaveOccurred())
		buf1 := bytes.NewBuffer(nil)

		Expect(tmpl1.Execute(buf1, struct{}{})).To(Succeed())
		Expect(buf1.String()).To(Equal("<body><main>page1</main></body>"))
	})

	It("should allow registering", func() {
		tmpl1 := template.Must(template.New("foo").Parse(`foo`))
		desc1 := (&blogv1.BlogIndex{}).ProtoReflect().Descriptor()
		desc2 := (&blogv1.BlogIndex{}).ProtoReflect().Descriptor()

		Expect(view1.RegisterEmbedded(desc1, tmpl1)).To(Succeed())
		Expect(view1.Embedded(desc1)).To(Equal(tmpl1))
		Expect(view1.Embedded(desc2)).To(Equal(tmpl1))
	})

	It("should error if already registered", func() {
		tmpl1 := template.Must(template.New("foo").Parse(`foo`))
		desc1 := (&blogv1.BlogIndex{}).ProtoReflect().Descriptor()
		desc2 := (&blogv1.BlogIndex{}).ProtoReflect().Descriptor()
		Expect(view1.RegisterEmbedded(desc1, tmpl1)).To(Succeed())
		Expect(view1.RegisterEmbedded(desc2, tmpl1)).To(MatchError(gossr.ErrTmplAlreadyRegistered))
	})

	It("should error if not registered", func() {
		desc1 := (&blogv1.BlogIndex{}).ProtoReflect().Descriptor()
		_, err := view1.Embedded(desc1)
		Expect(err).To(MatchError(gossr.ErrTmplNotRegistered))
	})

	It("should return empty livedir by default", func() {
		Expect(view1.LiveDir()).To(BeEmpty())
	})
})
