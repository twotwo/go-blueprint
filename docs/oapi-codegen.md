# ogen Quick Started

[oapi-codegen/oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)

- 支持 OpenAPI v3
- 支持生成多种服务器端代码 [Supported Servers](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#generating-server-side-boilerplate)
- Github 6.9k stars

## Installation

```bash
# for the binary install
$ go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
$ oapi-codegen -version
github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
v2.4.1
```

## Generate code

[Spec 生成 Chi 代码](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#chi)

```bash
# generate chi server code
$ cd app/api && oapi-codegen --config=config.yaml ../../api.yaml && cd -
# start server
$ go run cmd/api/main.go
```
