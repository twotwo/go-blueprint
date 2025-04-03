package model

import (
	"gorm.io/datatypes"
)

type UsersModel struct {
	BaseModel
	Username                 string `gorm:"default:null;comment:Username" json:"username"`
	PrimaryEmail             string `gorm:"default:null;comment:邮箱"`
	PrimaryPhone             string `gorm:"default:null;comment:电话"`
	PasswordEncrypted        string
	PasswordEncryptionMethod string
	Name                     string
	Avatar                   string
	Profile                  datatypes.JSON `gorm:"default:'{}'::jsonb"`
	ApplicationID            string
	IsSuspended              bool `gorm:"default:false"`
}

// 表名
func (u *UsersModel) TableName() string {
	return "users"
}
