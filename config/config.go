package config

import (
	"golang.org/x/sync/singleflight"

	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"gorm.io/gorm"
)

var (
	PrintVersion     bool
	Debug            bool
	Listen           string
	AuthUser         string
	AuthPassword     string
	TelegramBotToken string

	DBTYPE       string
	DSN          string
	LogPath      string
	JwtSecretKey string

	Version string = "1.0.7"
)

// 全局参数配置
var (
	DB                 *gorm.DB
	TeleBot            *telego.Bot
	Routers            gin.RoutesInfo
	BlackCache         local_cache.Cache
	ConcurrencyControl = &singleflight.Group{}
)

// 自定义跨域请求配置 写死
var GlobalCorsConfig []WhiteListConfig = []WhiteListConfig{
	{
		AllowOrigin:      "*",
		AllowCredentials: true,
		AllowHeaders:     "*",
		AllowMethods:     "*",
		ExposeHeaders:    "*",
	},
}

var GlobalJWTConfig = JWTConfig{
	BufferTime:  "1h",
	ExpiresTime: "24h",
	Issuer:      "CloudGram-Go",
}

type WhiteListConfig struct {
	AllowOrigin      string `mapstructure:"allow_origin,omitempty" yaml:"allow_origin,omitempty" json:"allow_origin,omitempty"`
	AllowCredentials bool   `mapstructure:"allow_credentials,omitempty" yaml:"allow_credentials,omitempty" json:"allow_credentials,omitempty"`
	AllowHeaders     string `mapstructure:"allow_headers,omitempty" yaml:"allow_headers,omitempty" json:"allow_headers,omitempty"`
	AllowMethods     string `mapstructure:"allow_methods,omitempty" yaml:"allow_methods,omitempty" json:"allow_methods,omitempty"`
	ExposeHeaders    string `mapstructure:"expose_headers,omitempty" yaml:"expose_headers,omitempty" json:"expose_headers"`
}

type JWTConfig struct {
	Issuer      string
	ExpiresTime string
	BufferTime  string
}
