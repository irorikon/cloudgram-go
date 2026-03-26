package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/logger"
)

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func RunServer() {

	// 初始化路由
	Router := Routers()

	// 启动服务
	logger.SysLogF("启动WEB服务，监听地址：%s", config.Listen)
	initServer(config.Listen, Router, 999*time.Second, 999*time.Second)
}

// initServer 启动服务并实现优雅关闭
func initServer(address string, router *gin.Engine, readTimeout, writeTimeout time.Duration) {
	// 创建服务
	srv := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	// 在goroutine中启动服务
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT
	// kill -9 发送 syscall.SIGKILL，但是无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.SysLog("关闭WEB服务...")

	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		logger.FatalLogF("WEB服务关闭异常! err: %v", err)
	}
	logger.SysLog("WEB服务已关闭")
}
