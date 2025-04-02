# API 设计最佳实践(Go 语言版)

API 优先（API-First）是一种以 API 设计为核心的开发理念，强调在系统构建初期优先定义接口规范，再围绕 API 实现功能。

采用 OpenAPI 作为系统间的集成方式。标准化的 OpenAPI 更有利于系统间的集成，因为 OpenAPI 有明确的契约描述或接口规格描述，且提供了各种开放的工具，可以用来做 IoT（连通性测试）、SIT（集成测试）等。同时，由于其开放接口（比如，基于 RESTful）的特性，可以实现快速集成，从而提升集成的效率。

声明式 API。很多软件交付都是“告诉”系统需要做什么，特别是脚本中往往会写明如何进行部署；而声明式 API 首先是“告诉”系统期望的目标状态是什么，比如，在这种环境下部署需要用到两个实例，其次才是脚本或工具需要做什么才能交付这个目标状态（即如何做）。声明式 API 本身并不复杂，实际上它是一种开发理念的彻底升级，因为系统更多的是关注需要什么（达到什么状态），所有的“如何做”都是围绕这个目标状态来服务的。

## 设计原则

1. 接口先行
    在编写业务逻辑前完成 API 设计，确保接口规范成为系统构建的基础。

2. 契约驱动
    通过 OpenAPI 等规范定义请求/响应格式，使前后端团队可并行开发（Swagger 工具支持 YAML 编写 API 契约）。

3. 声明式 API
    告诉”系统期望的目标状态是什么，而不是需要做什么。

4. 安全性增强
    API 网关可集中管理身份验证（如 OAuth 2.0/ OIDC）、流量监控和攻击防护。

## 实施工具推荐

1. **Apifox** 2025年首选工具，覆盖API全生命周期，支持智能文档生成和实时协作（免费版满足初创团队需求）。
2. **Swagger Editor** 通过 YAML 编写 OpenAPI 规范，自动生成交互式文档和客户端 SDK。
3. **Mockoon** 支持 OpenAPI v3 Spec 的导入导出，启动 mock server。
4. **oapi-codegen** 一个命令行工具，能够将 OpenAPI 规范转换为 Go 代码

### oapi-codegen

<https://github.com/oapi-codegen/oapi-codegen>

支持 OpenAPI v3，支持生成多种服务器端代码 [Supported Servers](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#generating-server-side-boilerplate)

```bash
# for the binary install
$ go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
$ oapi-codegen -version
github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
v2.4.1
```

[Spec 生成 Chi 代码](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#chi)

```bash
# generate chi server code
$ cd app/api && oapi-codegen --config=config.yaml ../../api.yaml && cd -
# start server
$ go run cmd/api/main.go
```

`$ oapi-codegen --generate chi-server -o petstore.gen.go -package petstore api/swagger.yaml`

- `--generate chi-server`：​指定生成 `chi` 路由器的服务端代码
- `-o petstore.gen.go`：​指定输出文件名
- `-package petstore`：​指定生成代码的包名。​
- `swagger.yaml`：​您的 OpenAPI 规范文件。
