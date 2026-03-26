package utils

// utils/files.go
// 文件操作工具函数
// 提供文件和目录的创建、删除、检查等功能

import "os"

// 检测文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 创建文件夹
func CreateDirIfNotExists(dir string) error {
	if !FileExists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// 创建文件，如果文件已存在则不做任何操作
func CreateFileIfNotExists(filePath string) (*os.File, error) {
	if !FileExists(filePath) {
		return os.Create(filePath)
	}
	return nil, nil
}

// 删除文件
func DeleteFile(filePath string) error {
	if FileExists(filePath) {
		return os.Remove(filePath)
	}
	return nil
}
