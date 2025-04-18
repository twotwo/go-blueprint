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

### For Go 1.24+

建议遵循 Go 1.24+ 中提供的 [go tool](https://www.jvt.me/posts/2025/01/27/go-tools-124/) 支持来管理 `oapi-codegen` 与核心应用程序的依赖关系。

```bash
$ go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
# this will then modify your `go.mod`
$ cat server/oapi/generate.go # 1.24 内置了对运行 tool 的缓存
package oapi

//go:generate go tool oapi-codegen -config cfg.yaml api.yaml
$ go generate -x ./...
# 等价于以下命令
$ cd docs && go tool oapi-codegen -config oapi-codegen.yaml ../api.yaml && cd -
```

## Generate code

[Spec 生成 Chi 代码](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#chi)

<https://github.com/oapi-codegen/oapi-codegen/tree/main/examples/minimal-server>

- chi/
- api.yaml

```bash
# generate chi server code
$ go generate -x ./...
# start server
$ go run cmd/oapi/main.go
```

## 访问服务接口

1. VSCode 安装 [vscode-openapi-viewer](https://marketplace.visualstudio.com/items?itemName=AndrewButson.vscode-openapi-viewer) 扩展
2. 打开 api.yaml 后右上出现绿色靶子图标
3. 点击后可以预览 API 并发送请求
