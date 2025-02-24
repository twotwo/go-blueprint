package variable

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	BasePath           string                     // 定义项目的根目录
	DateFormat         = "2025-02-17 15:04:05"    //  设置全局日期时间格式
	ResourceServerName = "https://www.gpufree.cn" // 资源服务器地址
	// 服务相关配置
	Port int    // 服务运行的端口号，从环境变量 PORT 加载，默认值为 8080
	Name string // 服务名称，从环境变量 NAME 加载，默认值为 example-service

	// 数据库相关配置，参见 https://gorm.io/zh_CN/docs/connecting_to_the_database.html#PostgreSQL
	DSN              string       // 数据库配置字符串：Data Source Name
	ConnMaxIdleConns int    = 1   // 连接池最大空闲连接数
	ConnMaxOpenConns int    = 5   // 连接池最大连接数
	ConnMaxLifetime  int    = 30  // 连接池连接的最大存活时间（秒）
	ConnMaxIdleTime  int    = 300 // 连接池连接的最大空闲时间（秒）

	// JWT 公钥文件路径
	PUBLIC_KEY_FILE string // 公钥文件路径
)

// envMap 存储环境变量 (忽略大小写)
var envMap map[string]string

// initEnvMap 初始化环境变量 map（key 统一转换为小写）
func initEnvMap() {
	envMap = make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2) // 拆分 "KEY=VALUE"
		if len(parts) == 2 {
			key := strings.ToLower(parts[0]) // 统一转小写
			envMap[key] = parts[1]
		}
	}
}

// getEnvIgnoreCase 进行大小写无关的环境变量查找
func getEnvIgnoreCase(key string, defaultValue string) string {
	lowerKey := strings.ToLower(key)
	if value, exists := envMap[lowerKey]; exists {
		return value
	}
	return defaultValue
}

func stringsToInt(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		fmt.Fprintln(os.Stdout, "current path:", curPath)
	} else {
		panic(fmt.Sprintf("failed to get current path: %v", err))
	}
	BasePath, _ = os.Getwd()

	// 初始化环境变量 map
	initEnvMap()

	// 加载 Port
	portStr := getEnvIgnoreCase("PORT", "8080")
	if portStr != "" {
		// 将字符串类型的端口号转换为整数.
		port, err := strconv.Atoi(portStr)
		if err != nil {
			// 如果转换失败，则返回错误.
			panic(fmt.Sprintf("failed to convert PORT to int: %v", err))
		}
		Port = port
	}
	Name = getEnvIgnoreCase("NAME", "inferring-api")
	DSN = getEnvIgnoreCase("DSN", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai")
	ConnMaxIdleConns = stringsToInt(getEnvIgnoreCase("CONN_MAX_IDLE_CONNS", "1"), 1)
	ConnMaxOpenConns = stringsToInt(getEnvIgnoreCase("CONN_MAX_OPEN_CONNS", "5"), 5)
	ConnMaxLifetime = stringsToInt(getEnvIgnoreCase("CONN_MAX_LIFETIME", "30"), 30)
	ConnMaxIdleTime = stringsToInt(getEnvIgnoreCase("CONN_MAX_IDLE_TIME", "300"), 300)

	PUBLIC_KEY_FILE = getEnvIgnoreCase("PUBLIC_KEY_FILE", "/opt/local/ide/workspaces/w2024/backend/gpufree-inferring-api/public_key.pem")

	if getEnvIgnoreCase("Resource_Server_Name", "") != "" {
		ResourceServerName = getEnvIgnoreCase("Resource_Server_Name", "")
	}
}
