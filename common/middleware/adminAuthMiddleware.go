package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 定义角色常量
const (
	Admin      = "admin"
	SuperAdmin = "superAdmin"
)

// AdminAuthMiddleware 管理员权限验证中间件
func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从上下文获取用户角色（依赖JWT中间件已执行）
		// 注意：必须在JWT中间件之后使用此中间件，否则无法获取role
		role := GetRoleFromCtx(r.Context())
		if role == "" {
			logx.Errorf("管理员权限验证失败：上下文未获取到用户角色，path: %s", r.URL.Path)
			httpx.WriteJson(w, http.StatusForbidden, map[string]string{
				"code":    "403",
				"message": "权限验证失败：用户角色信息缺失",
			})
			return
		}

		// 检查角色是否为管理员
		if role != Admin && role != SuperAdmin {
			logx.Errorf("管理员权限验证失败：用户角色为 %s，path: %s", role, r.URL.Path)
			httpx.WriteJson(w, http.StatusForbidden, map[string]string{
				"code":    "403",
				"message": "权限不足：需要管理员权限",
			})
			return
		}

		// 权限验证通过，继续执行下一个处理器
		next.ServeHTTP(w, r)
	}
}
