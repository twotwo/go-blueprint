package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/twotwo/go-blueprint/server/message"
	"github.com/twotwo/go-blueprint/server/user"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r chi.Router) {
	// API版本前缀
	r.Route("/api/v1", func(r chi.Router) {
		// 注册用户资源路由
		user.RegisterRoutes(r)

		// 注册消息资源路由
		message.RegisterRoutes(r)
	})
}
