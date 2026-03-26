package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model/response"
)

type SystemApi struct {
}

// health
func (s *SystemApi) Health(c *gin.Context) {
	response.OKWithData(gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   config.Version,
	}, c)
}

// Telegram Bot Check
func (s *SystemApi) TelegramBotCheck(c *gin.Context) {
	// 初始化 Bot
	if err := telegramService.CheckMe(c, config.TeleBot); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OK(c)
}
