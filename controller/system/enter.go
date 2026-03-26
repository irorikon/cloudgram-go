package system

import "github.com/irorikon/cloudgram-go/service"

type ApiGroup struct {
	SystemAuthApi
	SystemApi
}

var (
	telegramService = service.ServiceGroupApp.TelegramServiceGroup.TelegramBotService
)
