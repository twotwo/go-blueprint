package database

import (
	"fmt"
	"log"
	"time"

	"github.com/twotwo/go-blueprint/pkg/variables"
	"github.com/twotwo/go-blueprint/server/user"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Setup 初始化数据库连接并返回GORM DB实例
func Setup() (*gorm.DB, error) {
	dbType := variables.GetEnv("DB_TYPE", "sqlite")

	var dialector gorm.Dialector

	switch dbType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			variables.GetEnv("DB_USER", "root"),
			variables.GetEnv("DB_PASSWORD", ""),
			variables.GetEnv("DB_HOST", "localhost"),
			variables.GetEnv("DB_PORT", "3306"),
			variables.GetEnv("DB_NAME", "blueprint"),
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			variables.GetEnv("DB_HOST", "localhost"),
			variables.GetEnv("DB_PORT", "5432"),
			variables.GetEnv("DB_USER", "postgres"),
			variables.GetEnv("DB_PASSWORD", ""),
			variables.GetEnv("DB_NAME", "blueprint"),
		)
		dialector = postgres.Open(dsn)
	default: // sqlite
		dialector = sqlite.Open(variables.GetEnv("DB_FILE", "blueprint.db"))
	}

	// 配置GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(variables.GetEnvInt("DB_MAX_IDLE_CONNS", 10))

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(variables.GetEnvInt("DB_MAX_OPEN_CONNS", 100))

	// 设置连接最大生存时间
	sqlDB.SetConnMaxLifetime(time.Duration(variables.GetEnvInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second)

	// 自动迁移模式
	if variables.GetEnvBool("DB_AUTO_MIGRATE", true) {
		log.Println("正在自动迁移数据库模式...")
		err = db.AutoMigrate(
			&user.UserModel{},
			// 添加其他需要迁移的模型
		)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
