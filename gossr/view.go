// Package gossr provides shared coded for generated rendering code.
package gossr

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// ErrTmplAlreadyRegistered is returned when a template for a proto message is already registered.
	ErrTmplAlreadyRegistered = errors.New("template is already registered, check if .InitRender is not called twice")
	// ErrTmplNotRegistered is returned when a template for a proto message has not been registered.
	ErrTmplNotRegistered = errors.New("message template is not registered, check if .InitRender is called")
)

// View provides template rendering functionality.
type View struct {
	liveDir string
	funcs   template.FuncMap
	embeds  map[protoreflect.MessageDescriptor]*template.Template
	mu      sync.RWMutex
}

// New inits the view.
func New(liveDir string, funcs template.FuncMap) *View {
	return &View{
		liveDir: liveDir,
		embeds:  make(map[protoreflect.MessageDescriptor]*template.Template),
		funcs:   funcs,
	}
}

// Parse should parse templates from the provided filesystem. Funcs should be to the template before parsing.
func (v *View) Parse(tmplfs fs.FS, name string, names ...string) (*template.Template, error) {
	var root *template.Template

	for idx, name := range append([]string{name}, names...) {
		data, err := fs.ReadFile(tmplfs, name)
		if err != nil {
			return nil, fmt.Errorf("failed to read template '%s': %w", name, err)
		}

		if idx == 0 {
			root, err = template.New(name).Funcs(v.funcs).Parse(string(data))
		} else {
			_, err = root.New(name).Parse(string(data))
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse template '%s': %w", name, err)
		}
	}

	return root, nil
}

// LiveDir should return a non-empty string when rendering should re-parse the ondisk template. Mostly for
// development purposes.
func (v *View) LiveDir() string {
	return v.liveDir
}

// RegisterEmbedded will store the embedded template so it can be retrieved during rendering if live-mode
// is not enabled.
func (v *View) RegisterEmbedded(rfm protoreflect.MessageDescriptor, tmpl *template.Template) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, ok := v.embeds[rfm]; ok {
		return fmt.Errorf("%w", ErrTmplAlreadyRegistered)
	}

	v.embeds[rfm] = tmpl

	return nil
}

// Embedded homepage should return the template as parsed during initialization. Parsed from
// embedded resources.
func (v *View) Embedded(rfm protoreflect.MessageDescriptor) (*template.Template, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	tmpl, ok := v.embeds[rfm]
	if !ok {
		return nil, fmt.Errorf("%w", ErrTmplNotRegistered)
	}

	return tmpl, nil
}
