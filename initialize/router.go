package initialize

import (
	"io/fs"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/logger"
	"github.com/irorikon/cloudgram-go/middleware"
	"github.com/irorikon/cloudgram-go/router"
)

var staticFS fs.FS

func SetStaticFS(fs fs.FS) {
	staticFS = fs
}

// 初始化总路由
func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Logger(), gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}
	authRouter := router.RouterGroup.AuthRouter
	baseRouter := router.RouterGroup.BaseRouter
	chunkRouter := router.RouterGroup.ChunkRouter
	channelRouter := router.RouterGroup.ChannelRouter
	downloadRouter := router.RouterGroup.DownloadRouter
	fileRouter := router.RouterGroup.FileRouter
	systemRouter := router.RouterGroup.SystemRouter
	uploadRouter := router.RouterGroup.UploadRouter

	// 然后执行打包命令 npm run build。在打开下面3行注释
	Router.StaticFile("/favicon.ico", "./dist/favicon.ico")
	Router.Static("/assets", "./dist/assets")   // dist里面的静态资源
	Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	// 跨域，如需跨域可以打开下面的注释
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求

	// 方便统一添加路由组前缀 多服务器上线使用
	PublicGroup := Router.Group("api")
	PrivateGroup := Router.Group("api")

	{
		systemRouter.InitSystemRouter(PublicGroup)
		baseRouter.InitBaseRouter(PublicGroup)
	}
	{
		authRouter.InitAuthRouter(PrivateGroup)
		channelRouter.InitChannelRouter(PrivateGroup)
		chunkRouter.InitChunkRouter(PrivateGroup)
		downloadRouter.InitDownloadRouter(PrivateGroup)
		fileRouter.InitFileRouter(PrivateGroup)
		uploadRouter.InitUploadRouter(PrivateGroup)
	}

	config.Routers = Router.Routes()
	logger.SysLog("register router success")
	return Router
}
