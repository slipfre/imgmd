package internal

import "os"

// CreateDirectory Create directory recursively
func CreateDirectory(dirPath string) error {
	if !IsFileExist(dirPath) {
		err := os.MkdirAll(dirPath, os.ModePerm)
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
