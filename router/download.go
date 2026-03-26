package router

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/middleware"
)

type DownloadRouter struct {
}

func (d *DownloadRouter) InitDownloadRouter(Router *gin.RouterGroup) {
	downloadRouter := Router.Group("download").Use(middleware.JWTAuth())
	{
		downloadRouter.GET("chunk/:telegram_file_id", downloadApi.ProxyDownload)
	}
}
