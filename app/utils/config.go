package utils

import (
	"errors"
	"os"
	"strconv"
	"sync"
)

// Config holds the application's configuration.
// 该结构体用于保存应用程序的配置信息，包括服务配置和数据库配置.
type Config struct {
	// 服务相关配置
	ServicePort int    `env:"SERVICE_PORT" default:"8080"`            // 服务运行的端口号，从环境变量 SERVICE_PORT 加载，默认值为 8080
	ServiceName string `env:"SERVICE_NAME" default:"example-service"` // 服务名称，从环境变量 SERVICE_NAME 加载，默认值为 example-service

	// 数据库相关配置
	DBHost     string `env:"DB_HOST" default:"localhost"`    // 数据库主机名，从环境变量 DB_HOST 加载，默认值为 localhost
	DBPort     int    `env:"DB_PORT" default:"5432"`         // 数据库端口号，从环境变量 DB_PORT 加载，默认值为 5432
	DBUser     string `env:"DB_USER" default:"postgres"`     // 数据库用户名，从环境变量 DB_USER 加载，默认值为 postgres
	DBPassword string `env:"DB_PASSWORD" default:"password"` // 数据库密码，从环境变量 DB_PASSWORD 加载，默认值为 password
	DBName     string `env:"DB_NAME" default:"example_db"`   // 数据库名称，从环境变量 DB_NAME 加载，默认值为 example_db

	// 其他配置项按需添加
}

var (
	cfg  *Config   // 全局配置实例
	once sync.Once // 用于确保配置只加载一次
)

// GetInstance returns the global configuration instance.
// GetInstance 返回全局唯一的配置实例，如果还未加载则进行加载.
// 注意：该函数采用 sync.Once 确保 load() 仅执行一次.
func GetConfigInstance() (*Config, error) {
	once.Do(func() {
		cfg = new(Config)
		// 当加载配置出错时，直接 panic 退出程序，并输出错误信息.
		err := cfg.load()
		if err != nil {
			panic("failed to load configuration: " + err.Error())
		}
	})
	return cfg, nil
}

// getEnv 获取环境变量值，如果为空则使用默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// load loads environment variables into the Config struct.
func (c *Config) load() error {
	// 加载 ServicePort
	portStr := getEnv("SERVICE_PORT", "8080")
	if portStr != "" {
		// 将字符串类型的端口号转换为整数.
		port, err := strconv.Atoi(portStr)
		if err != nil {
			// 如果转换失败，则返回错误.
			return errors.New("invalid SERVICE_PORT value")
		}
		c.ServicePort = port
	}

	// Load other fields similarly...
	// 例如: ServiceName, DBHost, DBPort, DBUser, DBPassword, DBName 等
	c.ServiceName = getEnv("SERVICE_NAME", "rest-api")
	c.DBHost = getEnv("DB_HOST", "psql_bp")
	c.DBPort, _ = strconv.Atoi(getEnv("DB_PORT", "5432"))
	c.DBUser = getEnv("DB_USER", "melkey")
	c.DBPassword = getEnv("DB_PASSWORD", "password")
	c.DBName = getEnv("DB_NAME", "blueprint")

	return nil
}
