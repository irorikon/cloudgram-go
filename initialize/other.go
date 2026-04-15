package initialize

import (
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/utils"
	"github.com/mymmrac/telego"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

func OtherInit() {
	// 初始化 JWT
	if _, err := utils.ParseDuration(config.JwtExpiresTime); err != nil {
		panic(err)
	} else {
		config.GlobalJWTConfig.ExpiresTime = config.JwtExpiresTime
	}
	// 初始化其他配置
	dr, err := utils.ParseDuration(config.GlobalJWTConfig.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(config.GlobalJWTConfig.BufferTime)
	if err != nil {
		panic(err)
	}

	config.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)

	// 初始化 Telegram Bot
	config.TeleBot, err = telego.NewBot(config.TelegramBotToken, telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}
}
