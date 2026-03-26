package file

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/logger"
	"github.com/irorikon/cloudgram-go/model/response"
)

type DownloadApi struct{}

// 代理下载文件分片
func (d *DownloadApi) ProxyDownload(c *gin.Context) {
	telegramFileID := c.Param("telegram_file_id")
	if telegramFileID == "" {
		response.FailWithMessage("file_id is required", c)
		return
	}
	result, err := telegramService.DownloadTelegramFile(c, config.TeleBot, telegramFileID)
	if err != nil {
		logger.SysError(err.Error())
		response.FailWithMessage("Proxy download failed", c)
		return
	}
	logger.SysLogF("Size: %d", len(result))

	// 设置响应头并返回文件数据
	c.Data(200, "application/octet-stream", result)
}
