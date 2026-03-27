package file

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model/request"
	"github.com/irorikon/cloudgram-go/model/response"
)

type FileApi struct{}

// GetFileDetail 获取文件详情
func (f *FileApi) GetFileDetail(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		response.FailWithMessage("file_id is required", c)
		return
	}

	fileID, ok := f.parseUUID(c, fileIDStr, "file_id")
	if !ok {
		return
	}

	fileRecord, err := fileService.GetFileInfo(fileID)
	if err != nil {
		response.FailWithMessage("file not found", c)
		return
	}
	response.OKWithData(fileRecord, c)
}

// ListFiles 列出文件
func (f *FileApi) ListFiles(c *gin.Context) {
	var req request.FileRequest
	if !f.parseRequest(c, &req) {
		return
	}

	parentID, ok := f.parseParentID(c, req.ParentID)
	if !ok {
		return
	}

	files, err := fileService.GetFileList(parentID)
	if err != nil {
		response.FailWithMessage("failed to get file list", c)
		return
	}
	response.OKWithData(files, c)
}

// ListDirs 列出目录
func (f *FileApi) ListDirs(c *gin.Context) {
	var req request.FileRequest
	if !f.parseRequest(c, &req) {
		return
	}

	parentID, ok := f.parseParentID(c, req.ParentID)
	if !ok {
		return
	}

	folders, err := fileService.GetFoldersByParentId(parentID)
	if err != nil {
		response.FailWithMessage("failed to get folder list", c)
		return
	}
	response.OKWithData(folders, c)
}

// FileExists 检查文件是否存在
func (f *FileApi) FileExists(c *gin.Context) {
	var req request.FileRequest
	if !f.parseRequest(c, &req) {
		return
	}

	if req.FileName == "" {
		response.FailWithMessage("fileName is required", c)
		return
	}

	parentID, ok := f.parseParentID(c, req.ParentID)
	if !ok {
		return
	}

	existingFile, err := fileService.GetFileInfoByName(req.FileName, parentID)
	if err != nil {
		response.FailWithMessage("database error", c)
		return
	}
	response.OKWithData(gin.H{
		"exists": existingFile != nil,
	}, c)
}

// CreateDir 创建目录
func (f *FileApi) CreateDir(c *gin.Context) {
	var req request.FileRequest
	if !f.parseRequest(c, &req) {
		return
	}

	if req.FileName == "" {
		response.FailWithMessage("fileName is required", c)
		return
	}

	parentID, ok := f.parseParentID(c, req.ParentID)
	if !ok {
		return
	}

	fileRecord, err := fileService.AddFileInfo(req.FileName, "", parentID, true, 0)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithData(fileRecord, c)
}

// UpdateFile 更新文件
func (f *FileApi) UpdateFile(c *gin.Context) {
	var req request.FileRequest
	if !f.parseRequest(c, &req) {
		return
	}

	if req.Active == "" {
		response.FailWithMessage("active is required", c)
		return
	}

	if req.Active != "move" && req.Active != "rename" {
		response.FailWithMessage("active must be 'move' or 'rename'", c)
		return
	}

	fileID, ok := f.parseUUID(c, req.ID, "id")
	if !ok {
		return
	}

	var newParentID *uuid.UUID
	if req.Active == "move" {
		// 移动操作需要解析父目录ID
		if req.ParentID != "" {
			parsedID, ok := f.parseParentID(c, req.ParentID)
			if !ok {
				return
			}
			newParentID = parsedID
		} else {
			// 如果ParentID为空，表示移动到根目录，设置为nil
			newParentID = nil
		}
	} else {
		// 重命名操作不需要ParentID
		newParentID = nil
	}

	if req.Active == "rename" && req.NewName == "" {
		response.FailWithMessage("newName is required for rename", c)
		return
	}

	updatedFile, err := fileService.UpdateFileInfo(fileID, req.Active, req.NewName, newParentID)
	if err != nil {
		if err.Error() == "file not found" {
			response.FailWithMessage("file not found", c)
		} else {
			response.FailWithMessage(err.Error(), c)
		}
		return
	}
	response.OKWithData(updatedFile, c)
}

// DeleteFile 删除文件
func (f *FileApi) DeleteFile(c *gin.Context) {
	var req request.FileRequest
	if !f.parseRequest(c, &req) {
		return
	}

	fileID, ok := f.parseUUID(c, req.ID, "id")
	if !ok {
		return
	}

	if req.DeleteWithTG {
		fileIDs, err := fileService.GetFileIDRecursively(fileID)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		for _, fileID := range fileIDs {
			fileChunks, err := fileChunkService.GetFileChunksByFileIDs([]uuid.UUID{fileID})
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			if len(fileChunks) > 0 {
				var msgIDs []int
				var channelID int64
				for _, fileChunk := range fileChunks {
					channelID = fileChunk.ChannelID
					msgIDs = append(msgIDs, fileChunk.TelegramMsgID)
				}
				telegramService.DeleteTelegramFiles(c, config.TeleBot, channelID, msgIDs)
			}
		}
	}
	err := fileService.DeleteFile(fileID, req.Recursive)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithData(true, c)
}

// parseUUID 解析UUID参数
func (f *FileApi) parseUUID(c *gin.Context, uuidStr string, paramName string) (uuid.UUID, bool) {
	if uuidStr == "" {
		response.FailWithMessage(paramName+" is required", c)
		return uuid.Nil, false
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		response.FailWithMessage("invalid "+paramName+" format", c)
		return uuid.Nil, false
	}
	return parsedUUID, true
}

// parseRequest 解析请求体
func (f *FileApi) parseRequest(c *gin.Context, req *request.FileRequest) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.FailWithMessage("invalid request body", c)
		return false
	}
	return true
}

// parseParentID 解析父目录ID
func (f *FileApi) parseParentID(c *gin.Context, parentIDStr string) (*uuid.UUID, bool) {
	if parentIDStr == "" {
		// 空字符串表示根目录，返回nil指针
		return nil, true
	}

	parsedUUID, err := uuid.Parse(parentIDStr)
	if err != nil {
		response.FailWithMessage("invalid parent_id format", c)
		return nil, false
	}

	// 返回指针
	return &parsedUUID, true
}
