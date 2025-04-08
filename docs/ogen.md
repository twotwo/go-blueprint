# ogen Quick Started

[ogen-go/ogen](https://github.com/ogen-go/ogen) 是一个 OpenAPI v3 代码生成器，专为 Go 语言设计，Github 1.7k stars。

参见： <https://ogen.dev/docs/intro/>

## Installation

`go install github.com/ogen-go/ogen/cmd/ogen@latest`

## Generate code

见 [generate.go](../examples/generate.go)

`go generate ./...`

## Using generated server

`examples/petstore/oas_server_gen.go` 实现其中的接口

## Run server

`go run examples/main.go`
