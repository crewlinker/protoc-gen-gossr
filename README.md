# protoc-gen-gossr
Fearless server-side rendering in Go using Protobuf.

## Introduction
This project aims to enhance the safety and simplicity of developing Go applications that utilize the html/template package by leveraging Protocol Buffers (Protobuf).

The html/template (or text/template) package in Go provides a powerful and flexible way to generate HTML, XML, and other text-based formats. However, working directly with strings and template syntax can be error-prone, leading to security vulnerabilities and maintenance challenges in larger codebases.

This code generator offers a solution by utilizing Protobuf, a language-agnostic binary serialization format, to define structured data models for your HTML templates. By representing your template data in a structured manner using Protobuf messages, you can benefit from improved code maintainability and automated testing.

## Specification
- MUST support specifying Go templates files next to your .proto 
- MUST have the same name as the proto file, but with a .gossr extension
- MUST each Protobuf message X that has template named gossr_X is called a "frame"
- MUST for frame X: generate a public render method on X's go representation
- MUST for frame X: generate template funcs that allow the frame to be rendered from other templates
- SHOULD only generate the template func of message X when the c message's fields (in)directly reference the frame's message
- SHOULD for each frame: generate a test function that automatically fuzzes the render method
- SHOULD for each frame fuzz test: generate the assertion that it is valid html (no open tags)
- COULD have the generated code only be dependant on the std library
- COULD generate visual (regression) tests for example representations of each frame (corpus)
- COULD generate a http.Handler that renders each example in isolation (for visual testing)