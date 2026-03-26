package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
)

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization, Token")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// 按照配置参数放行
func CorsWithConfig() gin.HandlerFunc {
	// 放行全部
	if config.Debug {
		return Cors()
	}

	return func(c *gin.Context) {
		whitelist := checkCORS(c.Request.Header.Get("Origin"))
		if whitelist != nil {
			c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
			c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)
			c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)
			c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
			if whitelist.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		// 非严格白名单模式，无论是否通过检查均放行所有 OPTIONS 方法
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func checkCORS(currentOrigin string) *config.WhiteListConfig {
	for _, whitelist := range config.GlobalCorsConfig {
		if whitelist.AllowOrigin == currentOrigin {
			return &whitelist
		}
	}
	return nil
}
