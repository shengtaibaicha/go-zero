// middleware/headers.go
package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func HeadersMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 获取请求源
			origin := r.Header.Get("Origin")

			// 只对有Origin头的请求设置CORS（避免非跨域请求的问题）
			if origin != "" {
				// 开发环境可固定为前端地址，生产环境可根据需要限制
				w.Header().Set("Access-Control-Allow-Origin", origin)      // 携带Cookie的情况下不能为*
				w.Header().Set("Access-Control-Allow-Credentials", "true") // 携带Cookie
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Expose-Headers", "redisKey, X-Custom-Header") // 暴露给前端的头
				w.Header().Set("Access-Control-Max-Age", "86400")                            // 预检请求缓存24小时
			}

			// 专门处理OPTIONS请求
			if r.Method == http.MethodOptions {
				// 直接返回204 No Content，不调用next
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// 继续处理请求
			next(w, r)
		}
	}
}
