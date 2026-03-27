package initialize

import (
	"io/fs"
	"strings"

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
	Router.Use(gin.Recovery())
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

	// 跨域中间件，注意顺序：在静态文件之前
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求

	// 1. 静态文件服务（必须在API路由之前）
	// 然后执行打包命令 npm run build。在打开下面3行注释
	Router.StaticFile("/favicon.ico", "./dist/favicon.ico")
	Router.Static("/assets", "./dist/assets") // dist里面的静态资源

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

	// 前端路由处理
	// 注意：这必须在所有其他路由之后注册
	// 方法：使用通配符路由处理所有未匹配的前端路由
	Router.NoRoute(func(c *gin.Context) {
		// 如果是 API 请求，返回 404
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}

		// 如果是静态资源请求，也返回 404
		if strings.HasPrefix(c.Request.URL.Path, "/assets") ||
			strings.HasSuffix(c.Request.URL.Path, ".js") ||
			strings.HasSuffix(c.Request.URL.Path, ".css") ||
			strings.HasSuffix(c.Request.URL.Path, ".png") ||
			strings.HasSuffix(c.Request.URL.Path, ".jpg") ||
			strings.HasSuffix(c.Request.URL.Path, ".ico") {
			c.JSON(404, gin.H{"error": "Static resource not found"})
			return
		}

		// 否则返回前端入口文件
		c.File("./dist/index.html")
	})

	config.Routers = Router.Routes()
	logger.SysLog("register router success")
	return Router
}
