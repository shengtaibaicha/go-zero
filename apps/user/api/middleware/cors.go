// middleware/headers.go
package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func HeadersMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 设置 CORS 头
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 暴露自定义响应头（多个头用逗号分隔）
			w.Header().Set("Access-Control-Expose-Headers", "redisKey, X-Custom-Header")

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// 继续处理请求
			next(w, r)
		}
	}
}
