package file

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/logger"
	"github.com/irorikon/cloudgram-go/model/request"
	"github.com/irorikon/cloudgram-go/model/response"
)

type UploadChunkApi struct{}

func (u *UploadChunkApi) UploadChunk(c *gin.Context) {
	var req request.UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage("Request body must be valid From", c)
		return
	}

	if req.ChannelID == 0 || req.ChunkIndex < 0 || req.ChunkSize < 0 || req.File == nil || req.UploadID == "" {
		response.FailWithMessage("Invalid request body", c)
		return
	}
	file, err := req.File.Open()
	if err != nil {
		response.FailWithMessage("File open failed", c)
		return
	}
	defer file.Close()

	res, err := telegramService.UploadFileToTelegram(c, config.TeleBot, req.ChannelID, file, req.FileName)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(strings.ToLower(errorMsg), "network") ||
			strings.Contains(strings.ToLower(errorMsg), "timeout") ||
			strings.Contains(strings.ToLower(errorMsg), "connection") {
			response.FailWithMessage("Network connection failed, please try again.", c)
			return
		}
		response.FailWithMessage("Upload failed", c)
		return
	}
	if res.Document.FileID == "" {
		response.FailWithMessage("Upload failed", c)
		return
	}
	// 创建临时分片记录
	fileChunkRecord, err := fileChunkService.CreateTempChunkRecord(
		req.UploadID,
		req.ChunkIndex,
		req.ChunkSize,
		req.ChannelID,
		res.Document.FileID,
		res.MessageID,
	)
	if err != nil {
		response.FailWithMessage("CreateTempChunkRecord failed", c)
		return
	}
	_, err = channelService.UpdateTGChannelMessageId(req.ChannelID, res.MessageID, false)
	if err != nil {
		response.FailWithMessage("Update channel message id failed", c)
		return
	}
	response.OKWithData(fileChunkRecord, c)
}

func (u *UploadChunkApi) MergeChunks(c *gin.Context) {
	var req request.MergeChunksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SysError(err.Error())
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}

	var parentID *uuid.UUID
	if req.ParentID != "" {
		parsedID, err := uuid.Parse(req.ParentID)
		if err != nil {
			response.FailWithMessage("Invalid parentId format", c)
			return
		}
		parentID = &parsedID
	} else {
		// 父目录ID为空字符串时，表示根目录，设置为nil
		parentID = nil
	}

	totalChunks, err := fileChunkService.GetTotalChunksByUploadID(req.UploadID)
	if err != nil {
		response.FailWithMessage("Failed to get total chunks", c)
		return
	}

	if totalChunks != int64(req.UploadCount) {
		response.FailWithMessage("Total chunks mismatch", c)
		return
	}

	file, err := fileService.AddFileInfo(
		req.FileName,
		req.MimeType,
		parentID, // 这里已经是指针类型
		false,
		req.Size,
	)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if _, err = channelService.UpdateTGChannelMessageId(req.ChannelID, req.MessageID, true); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if _, err := fileChunkService.TransferTempChunksToFile(req.UploadID, file.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithData(file, c)
}

func (u *UploadChunkApi) CleanChunks(c *gin.Context) {
	var req request.ClearTempChunksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	tempChunks, err := fileChunkService.GetTempChunkRecordsByUploadID(req.UploadID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var msgIDs []int
	for _, tempChunk := range tempChunks {
		msgIDs = append(msgIDs, tempChunk.TelegramMsgID)
	}
	if err := telegramService.DeleteTelegramFiles(c, config.TeleBot, req.ChannelID, msgIDs); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := fileChunkService.DeleteTempChunkRecordsByUploadID(req.UploadID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithData(true, c)
}
