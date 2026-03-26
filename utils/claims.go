package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/logger"
	systemReq "github.com/irorikon/cloudgram-go/model/request"
)

// GetToken 从请求头获取token
func GetToken(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// 移除Bearer前缀
	if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:]
	}

	// 或者按空格分割
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	return authHeader
}

// GetClaims 从token解析声明
func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := GetToken(c)
	if token == "" {
		return nil, errors.New("未找到token")
	}

	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		logger.SysErrorF("从Authorization头解析JWT失败: %v, token: %s", err, maskToken(token))
		return nil, err
	}
	return claims, nil
}

// GetUserName 获取用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); exists {
		if waitUse, ok := claims.(*systemReq.CustomClaims); ok {
			return waitUse.Username
		}
	}

	// 如果上下文中没有，尝试从token解析
	claims, err := GetClaims(c)
	if err != nil {
		return ""
	}
	return claims.Username
}

// LoginToken 生成登录token
func LoginToken(username string) (string, systemReq.CustomClaims, error) {
	j := NewJWT()
	claims := j.CreateClaims(systemReq.BaseClaims{
		Username: username,
	})
	token, err := j.GenerateToken(claims)
	return token, claims, err
}

// maskToken 安全地显示token（只显示前10位和最后4位）
func maskToken(token string) string {
	if len(token) <= 14 {
		return "***"
	}
	return token[:10] + "..." + token[len(token)-4:]
}
