package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type AuthRouter struct {
}

func (a *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) {
	authRouter := Router.Group("auth").Use(middleware.JWTAuth())
	{
		authRouter.GET("refresh", systemAuthApi.Refresh)
	}
}
