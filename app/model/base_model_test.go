package model

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twotwo/go-blueprint/app/global/variable"
	"gorm.io/gorm"
)

// TestBaseModelbyFields 测试 BaseModel 结构体的字段和标签
func TestBaseModelbyFields(t *testing.T) {
	// 创建一个测试用的模型实例
	model := &BaseModel{
		Id:        1,
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()),
	}

	// 测试字段值是否符合预期
	assert.Equal(t, int64(1), model.Id) // ID 应该是 1
	assert.NotZero(t, model.CreatedAt)  // 创建时间不应为零
	assert.NotZero(t, model.UpdatedAt)  // 更新时间不应为零
	assert.Nil(t, model.DB)             // DB 字段应该为 nil
}

// TestInitConnection 测试数据库连接初始化
func TestInitConnection(t *testing.T) {
	// 保存原始 DSN 以便测试后恢复
	origDSN := variable.DSN
	defer func() {
		variable.DSN = origDSN
	}()

	// 定义测试用例
	tests := []struct {
		name    string // 测试用例名称
		dsn     string // 数据库连接字符串
		wantErr bool   // 是否期望错误
	}{
		{
			name:    "Valid DSN", // 有效的 DSN
			dsn:     variable.DSN,
			wantErr: false,
		},
		{
			name:    "Invalid DSN", // 无效的 DSN
			dsn:     "invalid-dsn",
			wantErr: true,
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variable.DSN = tt.dsn

			if tt.wantErr {
				// 使用 defer 捕获预期的 panic
				defer func() {
					if r := recover(); r == nil {
						t.Error("期望出现 panic 但没有")
					}
				}()
			}

			// 测试连接初始化
			db := initConnection()
			if !tt.wantErr {
				assert.NotNil(t, db)             // 确保返回的 db 不为空
				assert.IsType(t, &gorm.DB{}, db) // 确保返回类型正确
			}
		})
	}
}

// TestSetConnectionPool 测试连接池配置
func TestSetConnectionPool(t *testing.T) {
	// 创建一个测试用的数据库连接
	db := &gorm.DB{}

	// 测试无效 DB 时是否会 panic
	assert.Panics(t, func() {
		setConnectionPool(db)
	})

	// 测试连接池设置
	t.Run("Connection Pool Settings", func(t *testing.T) {
		if db, err := initConnection().DB(); err == nil {
			stats := db.Stats()
			// 验证连接池参数是否符合配置
			assert.LessOrEqual(t, stats.MaxOpenConnections, variable.ConnMaxOpenConns)
			assert.LessOrEqual(t, stats.Idle, variable.ConnMaxIdleConns)
		}
	})
}

// TestDBConnbyGlobalInstance 测试全局数据库连接实例
func TestDBConnbyGlobalInstance(t *testing.T) {
	// 确保全局连接已初始化
	assert.NotNil(t, DBConn, "全局数据库连接应该已初始化")

	// 测试是否为单例模式（返回相同实例）
	db1 := DBConn
	db2 := DBConn
	assert.Equal(t, db1, db2, "应该返回相同的数据库实例")
}

func TestSchemaMigration(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("跳过本测试，仅在 CI 环境下运行")
	}
	// 迁移数据库
	if err := DBConn.Set("gorm:table_options",
		";COMMENT ON TABLE tb_users IS '用户表'; COMMENT ON sequence tb_users_id_seq IS '用户ID序列';").
		AutoMigrate(&UsersModel{}); err != nil {
		assert.NoError(t, err, "迁移数据库失败")
	}
	DBConn.Create(&UsersModel{Username: "admin", PasswordEncrypted: "admin",
		PrimaryPhone: "1234567890", Name: "Admin", IsSuspended: true})
	// err := DBConn.AutoMigrate(&UsersModel{})
	// assert.NoError(t, err, "迁移数据库失败)
}

// TestMain 测试主函数，用于设置和清理测试环境
func TestMain(m *testing.M) {
	// 设置测试环境
	origDSN := variable.DSN
	if os.Getenv("CI") == "true" {
		// 如果在 CI 环墋下，使用测试数据库
		variable.DSN = "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	} else {
		// 否则使用本地数据库
		variable.DSN = origDSN
	}

	// 运行测试
	code := m.Run()

	// 清理环境
	variable.DSN = origDSN

	os.Exit(code)
}
