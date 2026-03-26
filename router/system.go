package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type SystemRouter struct {
}

func (s *SystemRouter) InitSystemRouter(Router *gin.RouterGroup) {
	systemRouterGroup := Router.Group("system").Use(middleware.JWTAuth())
	{
		systemRouterGroup.GET("status", systemApi.TelegramBotCheck)
	}
}
