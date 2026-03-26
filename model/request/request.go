package request

import (
	"mime/multipart"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/irorikon/cloudgram-go/model"
)

// CustomClaims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	Username string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// FileRequest 文件列表请求
type FileRequest struct {
	Active       string `json:"active,omitempty"`
	ID           string `json:"id,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	NewName      string `json:"new_name,omitempty"`
	ParentID     string `json:"parent_id,omitempty"`
	IsDir        bool   `json:"is_dir,omitempty"`
	Recursive    bool   `json:"recursive,omitempty"`
	DeleteWithTG bool   `json:"delete_with_tg,omitempty"`
}

type ChannelRequest struct {
	ChannelID   int64  `json:"channel_id,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
	Limited     bool   `json:"limited,omitempty"`
	MessageID   int    `json:"message_id,omitempty"`
}

type ChunkRequest struct {
	ChannelID       int64             `json:"channel_id,omitempty"`
	FileID          string            `json:"file_id,omitempty"`
	FileChunkRecord []model.FileChunk `json:"file_chunk_record,omitempty"`
}

type UploadRequest struct {
	ChannelID   int64                 `form:"channel_id" binding:"required"`
	UploadID    string                `form:"upload_id" binding:"required"`
	ChunkIndex  int                   `form:"chunk_index,omitempty"`
	ChunkSize   int                   `form:"chunk_size" binding:"required"`
	TotalChunks int                   `form:"total_chunks" binding:"required"`
	FileName    string                `form:"file_name" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

type MergeChunksRequest struct {
	UploadID    string `json:"upload_id" binding:"required"`
	ChannelID   int64  `json:"channel_id" binding:"required"`
	MessageID   int    `json:"message_id" binding:"required"`
	FileName    string `json:"file_name" binding:"required"`
	ParentID    string `json:"parent_id,omitempty"`
	Size        int64  `json:"size" binding:"required"`
	MimeType    string `json:"mime_type,omitempty"`
	UploadCount int    `json:"upload_count" binding:"required"`
}

type ClearTempChunksRequest struct {
	UploadID  string `json:"upload_id" binding:"required"`
	ChannelID int64  `json:"channel_id" binding:"required"`
}
