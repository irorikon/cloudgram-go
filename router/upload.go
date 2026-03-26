package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type UploadRouter struct {
}

func (u *UploadRouter) InitUploadRouter(Router *gin.RouterGroup) {
	uploadRouter := Router.Group("upload").Use(middleware.JWTAuth())
	{
		uploadRouter.POST("chunk", uploadChunkApi.UploadChunk)
		uploadRouter.POST("merge", uploadChunkApi.MergeChunks)
		uploadRouter.POST("cleanup", uploadChunkApi.CleanChunks)
	}
}
