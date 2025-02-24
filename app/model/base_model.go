package model

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/twotwo/go-blueprint/app/global/variable"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

type BaseModel struct {
	*gorm.DB  `gorm:"-" json:"-"`
	Id        int64 `gorm:"primaryKey;comment:Primary Key" json:"id"`
	CreatedAt int   `gorm:"autoCreateTime;comment:Created Time" json:"created_at"`
	UpdatedAt int   `gorm:"autoCreateTime;comment:Updated Time" json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at"`   // 使用软删除功能，打开本行注释掉的代码即可；同时需要在数据库的所有表增加字段deleted_at 类型为 datetime
}

func initConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(variable.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db: " + err.Error())
	}
	setConnectionPool(db)

	return db
}

func setConnectionPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("failed to get db connection: %v", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(variable.ConnMaxIdleConns)
	sqlDB.SetMaxOpenConns(variable.ConnMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(variable.ConnMaxLifetime))
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(variable.ConnMaxIdleTime))
}

func init() {
	DBConn = initConnection()
}
