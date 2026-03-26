package file

import "github.com/irorikon/cloudgram-go/service"

type ApiGroup struct {
	FileApi
	FileChunkApi
	DownloadApi
	UploadChunkApi
}

var (
	fileService      = service.ServiceGroupApp.FileServiceGroup.FileService
	fileChunkService = service.ServiceGroupApp.FileServiceGroup.FileChunkService
	channelService   = service.ServiceGroupApp.ChannelServiceGroup.ChannelService
	telegramService  = service.ServiceGroupApp.TelegramServiceGroup.TelegramBotService
)
