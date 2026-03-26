package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type ChunkRouter struct {
}

func (c *ChunkRouter) InitChunkRouter(Router *gin.RouterGroup) {
	chunkRouter := Router.Group("chunk").Use(middleware.JWTAuth())
	{
		chunkRouter.GET("query/:telegram_file_id", fileChunkApi.QueryDownloadUrl)
		chunkRouter.POST("list", fileChunkApi.ListChunks)
	}
}
