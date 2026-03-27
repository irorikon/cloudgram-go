package file

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model"
	"gorm.io/gorm"
)

type FileService struct{}

var FileApp = new(FileService)

// GetFileList 根据 ParentID 获取文件列表
func (f *FileService) GetFileList(parentID *uuid.UUID) ([]*model.File, error) {
	var files []*model.File
	query := config.DB

	// 处理父ID为nil的情况（根目录）
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	err := query.Order("is_dir DESC, name ASC").Find(&files).Error
	return files, err
}

// GetFoldersByParentId 根据父目录ID获取文件夹列表
func (f *FileService) GetFoldersByParentId(parentID *uuid.UUID) ([]*model.File, error) {
	var folders []*model.File
	query := config.DB.Where("is_dir = ?", true)

	// 处理父ID为nil的情况（根目录）
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	err := query.Order("name ASC").Find(&folders).Error
	return folders, err
}

// GetFileInfo 根据ID获取文件信息
func (f *FileService) GetFileInfo(fileID uuid.UUID) (*model.File, error) {
	var file model.File
	err := config.DB.First(&file, "id = ?", fileID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("file not found")
	}
	return &file, err
}

// GetFileIDRecursively 递归获取文件ID
// 如果传入的是文件，返回包含该文件ID的切片
// 如果传入的是目录，返回该目录下所有文件的ID
func (f *FileService) GetFileIDRecursively(id uuid.UUID) ([]uuid.UUID, error) {
	query := `
		WITH RECURSIVE descendants(id, is_dir) AS (
			SELECT id, is_dir FROM files WHERE id = ?
			UNION ALL
			SELECT f.id, f.is_dir FROM files f
			INNER JOIN descendants d ON f.parent_id = d.id
		)
		SELECT id FROM descendants WHERE is_dir = false
	`

	var fileIDs []uuid.UUID
	err := config.DB.Raw(query, id).Scan(&fileIDs).Error
	return fileIDs, err
}

// GetFileDetailedRecursively 递归获取文件详细信息
func (f *FileService) GetFileDetailedRecursively(id uuid.UUID) ([]*model.File, error) {
	query := `
		WITH RECURSIVE descendants(id, is_dir) AS (
			SELECT id, is_dir FROM files WHERE id = ?
			UNION ALL
			SELECT f.id, f.is_dir FROM files f
			INNER JOIN descendants d ON f.parent_id = d.id
		)
		SELECT f.* FROM files f
		INNER JOIN descendants d ON f.id = d.id
		WHERE d.is_dir = false
	`

	var files []*model.File
	err := config.DB.Raw(query, id).Scan(&files).Error
	return files, err
}

// GetFileInfoByName 根据文件名和父目录查询文件
func (f *FileService) GetFileInfoByName(name string, parentID *uuid.UUID) (*model.File, error) {
	var file model.File
	query := config.DB.Where("name = ?", name)

	// 处理父ID为nil的情况
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	err := query.First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

// AddFileInfo 添加文件
func (f *FileService) AddFileInfo(name, mimeType string, parentID *uuid.UUID, isDir bool, size int64) (*model.File, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("invalid file name")
	}

	// 检查文件是否已存在
	existing, err := f.GetFileInfoByName(name, parentID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("file already exists")
	}

	// 创建新文件
	file := &model.File{
		ID:       uuid.New(),
		Name:     name,
		ParentID: parentID,
		IsDir:    isDir,
		Size:     size,
		MimeType: mimeType,
	}
	err = config.DB.Create(file).Error
	return file, err
}

// UpdateFileInfo 更新文件
func (f *FileService) UpdateFileInfo(fileID uuid.UUID, action, newName string, newParentID *uuid.UUID) (*model.File, error) {
	switch action {
	case "rename":
		return f.renameFile(fileID, newName)
	case "move":
		return f.moveFile(fileID, newParentID)
	}
	return nil, errors.New("invalid action")
}

