package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type FileRouter struct {
}

func (f *FileRouter) InitFileRouter(Router *gin.RouterGroup) {
	fileRouter := Router.Group("file").Use(middleware.JWTAuth())
	{
		fileRouter.GET("detail/:file_id", fileApi.GetFileDetail)
		fileRouter.POST("list", fileApi.ListFiles)
		fileRouter.POST("dir", fileApi.ListDirs)
		fileRouter.POST("exists", fileApi.FileExists)
		fileRouter.POST("create", fileApi.CreateDir)
		fileRouter.POST("update", fileApi.UpdateFile)
		fileRouter.POST("delete", fileApi.DeleteFile)
	}
}
