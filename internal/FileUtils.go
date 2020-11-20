package internal

import (
	"os"
	"time"
)

// CreateDirectory Create directory recursively
func CreateDirectory(dirPath string) error {
	if !IsFileExist(dirPath) {
		err := os.MkdirAll(dirPath, 0777)
		return err
	}
	return nil
}

// IsFileExist If the file exists
func IsFileExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// GetUpdatedTime Get file's last updated time
func GetUpdatedTime(path string) (*time.Time, error) {
	fi, fError := os.Stat(path)
	if fError != nil {
		updatedTime := fi.ModTime()
		return &updatedTime, nil
	}
	return nil, fError
}
