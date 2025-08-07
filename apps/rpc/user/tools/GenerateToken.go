package tools

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtCustomClaims struct {
	UserId             string `json:"userId"`   // 自定义字段：用户ID
	Username           string `json:"username"` // 自定义字段：用户名
	Role               string `json:"role"`     // 自定义字段：用户角色
	jwt.StandardClaims        // 标准声明，包含过期时间等信息
}

// 生成JWT令牌
func GenerateToken(secretKey, userId, username, role string, expireHours int64) (string, error) {
	// 1. 定义过期时间（当前时间+有效期）
	expireTime := time.Now().Add(time.Hour * time.Duration(expireHours)).Unix()

	// 2. 创建声明对象（必须使用与解析时相同的结构体）
	claims := JwtCustomClaims{
		UserId:   userId,   // 自定义字段：用户ID
		Username: username, // 自定义字段：用户名
		Role:     role,     // 自定义字段：用户角色
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,        // 标准字段：过期时间
			IssuedAt:  time.Now().Unix(), // 标准字段：签发时间
			Issuer:    "baicha",          // 标准字段：签发者（可选）
		},
	}

	// 3. 使用HS256算法签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. 生成令牌字符串
	return token.SignedString([]byte(secretKey))
}
