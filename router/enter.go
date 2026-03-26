package router

import "github.com/irorikon/cloudgram-go/controller"

var RouterGroup struct {
	AuthRouter
	BaseRouter
	ChunkRouter
	ChannelRouter
	DownloadRouter
	FileRouter
	SystemRouter
	UploadRouter
}

var (
	channelApi     = controller.ApiGroupApp.ChannelApiGroup.ChannelApi
	fileApi        = controller.ApiGroupApp.FileApiGroup.FileApi
	fileChunkApi   = controller.ApiGroupApp.FileApiGroup.FileChunkApi
	systemApi      = controller.ApiGroupApp.SystemApiGroup.SystemApi
	systemAuthApi  = controller.ApiGroupApp.SystemApiGroup.SystemAuthApi
	downloadApi    = controller.ApiGroupApp.FileApiGroup.DownloadApi
	uploadChunkApi = controller.ApiGroupApp.FileApiGroup.UploadChunkApi
)
