package server

import (
	"crypto/subtle"
	"fmt"
	"net/http"
)

// BasicAuth 实现了一个简单的中间件，用于为路由添加 HTTP 基础认证功能。
// 参数说明：
//
//	realm - 用于指定认证领域（浏览器弹出登录框显示）
//	creds - 包含用户名和密码映射的字典，用于对比认证信息
//
// 返回一个中间件，该中间件会在请求进入下一个处理器之前执行认证检查。
func BasicAuth(realm string, creds map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 尝试从请求中解析 Basic Auth 信息
			user, pass, ok := r.BasicAuth()
			if !ok {
				// 未提供认证信息，认证失败
				basicAuthFailed(w, realm)
				return
			}

			// 根据用户名查找预设的密码
			credPass, credUserOk := creds[user]
			// 使用常量时间比较函数对比密码，防止计时攻击
			if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				// 用户名不存在或密码不正确，认证失败
				basicAuthFailed(w, realm)
				return
			}

			// 认证成功，继续执行下一个处理器
			next.ServeHTTP(w, r)
		})
	}
}

// basicAuthFailed 用于返回认证失败的响应，提示客户端进行基础认证。
// 该函数设置了 WWW-Authenticate 响应头，并发送 401 未授权状态码。
func basicAuthFailed(w http.ResponseWriter, realm string) {
	// 设置认证提示信息，浏览器收到该响应会提示用户输入用户名和密码
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	// 返回 401 状态码，表示请求需要认证
	w.WriteHeader(http.StatusUnauthorized)
}
