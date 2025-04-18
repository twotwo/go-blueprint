# Project rest_api

~~用 [go-blueprint](github.com/melkeydev/go-blueprint) 构建的 RESTful API.~~

一个 RESTful API 服务底座，可以快速嵌入新的 Resource 接口。

采用 `chi` 做路由管理，目录结构如下

```bash
./go-blueprint
├── cmd                   # 应用程序入口点
│   ├── oapi/main.go      # oapi-codegen 示例
│   └── server/main.go    # RESTful 服务
├── deploy
├── docs
├── pkg                   # 通用功能包
│   ├── errors
│   └── variables
└── server                # Resources of API
    ├── message
    ├── oapi
    └── user

```

## VSCode

对 Go 插件来说，GOROOT 环境变量还是必须的

安装 Go 语言插件：`golang.go` 扩展

安装 [vscode-openapi-viewer](https://marketplace.visualstudio.com/items?itemName=AndrewButson.vscode-openapi-viewer) 扩展

## Getting Started

### 如何添加新资源

要添加新的资源接口，可以参照 `user` 和 `message` 资源的实现模式：

1. 在 server 目录下创建新的资源目录，例如 product
2. 在该目录下创建 routes.go 和 handlers.go 文件
3. 实现资源的路由注册和处理函数
4. 在 server/routes.go 中导入新模块并注册新资源的路由

更多说明，见 [oapi-codegen.md](./docs/oapi-codegen.md)

### 安装依赖

`go mod tidy`

### 运行

`make run`

## 单元测试

### 接口测试

[testing chi](https://go-chi.io/#/pages/testing) chi 路由测试

## API 文档

### ~~swaggo/swag~~ 仅支持 OpenAPI v2

<https://github.com/swaggo/swag>

与 chi 集成 [swaggo/http-swagger/v2](https://github.com/swaggo/http-swagger/)

按照 [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format) 在 API 代码中添加注释(app/server/routes.go)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g app/server/routes.go # 生成 docs，指定路由配置文件(默认是 main.go)
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
swagger serve -F=swagger docs/swagger.yaml # 启动 spec 服务
```

[API操作](https://github.com/swaggo/swag/blob/master/README_zh-CN.md#api%E6%93%8D%E4%BD%9C) 重要注释

- `description` 操作行为的详细说明。
- `summary` 该操作的简短摘要
- `tags` 每个API操作的标签列表，以逗号分隔
- `router` 以空格分隔的路径定义。 `path,[httpMethod]`