func (f *FileService) renameFile(fileID uuid.UUID, newName string) (*model.File, error) {
	newName = strings.TrimSpace(newName)
	if newName == "" || fileID == uuid.Nil {
		return nil, errors.New("invalid params")
	}

	file, err := f.GetFileInfo(fileID)
	if err != nil || file == nil {
		return nil, fmt.Errorf("file not found: %s", fileID.String())
	}

	// 检查目标位置是否已存在同名文件
	if exists, _ := f.GetFileInfoByName(newName, file.ParentID); exists != nil {
		return nil, errors.New("file name already exists")
	}

	err = config.DB.Model(&model.File{}).Where("id = ?", fileID).
		Update("name", newName).Error
	return f.GetFileInfo(fileID)
}

// moveFile 移动文件
func (f *FileService) moveFile(fileID uuid.UUID, newParentID *uuid.UUID) (*model.File, error) {
	if fileID == uuid.Nil {
		return nil, errors.New("invalid file id")
	}

	file, err := f.GetFileInfo(fileID)
	if err != nil || file == nil {
		return nil, fmt.Errorf("file not found: %s", fileID.String())
	}

	// 检查不能移动到自身
	if newParentID != nil && *newParentID == fileID {
		return nil, errors.New("cannot move to self")
	}

	// 如果是目录，检查循环依赖（目标目录不是根目录时才检查）
	if file.IsDir && newParentID != nil {
		// 检查 newParentID 是否是 fileID 的子孙
		isDescendant, err := f.isDescendant(fileID, *newParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check directory hierarchy: %v", err)
		}
		if isDescendant {
			return nil, errors.New("cannot move a directory into its own subdirectory")
		}
	}

	// 检查目标目录是否存在同名文件
	if exists, _ := f.GetFileInfoByName(file.Name, newParentID); exists != nil && exists.ID != fileID {
		return nil, errors.New("file name already exists in the target directory")
	}

	// 更新父目录
	err = config.DB.Model(&model.File{}).Where("id = ?", fileID).
		Update("parent_id", newParentID).Error
	if err != nil {
		return nil, err
	}

	return f.GetFileInfo(fileID)
}

// isDescendant 检查 descendant 是否是 ancestor 的子孙
func (f *FileService) isDescendant(ancestor, descendant uuid.UUID) (bool, error) {
	if ancestor == uuid.Nil || descendant == uuid.Nil {
		return false, nil
	}

	query := `
		WITH RECURSIVE children(id) AS (
			SELECT id FROM files WHERE id = ?
			UNION ALL
			SELECT f.id FROM files f
			INNER JOIN children c ON f.parent_id = c.id
		)
		SELECT 1 FROM children WHERE id = ?
		LIMIT 1
	`

	var exists bool
	err := config.DB.Raw(query, ancestor, descendant).Scan(&exists).Error
	return exists, err
}

// DeleteFile 删除文件及其分片信息
func (f *FileService) DeleteFile(fileID uuid.UUID, recursive bool) error {
	if fileID == uuid.Nil {
		return errors.New("invalid file id")
	}

	var file model.File
	if err := config.DB.Select("is_dir").Where("id = ?", fileID).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file not found")
		}
		return err
	}

	if file.IsDir && !recursive {
		return errors.New("cannot delete a directory without recursive flag")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 获取要删除的文件ID
		var ids []uuid.UUID
		if file.IsDir {
			query := `WITH RECURSIVE descendants(id) AS (
				SELECT id FROM files WHERE id = ?
				UNION ALL
				SELECT f.id FROM files f
				INNER JOIN descendants d ON f.parent_id = d.id
			) SELECT id FROM descendants`
			if err := tx.Raw(query, fileID).Scan(&ids).Error; err != nil {
				return err
			}
		} else {
			ids = []uuid.UUID{fileID}
		}

		// 如果没有文件要删除，直接返回
		if len(ids) == 0 {
			return nil
		}

		// 删除分片
		if err := tx.Delete(&model.FileChunk{}, "file_id IN ?", ids).Error; err != nil {
			return err
		}

		// 删除文件
		return tx.Delete(&model.File{}, "id IN ?", ids).Error
	})
}
