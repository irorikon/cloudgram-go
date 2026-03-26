package router

import "github.com/gin-gonic/gin"

type BaseRouter struct {
}

func (b *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	{
		baseRouter.POST("login", systemAuthApi.Login)
		baseRouter.GET("health", systemApi.Health)
	}
}
