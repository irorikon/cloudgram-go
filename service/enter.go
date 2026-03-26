package service

import (
	"github.com/irorikon/cloudgram-go/service/channel"
	"github.com/irorikon/cloudgram-go/service/file"
	"github.com/irorikon/cloudgram-go/service/telegram"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	FileServiceGroup     file.ServiceGroup
	TelegramServiceGroup telegram.ServiceGroup
	ChannelServiceGroup  channel.ServiceGroup
}
