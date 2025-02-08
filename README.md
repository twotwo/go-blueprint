# Project rest_api

用 [go-blueprint](github.com/melkeydev/go-blueprint) 构建的 RESTful API.

## VSCode

对 Go 插件来说，GOROOT 环境变量还是必须的

## Getting Started

### 添加环境变量

```bash
$ cat .env
PORT=8080
BLUEPRINT_DB_HOST=psql_bp
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=blueprint
BLUEPRINT_DB_USERNAME=melkey
BLUEPRINT_DB_PASSWORD=password1234
BLUEPRINT_DB_SCHEMA=web
```

### 安装依赖

`go mod tidy`

### 运行

`make run`

## 单元测试

### 接口测试

[testing chi](https://go-chi.io/#/pages/testing) chi 路由测试

## API 文档

### swaggo/swag

<https://github.com/swaggo/swag>

与 chi 集成 [swaggo/http-swagger/v2](https://github.com/swaggo/http-swagger/)

按照 [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format) 在 API 代码中添加注释(internal/server/routes.go)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g internal/server/routes.go # 生成 docs，指定路由配置文件(默认是 main.go)
```

[API操作](https://github.com/swaggo/swag/blob/master/README_zh-CN.md#api%E6%93%8D%E4%BD%9C) 重要注释

- `description` 操作行为的详细说明。
- `summary` 该操作的简短摘要
- `tags` 每个API操作的标签列表，以逗号分隔
- `router` 以空格分隔的路径定义。 `path,[httpMethod]`
