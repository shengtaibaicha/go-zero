package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
	"time"
)

// 定义上下文存储的键常量，用于从上下文获取数据时的标识
// 这些键与JWT中解析出的用户信息字段对应
const (
	JwtUserIdKey   = "userId"   // 用户ID在上下文中的存储键
	JwtUsernameKey = "username" // 用户名在上下文中的存储键
	JwtRoleKey     = "role"     // 用户角色在上下文中的存储键
)

// JwtAuthConfig 是JWT中间件的配置结构体
// 用于传递中间件所需的核心参数
type JwtAuthConfig struct {
	SecretKey    string   // JWT签名密钥，用于验证令牌的合法性
	ExcludePaths []string // 不需要进行JWT验证的路径列表
	TokenPrefix  string   // 令牌前缀，通常为"Bearer"，与请求头中的格式对应
}

// JwtCustomClaims 自定义JWT声明结构体
// 用于解析JWT中包含的用户信息，需要继承jwt.StandardClaims以支持标准字段
type JwtCustomClaims struct {
	UserId             string `json:"userId"`   // 自定义字段：用户ID
	Username           string `json:"username"` // 自定义字段：用户名
	Role               string `json:"role"`     // 自定义字段：用户角色
	jwt.StandardClaims        // 标准声明，包含过期时间等信息
}

// NewJwtAuthMiddleware 创建一个JWT认证中间件
// 参数：cfg - JWT中间件的配置
// 返回：一个符合GoZero规范的HTTP中间件函数
func NewJwtAuthMiddleware(cfg JwtAuthConfig) func(http.HandlerFunc) http.HandlerFunc {
	// 如果未指定令牌前缀，默认使用"Bearer"（这是OAuth2.0的标准前缀）
	//if cfg.TokenPrefix == "" {
	//	cfg.TokenPrefix = "Bearer"
	//}

	// 返回中间件函数，遵循GoZero的中间件签名：func(http.HandlerFunc) http.HandlerFunc
	return func(next http.HandlerFunc) http.HandlerFunc {
		// 中间件的核心逻辑，接收http.ResponseWriter和*http.Request
		return func(w http.ResponseWriter, r *http.Request) {
			// 1. 检查当前请求路径是否在排除列表中
			// 如果是排除路径，直接跳过验证，执行下一个处理器
			if isExcluded(r.URL.Path, cfg.ExcludePaths) {
				next(w, r)
				return
			}

			// 2. 从请求头中获取Authorization字段
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// 日志记录：Authorization头为空的错误
				logx.Errorf("Authorization header is empty, path: %s", r.URL.Path)
				// 使用GoZero的httpx工具返回JSON格式的错误响应
				httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
					"code":    "401",
					"message": "Authorization header is required",
				})
				return // 终止请求处理
			}

			// 3. 验证Authorization头的格式是否正确
			// 标准格式为："Bearer <token>"，使用空格分割为两部分
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != cfg.TokenPrefix {
				logx.Errorf("Invalid authorization format, path: %s", r.URL.Path)
				httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
					"code":    "401",
					"message": fmt.Sprintf("Authorization format must be '%s <token>'", cfg.TokenPrefix),
				})
				return
			}

			// 4. 验证JWT令牌的有效性
			// parts[1]是提取出的纯令牌字符串
			claims, err := verifyJwtToken(parts[1], cfg.SecretKey)
			if err != nil {
				logx.Errorf("JWT verify failed: %v, path: %s", err, r.URL.Path)
				httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
					"code":    "401",
					"message": "Invalid or expired token",
				})
				return
			}

			// 5. 将解析后的用户信息存入GoZero上下文
			// 获取原始请求上下文
			ctx := r.Context()
			// 使用GoZero的contextx.WithValue存储键值对（线程安全）
			ctx = context.WithValue(ctx, JwtUserIdKey, claims.UserId)
			ctx = context.WithValue(ctx, JwtUsernameKey, claims.Username)
			ctx = context.WithValue(ctx, JwtRoleKey, claims.Role)

			// 更新请求的上下文为新的上下文
			r = r.WithContext(ctx)

			// 6. 执行下一个处理器（业务逻辑）
			next(w, r)
		}
	}
}

// verifyJwtToken 验证JWT令牌的有效性并解析出声明信息
// 参数：
//
//	tokenString - 待验证的JWT令牌字符串
//	secret - 签名密钥（与生成令牌时使用的密钥一致）
//
// 返回：
//
//	*JwtCustomClaims - 解析后的自定义声明
//	error - 验证过程中出现的错误（如签名无效、令牌过期等）
func verifyJwtToken(tokenString, secret string) (*JwtCustomClaims, error) {
	// 解析令牌：使用jwt.ParseWithClaims方法，传入自定义声明结构体的指针
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否为预期的HMAC算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回签名密钥（用于验证签名）
		return []byte(secret), nil
	})

	// 如果解析过程出错（如格式错误、签名验证失败等）
	if err != nil {
		return nil, err
	}

	// 验证令牌是否有效，并提取声明信息
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		// 检查令牌是否已过期
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, errors.New("token expired")
		}
		// 返回解析后的声明
		return claims, nil
	}

	// 令牌无效（如未通过验证）
	return nil, errors.New("invalid token")
}

// isExcluded 检查当前请求路径是否需要排除JWT验证
// 参数：
//
//	path - 当前请求的路径
//	excludePaths - 配置的排除路径列表
//
// 返回：
//
//	bool - true表示需要排除，false表示需要验证
func isExcluded(path string, excludePaths []string) bool {
	for _, p := range excludePaths {
		// 支持两种匹配模式：
		// 1. 精确匹配：路径完全相等（如"/login"匹配"/login"）
		// 2. 前缀匹配：排除路径以"/"结尾，且当前路径以此为前缀（如"/public/"匹配"/public/data"）
		if p == path || (strings.HasSuffix(p, "/") && strings.HasPrefix(path, p)) {
			return true
		}
	}
	return false
}

// GetUserIdFromCtx 从GoZero上下文中获取用户ID
// 参数：
//
//	ctx - 上下文对象
//
// 返回：
//
//	string - 用户ID（如果不存在或类型错误，返回空字符串）
func GetUserIdFromCtx(ctx context.Context) string {
	// 从上下文获取值，并进行类型断言为string
	// 使用空白标识符忽略断言失败的错误（在正常流程中不会失败）
	val, _ := ctx.Value(JwtUserIdKey).(string)
	return val
}

// GetUsernameFromCtx 从GoZero上下文中获取用户名
// 参数：
//
//	ctx - 上下文对象
//
// 返回：
//
//	string - 用户名（如果不存在或类型错误，返回空字符串）
func GetUsernameFromCtx(ctx context.Context) string {
	val, _ := ctx.Value(JwtUsernameKey).(string)
	return val
}

// GetRoleFromCtx 从GoZero上下文中获取用户角色
// 参数：
//
//	ctx - 上下文对象
//
// 返回：
//
//	string - 用户角色（如果不存在或类型错误，返回空字符串）
func GetRoleFromCtx(ctx context.Context) string {
	val, _ := ctx.Value(JwtRoleKey).(string)
	return val
}
