package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type ChannelRouter struct{}

func (c *ChannelRouter) InitChannelRouter(Router *gin.RouterGroup) {
	channelRouter := Router.Group("channel").Use(middleware.JWTAuth())
	{
		channelRouter.GET("list", channelApi.ListChannel)
		channelRouter.GET("get", channelApi.GetChannel)
		channelRouter.POST("create", channelApi.CreateChannel)
		channelRouter.POST("update", channelApi.UpdateChannelName)
		channelRouter.POST("delete", channelApi.DeleteTGChannel)
		channelRouter.POST("check", channelApi.CheckTelegramChannel)
	}
}
