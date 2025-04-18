package message

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes 注册消息相关路由
func RegisterRoutes(r chi.Router) {

	// 消息相关路由
	r.Route("/message", func(r chi.Router) {
		// POST /message - 创建新消息（需要JWT认证）
		r.With(AuthMiddleware).Post("/", CreateMessage)

		// GET /message/sms/{number} - 根据手机号查询短信
		r.Get("/sms/{number}", FindMessagesByNumber)

		// GET /message/sitemessage/{uid} - 根据用户ID查询站内消息
		r.Get("/sitemessage/{uid}", FindMessagesByUID)
	})
}

// AuthMiddleware 是一个简单的JWT认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// 检查Authorization头是否存在并且格式正确
		if token == "" || len(token) < 8 || token[:7] != "Bearer " {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Unauthorized: Missing or invalid token",
			})
			return
		}

		// 提取实际令牌部分（Bearer 之后的内容）
		tokenString := token[7:]

		// 这里应该添加实际的JWT验证逻辑
		// 简单示例仅做格式验证，实际应用中应验证签名等
		if len(tokenString) < 10 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Unauthorized: Invalid token format",
			})
			return
		}

		// 如果验证通过，继续处理请求
		next.ServeHTTP(w, r)
	})
}
