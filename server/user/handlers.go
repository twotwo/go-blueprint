package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/twotwo/go-blueprint/pkg/errors"
)

// DB 是一个全局数据库连接
// 在实际应用中，应该通过依赖注入或上下文来传递
var DB *gorm.DB

// CreateUser 处理创建用户的请求
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var apiUser User

	// 根据content-type解码请求体
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&apiUser); err != nil {
			apiErr := errors.BadRequest("无效的JSON格式")
			errors.WriteJSON(w, apiErr)
			return
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			apiErr := errors.BadRequest("无效的表单数据")
			errors.WriteJSON(w, apiErr)
			return
		}

		// 解析表单数据
		if r.Form.Get("username") != "" {
			username := r.Form.Get("username")
			apiUser.Username = &username
		}
		// 添加其他表单字段...
	}

	// 验证必要字段
	if apiUser.Username == nil || *apiUser.Username == "" {
		apiErr := errors.BadRequest("用户名不能为空")
		errors.WriteJSON(w, apiErr)
		return
	}

	if apiUser.Password == nil || *apiUser.Password == "" {
		apiErr := errors.BadRequest("密码不能为空")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 创建用户
	user, err := Create(DB, apiUser)
	if err != nil {
		// 检查是否是唯一性约束错误
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			apiErr := errors.BadRequest("用户名已存在")
			errors.WriteJSON(w, apiErr)
			return
		}

		apiErr := errors.InternalServer("创建用户失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 返回创建的用户
	apiResponse := user.ToAPI()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse)
}

// CreateUsersWithListInput 处理批量创建用户的请求
func CreateUsersWithListInput(w http.ResponseWriter, r *http.Request) {
	var apiUsers []User

	if err := json.NewDecoder(r.Body).Decode(&apiUsers); err != nil {
		apiErr := errors.BadRequest("无效的JSON格式")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 开始事务
	tx := DB.Begin()

	var createdUsers []UserModel

	for _, apiUser := range apiUsers {
		if apiUser.Username == nil || *apiUser.Username == "" {
			tx.Rollback()
			apiErr := errors.BadRequest("用户名不能为空")
			errors.WriteJSON(w, apiErr)
			return
		}

		if apiUser.Password == nil || *apiUser.Password == "" {
			tx.Rollback()
			apiErr := errors.BadRequest("密码不能为空")
			errors.WriteJSON(w, apiErr)
			return
		}

		user, err := Create(tx, apiUser)
		if err != nil {
			tx.Rollback()

			// 检查是否是唯一性约束错误
			if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
				apiErr := errors.BadRequest(fmt.Sprintf("用户名 %s 已存在", *apiUser.Username))
				errors.WriteJSON(w, apiErr)
				return
			}

			apiErr := errors.InternalServer("创建用户失败")
			errors.WriteJSON(w, apiErr)
			return
		}

		createdUsers = append(createdUsers, *user)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		apiErr := errors.InternalServer("创建用户失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 返回第一个创建的用户作为响应（按照API定义）
	if len(createdUsers) > 0 {
		apiResponse := createdUsers[0].ToAPI()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apiResponse)
		return
	}

	// 如果没有创建任何用户，返回空响应
	w.WriteHeader(http.StatusOK)
}

// LoginUser 处理用户登录
func LoginUser(w http.ResponseWriter, r *http.Request) {
	// 从查询参数获取用户名和密码
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 在数据库中查找用户
	var user UserModel
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		apiErr := errors.InternalServer("查询用户信息失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 验证密码（在实际应用中应该使用哈希比较）
	if user.Password != password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 设置响应头
	w.Header().Set("X-Rate-Limit", "100")
	expiresTime := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	w.Header().Set("X-Expires-After", expiresTime)

	// 返回登录令牌（在实际应用中应该生成JWT）
	token := "logged-in-token-would-go-here"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("\"%s\"", token)))
}

// LogoutUser 处理用户登出
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// 在实际应用中，可能需要使登录令牌失效
	w.WriteHeader(http.StatusOK)
}

// GetUserByName 根据用户名获取用户
func GetUserByName(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		apiErr := errors.BadRequest("用户名不能为空")
		errors.WriteJSON(w, apiErr)
		return
	}

	user, err := FindUserByUsername(DB, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiErr := errors.NotFound("用户不存在")
			errors.WriteJSON(w, apiErr)
			return
		}

		apiErr := errors.InternalServer("查询用户信息失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	apiResponse := user.ToAPI()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse)
}

// UpdateUser 更新用户信息
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		apiErr := errors.BadRequest("用户名不能为空")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 首先检查用户是否存在
	_, err := FindUserByUsername(DB, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiErr := errors.NotFound("用户不存在")
			errors.WriteJSON(w, apiErr)
			return
		}

		apiErr := errors.InternalServer("查询用户信息失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	var apiUser User

	// 根据content-type解码请求体
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&apiUser); err != nil {
			apiErr := errors.BadRequest("无效的JSON格式")
			errors.WriteJSON(w, apiErr)
			return
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			apiErr := errors.BadRequest("无效的表单数据")
			errors.WriteJSON(w, apiErr)
			return
		}

		// 解析表单数据
		if r.Form.Get("username") != "" {
			formUsername := r.Form.Get("username")
			apiUser.Username = &formUsername
		}
		// 添加其他表单字段...
	}

	// 更新用户
	err = Update(DB, username, apiUser)
	if err != nil {
		apiErr := errors.InternalServer("更新用户信息失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser 删除用户
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	if username == "" {
		apiErr := errors.BadRequest("用户名不能为空")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 首先检查用户是否存在
	_, err := FindUserByUsername(DB, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiErr := errors.NotFound("用户不存在")
			errors.WriteJSON(w, apiErr)
			return
		}

		apiErr := errors.InternalServer("查询用户信息失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 删除用户
	err = Delete(DB, username)
	if err != nil {
		apiErr := errors.InternalServer("删除用户失败")
		errors.WriteJSON(w, apiErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}
