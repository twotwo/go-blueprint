package user

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes 注册用户相关路由
func RegisterRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		// POST /user - 创建单个用户
		r.Post("/", CreateUser)

		// POST /user/createWithList - 批量创建用户
		r.Post("/createWithList", CreateUsersWithListInput)

		// GET /user/login - 用户登录
		r.Get("/login", LoginUser)

		// GET /user/logout - 用户登出
		r.Get("/logout", LogoutUser)

		// 针对特定用户名的操作
		r.Route("/{username}", func(r chi.Router) {
			// GET /user/{username} - 获取用户信息
			r.Get("/", GetUserByName)

			// PUT /user/{username} - 更新用户信息
			r.Put("/", UpdateUser)

			// DELETE /user/{username} - 删除用户
			r.Delete("/", DeleteUser)
		})
	})
}
