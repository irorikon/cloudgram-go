package controller

import (
	"github.com/irorikon/cloudgram-go/controller/channel"
	"github.com/irorikon/cloudgram-go/controller/file"
	"github.com/irorikon/cloudgram-go/controller/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	FileApiGroup    file.ApiGroup
	ChannelApiGroup channel.ApiGroup
	SystemApiGroup  system.ApiGroup
}
