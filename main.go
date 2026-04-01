package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/initialize"
	"github.com/irorikon/cloudgram-go/logger"
)

func init() {
	flag.StringVar(&config.Listen, "listen", ":5244", "Listen address")
	flag.StringVar(&config.Listen, "L", ":5244", "Listen address")
	flag.StringVar(&config.AuthUser, "user", "", "Authentication username")
	flag.StringVar(&config.AuthUser, "U", "", "Authentication username")
	flag.StringVar(&config.AuthPassword, "password", "", "Authentication password")
	flag.StringVar(&config.AuthPassword, "P", "", "Authentication password")
	flag.StringVar(&config.TelegramBotToken, "token", "", "Telegram bot token")
	flag.StringVar(&config.TelegramBotToken, "T", "", "Telegram bot token")
	flag.StringVar(&config.DBTYPE, "type", "sqlite", "Database type (sqlite, mysql, postgres)")
	flag.StringVar(&config.DSN, "dsn", "file:./cloudgram.db?cache=shared&mode=rwc", "Database connection string")
	flag.StringVar(&config.LogPath, "log", "", "Log file path")
	flag.StringVar(&config.JwtSecretKey, "jwt-secret", "cloudgram-secret-key", "JWT secret key")
	flag.BoolVar(&config.Debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&config.Debug, "d", false, "Enable debug mode")
	flag.BoolVar(&config.PrintVersion, "version", false, "Print version and exit")
	flag.BoolVar(&config.PrintVersion, "v", false, "Print version and exit")
	flag.Usage = usage
}

var usageStr = `
Usage: cloudgram-go [option]
Optional parameters:
    -L,   --listen                   Listen address (default :8080)
    -U,   --user                     Authentication username
    -P,   --password                 Authentication password
    -T,   --token                    Telegram bot token
          --type                     Database type (sqlite, mysql, postgres, default sqlite)
          --dsn                      Database connection string (default ./cloudgram.db)
          --log                      Log file path
          --jwt-secret               JWT secret key (default cloudgram-secret-key)
    -d,   --debug                    Enable debug mode
    -v,   --version                  Show software version
    -h,   --help                     Show help message
`

func usage() {
	fmt.Println(usageStr)
	os.Exit(0)
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 打印版本信息并退出
	if config.PrintVersion {
		fmt.Printf("cloudgram-go version %s\n", config.Version)
		os.Exit(0)
	}

	initialize.OtherInit()

	// 初始化日志
	logger.SetupGinLog()
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	logger.SysLogF("CloudGram-GO %v started", config.Version)

	// 初始化数据库
	if err := initialize.InitDB(); err != nil {
		logger.FatalLogF("Failed to initialize database: %v", err)
	}
	defer func() {
		err := initialize.CloseDB()
		if err != nil {
			logger.FatalLog(err)
		}
	}()

	if config.DB != nil {
		// 初始化表数据
		initialize.RegisterTables()
	}

	// 启动服务器
	initialize.RunServer()
}
