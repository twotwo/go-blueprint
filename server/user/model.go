package user

import (
	"time"

	"gorm.io/gorm"
)

// UserModel 定义用户在数据库中的表示
// 实现了 User API 模型到数据库模型的映射
type UserModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 用户基本信息
	Username   string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	FirstName  string `gorm:"size:50" json:"firstName"`
	LastName   string `gorm:"size:50" json:"lastName"`
	Email      string `gorm:"size:100" json:"email"`
	Password   string `gorm:"size:100;not null" json:"-"` // 密码不应在JSON中暴露
	Phone      string `gorm:"size:20" json:"phone"`
	UserStatus int32  `gorm:"default:0" json:"userStatus"` // 用户状态，0=非活跃, 1=活跃
}

// TableName 指定用户表名
func (UserModel) TableName() string {
	return "users"
}

// ToAPI 将数据库模型转换为API模型
func (u *UserModel) ToAPI() User {
	id := int64(u.ID)
	return User{
		Id:         &id,
		Username:   &u.Username,
		FirstName:  &u.FirstName,
		LastName:   &u.LastName,
		Email:      &u.Email,
		Phone:      &u.Phone,
		UserStatus: &u.UserStatus,
		// 注意：不返回密码
	}
}

// FromAPI 将API模型转换为数据库模型
func (u *UserModel) FromAPI(apiUser User) {
	if apiUser.Id != nil {
		u.ID = uint(*apiUser.Id)
	}
	if apiUser.Username != nil {
		u.Username = *apiUser.Username
	}
	if apiUser.FirstName != nil {
		u.FirstName = *apiUser.FirstName
	}
	if apiUser.LastName != nil {
		u.LastName = *apiUser.LastName
	}
	if apiUser.Email != nil {
		u.Email = *apiUser.Email
	}
	if apiUser.Password != nil {
		u.Password = *apiUser.Password // 在实际应用中应该先哈希处理
	}
	if apiUser.Phone != nil {
		u.Phone = *apiUser.Phone
	}
	if apiUser.UserStatus != nil {
		u.UserStatus = *apiUser.UserStatus
	}
}

// FindUserByUsername 根据用户名查找用户
func FindUserByUsername(db *gorm.DB, username string) (*UserModel, error) {
	var user UserModel
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser 创建新用户
func Create(db *gorm.DB, apiUser User) (*UserModel, error) {
	var user UserModel
	user.FromAPI(apiUser)

	// 在实际应用中应该进行密码哈希处理
	// user.Password = hashPassword(user.Password)

	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func Update(db *gorm.DB, username string, apiUser User) error {
	user, err := FindUserByUsername(db, username)
	if err != nil {
		return err
	}

	user.FromAPI(apiUser)
	return db.Save(user).Error
}

// DeleteUser 删除用户
func Delete(db *gorm.DB, username string) error {
	return db.Where("username = ?", username).Delete(&UserModel{}).Error
}

// ListUsers 获取所有用户
func ListUsers(db *gorm.DB) ([]UserModel, error) {
	var users []UserModel
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
