package file

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model/request"
	"github.com/irorikon/cloudgram-go/model/response"
)

type FileChunkApi struct{}

// 查询文件分片下载地址
func (f *FileChunkApi) QueryDownloadUrl(c *gin.Context) {
	telegramFileID := c.Param("telegram_file_id")
	if telegramFileID == "" {
		response.FailWithMessage("file_id is required", c)
		return
	}
	url, err := telegramService.GetTelegramFileUrl(c, config.TeleBot, telegramFileID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithData(gin.H{"url": url}, c)
}

// 查询所有文件分片
func (f *FileChunkApi) ListChunks(c *gin.Context) {
	var req request.ChunkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	if req.FileID == "" || strings.TrimSpace(req.FileID) == "" {
		response.FailWithMessage("file_ids is required", c)
		return
	}

	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	chunks, err := fileChunkService.GetFileChunksByFileIDs([]uuid.UUID{fileID})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithData(chunks, c)
}
