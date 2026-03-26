package model

import "github.com/google/uuid"

type FileChunk struct {
	BaseModel
	ID             int       `json:"id" gorm:"primaryKey;comment:分片ID"`
	FileID         uuid.UUID `json:"file_id" gorm:"type:uuid;index;comment:文件ID"`
	ChunkIndex     int       `json:"chunk_index" gorm:"index;comment:分块编号"`
	ChunkSize      int       `json:"chunk_size" gorm:"comment:分块大小"`
	ChannelID      int64     `json:"channel_id" gorm:"comment:频道ID"`
	TelegramFileID string    `json:"telegram_file_id" gorm:"type:text;comment:Telegram文件ID"`
	TelegramMsgID  int       `json:"telegram_msg_id" gorm:"comment:Telegram消息ID"`
	File           *File     `json:"file,omitempty" gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE;comment:文件"`
}

func (FileChunk) TableName() string {
	return "file_chunks"
}

type TempChunk struct {
	BaseModel
	ID             int    `json:"id" gorm:"primaryKey;comment:分片ID"`
	UploadID       string `json:"upload_id" gorm:"comment:上传任务 ID"`
	ChunkIndex     int    `json:"chunk_index" gorm:"comment:分块编号"`
	ChunkSize      int    `json:"chunk_size" gorm:"comment:分块大小"`
	ChannelID      int64  `json:"channel_id" gorm:"comment:频道ID"`
	TelegramFileID string `json:"telegram_file_id" gorm:"comment:Telegram文件ID"`
	TelegramMsgID  int    `json:"telegram_msg_id" gorm:"comment:Telegram消息ID"`
}

func (TempChunk) TableName() string {
	return "temp_chunks"
}
