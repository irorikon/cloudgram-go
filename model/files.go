package model

import "github.com/google/uuid"

type File struct {
	BaseModel
	ID       uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;comment:文件UUID"`
	Name     string      `json:"name" gorm:"not null;comment:文件名"`
	ParentID uuid.UUID   `json:"parent_id" gorm:"type:uuid;index;comment:父级ID"`
	IsDir    bool        `json:"is_dir" gorm:"not null;default:false;comment:是否是文件夹"`
	Size     int64       `json:"size" gorm:"not null;default:0;comment:文件大小"`
	MimeType string      `json:"mime_type" gorm:"comment:文件类型"`
	Parent   *File       `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE;comment:父级文件"`
	Chunks   []FileChunk `json:"chunks,omitempty" gorm:"foreignKey:FileID;references:ID;comment:文件分片"`
}

func (File) TableName() string {
	return "files"
}
