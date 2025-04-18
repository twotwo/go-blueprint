package errors

import (
	"encoding/json"
	"net/http"
)

// APIError 表示API错误响应
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现error接口
func (e APIError) Error() string {
	return e.Message
}

// New 创建一个新的API错误
func New(code int, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

// NotFound 返回404错误
func NotFound(message string) APIError {
	if message == "" {
		message = "资源未找到"
	}
	return New(http.StatusNotFound, message)
}

// BadRequest 返回400错误
func BadRequest(message string) APIError {
	if message == "" {
		message = "请求参数错误"
	}
	return New(http.StatusBadRequest, message)
}

// Unauthorized 返回401错误
func Unauthorized(message string) APIError {
	if message == "" {
		message = "未授权的请求"
	}
	return New(http.StatusUnauthorized, message)
}

// InternalServer 返回500错误
func InternalServer(message string) APIError {
	if message == "" {
		message = "服务器内部错误"
	}
	return New(http.StatusInternalServerError, message)
}

// WriteJSON 将错误写入HTTP响应
func WriteJSON(w http.ResponseWriter, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
