/*
 * @Author: iRorikon
 * @Date: 2025-08-21 19:14:49
 * @FilePath: \CMDB\logger\logger.go
 */
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
)

func SetupGinLog() {
	// gin 日志输出等级
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if config.LogPath != "" {
		LogPath := filepath.Join(config.LogPath, "cloudgram-go.log")
		// 确保日志目录存在
		if err := os.MkdirAll(config.LogPath, os.ModePerm); err != nil {
			log.Fatalf("创建日志目录失败: %v", err)
		}
		fd, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open log file" + fmt.Sprintf("%v", LogPath))
		}
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(os.Stdout, fd)
		gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, fd)
	} else {
		// 如果没有设置日志文件，则使用默认的输出
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}
}

func SysLog(s string) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), s)
}

func SysLogF(f string, v ...any) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), fmt.Sprintf(f, v...))
}

func SysError(s string) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), s)
}

func SysErrorF(f string, v ...any) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), fmt.Sprintf(f, v...))
}

func FatalLog(v ...any) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[FATAL] %v | %v \n", t.Format("2006/01/02 - 15:04:05"), v)
	os.Exit(1)
}

func FatalLogF(f string, v ...any) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[FATAL] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), fmt.Sprintf(f, v...))
	os.Exit(1)
}
