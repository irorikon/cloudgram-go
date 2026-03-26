package file

import (
	"errors"

	"github.com/google/uuid"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model"
	"gorm.io/gorm"
)

type FileChunkService struct{}

var FileChunkApp = new(FileChunkService)

// CreateFileChunk 创建文件分片记录
func (f *FileChunkService) CreateFileChunk(fileID uuid.UUID, chunkIndex int, chunkSize int, channelID int64, telegramFileID string, telegramMsgID int) (*model.FileChunk, error) {
	chunk := &model.FileChunk{
		FileID:         fileID,
		ChunkIndex:     chunkIndex,
		ChunkSize:      chunkSize,
		ChannelID:      channelID,
		TelegramFileID: telegramFileID,
		TelegramMsgID:  telegramMsgID,
	}

	return chunk, config.DB.Create(chunk).Error
}

// GetFileChunkByID 根据分片ID获取分片信息
func (f *FileChunkService) GetFileChunkByID(id int) (*model.FileChunk, error) {
	var chunk model.FileChunk
	err := config.DB.First(&chunk, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &chunk, err
}

// GetFileChunksByFileID 获取文件分片列表
func (f *FileChunkService) GetFileChunksByFileIDs(fileIDs []uuid.UUID) ([]*model.FileChunk, error) {
	var chunks []*model.FileChunk
	err := config.DB.Where("file_id IN ?", fileIDs).
		Order("file_id, chunk_index ASC").
		Find(&chunks).Error
	return chunks, err
}

// CreateTempChunkRecord 创建临时分片记录
func (f *FileChunkService) CreateTempChunkRecord(uploadID string, chunkIndex, chunkSize int, channelID int64, telegramFileID string, telegramMsgID int) (*model.TempChunk, error) {
	tempChunk := &model.TempChunk{
		UploadID:       uploadID,
		ChunkIndex:     chunkIndex,
		ChunkSize:      chunkSize,
		ChannelID:      channelID,
		TelegramFileID: telegramFileID,
		TelegramMsgID:  telegramMsgID,
	}
	return tempChunk, config.DB.Create(tempChunk).Error
}

// DeleteTempChunkRecordsByUploadID 清理临时分片记录
func (f *FileChunkService) DeleteTempChunkRecordsByUploadID(uploadID string) (int64, error) {
	result := config.DB.Where("upload_id = ?", uploadID).Delete(&model.TempChunk{})
	return result.RowsAffected, result.Error
}

// GetTempChunkRecordsByUploadID 获取临时分片记录
func (f *FileChunkService) GetTempChunkRecordsByUploadID(uploadID string) ([]model.TempChunk, error) {
	var tempChunks []model.TempChunk
	err := config.DB.Where("upload_id = ?", uploadID).Find(&tempChunks).Error
	return tempChunks, err
}

// TransferTempChunksToFile 将临时分片转移到正式文件分片表
func (f *FileChunkService) TransferTempChunksToFile(uploadID string, fileID uuid.UUID) (int, error) {
	var transferredCount int

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 查询临时分片
		var tempChunks []model.TempChunk
		if err := tx.Where("upload_id = ?", uploadID).Order("chunk_index ASC").Find(&tempChunks).Error; err != nil {
			return err
		}

		// 转换为正式分片
		fileChunks := make([]model.FileChunk, len(tempChunks))
		for i, tempChunk := range tempChunks {
			fileChunks[i] = model.FileChunk{
				FileID:         fileID,
				ChunkIndex:     tempChunk.ChunkIndex,
				ChunkSize:      tempChunk.ChunkSize,
				ChannelID:      tempChunk.ChannelID,
				TelegramFileID: tempChunk.TelegramFileID,
				TelegramMsgID:  tempChunk.TelegramMsgID,
				BaseModel:      tempChunk.BaseModel,
			}
		}

		// 批量插入
		if len(fileChunks) > 0 {
			if err := tx.CreateInBatches(fileChunks, 100).Error; err != nil {
				return err
			}
			transferredCount = len(fileChunks)
		}

		// 清理临时记录
		return tx.Where("upload_id = ?", uploadID).Delete(&model.TempChunk{}).Error
	})

	return transferredCount, err
}

// GetTotalChunksByUploadID 通过upload_id计算分片总数
func (f *FileChunkService) GetTotalChunksByUploadID(uploadID string) (int64, error) {
	var count int64
	err := config.DB.Model(&model.TempChunk{}).Where("upload_id = ?", uploadID).Count(&count).Error
	return count, err
}
