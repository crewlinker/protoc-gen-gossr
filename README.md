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
- SHOULD provide a way to declare and provide "context" data across templates in a hierarchy, request info, logged-in user etc
- SHOULD work in a dev environment without having to recompile when templates change: i.e: re-parse on re-render
- SHOULD only generate the template func of message X when the c message's fields (in)directly reference the partial's message
- SHOULD for each partial: generate a test function that automatically fuzzes the render method: take care of one_ofs
- SHOULD for each partial fuzz test: generate the assertion that it is valid html (no open tags)
- SHOULD have a benchmark for a large tree of partials, need to make sure it is not too heavy on the memory
- COULD have tooling that instantly shows feedback of changing the html, with styling, and across examples (e.g like a storybook)
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
- [x] MUST figure out what library we'll use for checking if the html is valid after rendering
  - Parse using the encoding/xml, as described here: https://stackoverflow.com/a/52410528
- [ ] MUST figure out how partials (or sub-sections) of the tree can overwrite higher values up in the tree. Such as
      meta tags, script/inport overwrites. etc.
  - If we have the concept of a "page" the page rendering method might be able to ovewrite this.
- [x] MUST figure out and design the best way is to declare and provide contextual data across parials in a hierarchy
  - IDEA2: Reference fields, shared context needs to include them all
    - PRO: context arguments can be narrow interfaces
    - PRO: very Go-like, protobuf messages already define method for their fields
    - PRO: Would be nice that we don't need to add a pre-processor to make render func ergonomic
    - PRO: No extra memory that needs to be initialized, just pass the context
    - CON: Cannot sure how to reference nested values in the context Context.foo.bar
    - CON: How to deal with multiple context messages, can only be one? per package?
    - CON: how to fuzz them? Fuzzing interfaces is hard
      - Custom fuzz func that just populates the fields for the context
    - CON: how to define them?
      - Options in format: ContextMessage.field
      - Statically check that the mentioned context message indeed as the declared fields
  - IDEA3: Reference a message, context arg is a union of the fields
    - CON: small messages/structs being spread around
    - PRO: No need to come up with the go type for the generated interface type has always been tricky
  - IDEA4: Options on message field, generated renders will read the value from the struct, into the message
    - CON: Not clear that they are populated through a side channel (not really props)
    - CON: lot of duplication if a context value is used a lot
    - PRO: might be able to use the "fieldmask" type and pkg to union and validate
    - PRO: No extra type/interface to fuzz
    - PRO: Can reference nested fields on the context
- [ ] COULD provide linting of HTML
  - Maybe with: https://github.com/wawandco/milo

## Big Ideas

- What if protobuf options allow for defining how data is loaded for the partial, like graphql resolvers. But how to
  deal with parameters, waterfall, etc.

## Backlog

- [ ] MUST Implement a Protobuf generator that generates Fuzz() methods on message types to support oneOfFields (select random oneOf)
- [ ] COULD build this generator such that it provides Fuzz Funcs with a corpos that can be loaded: https://pkg.go.dev/github.com/google/gofuzz?utm_source=godoc#Fuzzer.Funcs maybe through a json file
- [ ] COULD design a pre-processing step to make rendering partials from templates more ergonomic, especially with
      context values needing to be passed around
- [ ] COULD support generating (fuzzing)tests for different test runners, e.g: standard lib, and Ginkgo
