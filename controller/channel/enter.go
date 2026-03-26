package channel

import "github.com/irorikon/cloudgram-go/service"

type ApiGroup struct {
	ChannelApi
}

var (
	channelService  = service.ServiceGroupApp.ChannelServiceGroup.ChannelService
	telegramService = service.ServiceGroupApp.TelegramServiceGroup.TelegramBotService
)
