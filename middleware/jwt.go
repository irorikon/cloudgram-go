package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model/response"
	"github.com/irorikon/cloudgram-go/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 头获取 Bearer token
		token := utils.GetToken(c)
		if token == "" {
			response.NoAuth("未登录或非法访问，请登录", c)
			c.Abort()
			return
		}
		if isBlacklist(token) {
			response.NoAuth("您的帐户异地登陆或令牌失效", c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.ErrTokenExpired) {
				response.NoAuth("登录已过期，请重新登录", c)
				c.Abort()
				return
			}
			response.NoAuth(err.Error(), c)
			c.Abort()
			return
		}

		// 设置claims到上下文
		c.Set("claims", claims)
		c.Next()
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func isBlacklist(jwt string) bool {
	_, ok := config.BlackCache.Get(jwt)
	return ok
}
