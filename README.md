# protoc-gen-gossr

Fearless server-side rendering in Go using Protobuf.

## Introduction

This project aims to enhance the safety and simplicity of developing Go applications that utilize the html/template package by leveraging Protocol Buffers (Protobuf).

The html/template (or text/template) package in Go provides a powerful and flexible way to generate HTML, XML, and other text-based formats. However, working directly with strings and template syntax can be error-prone, leading to security vulnerabilities and maintenance challenges in larger codebases.

This code generator offers a solution by utilizing Protobuf, a language-agnostic binary serialization format, to define structured data models for your HTML templates. By representing your template data in a structured manner using Protobuf messages, you can benefit from improved code maintainability and automated testing.

## Goals

- Provide a minimal layer on top of two rock-solid and stable foundations: Protobuf and Go's stdlib templating

## Specification

- MUST support specifying Go templates files next to your .proto
- MUST have the same name as the proto file, but with a .gossr extension
- MUST each Protobuf message X that has template named gossr_X is called a "partial"
- MUST for partial X: generate a public render method on X's go representation
- MUST for partial X: generate template funcs that allow the partial to be rendered from other templates
- MUST allow for importing gossr protobuf messages from other packages, repos and it should just work
- SHOULD work in a dev environment without having to recompile when templates change: i.e: re-parse on re-render
- SHOULD only generate the template func of message X when the c message's fields (in)directly reference the partial's message
- SHOULD for each partial: generate a test function that automatically fuzzes the render method: take care of one_ofs
- SHOULD for each partial fuzz test: generate the assertion that it is valid html (no open tags)
- SHOULD have a benchmark for a large tree of partials, need to make sure it is not too heavy on the memory
- COULD allow for runtime inspection to allow dynamic code based on the partials
- COULD have the generated code only be dependant on the std library
- COULD generate visual (regression) tests for example representations of each partial (corpus)
- COULD generate a http.Handler that renders each example in isolation (for visual testing)
- COULD allow fuzz and assert generation to be configurable/disabled per message
- COULD allow specifying a specific corpus per proto message field for testing (but maybe that's a different project)
- COULD generate a custom fuzz method on each partial to facilitate calling it in go fuzz: for oneOf fields specifically

## Research Tasks

- [x] MUST figure out what fuzzing technique we'll actually generate for partials, how to generate random protobuf messages?
  - https://github.com/brianvoe/gofakeit#struct
  - https://adalogics.com/blog/structure-aware-go-fuzzing-complex-types
  - https://github.com/flyingmutant/rapid
- MUST figure out what library we'll use for checking if the html is valid after rendering
- SHOULD figure out what's the best way to generate visual regression tests for partials

## Big Ideas

- What if protobuf options allow for defining how data is loaded for the partial, like graphql resolvers. But how to
  deal with parameters, waterfall, etc.

## Backlog

- [ ] MUST Implement a Protobuf generator that generates Fuzz() methods on message types to support oneOfFields (select random oneOf)
- [ ] COULD build this generator such that it provides Fuzz Funcs with a corpos that can be loaded: https://pkg.go.dev/github.com/google/gofuzz?utm_source=godoc#Fuzzer.Funcs maybe through a json file
