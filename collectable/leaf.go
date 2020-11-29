package collectable

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/utils"
)

// LeafFile Collectable file which has no dependencies
type LeafFile struct {
	*FileAttrs
}

// NewLeafFile 创建一个 StandaloneCollector，它可以 collect 没有依赖项的文件
func NewLeafFile(parent, uri string) *LeafFile {
	reader, fError := utils.NewFileReader(uri)
	defer reader.Close()

	var updatedTimePtr *time.Time
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		var fi os.FileInfo
		if fError == nil {
			if fi, fError = os.Stat(uri); fError != nil {
				updatedTime := fi.ModTime()
				updatedTimePtr = &updatedTime
			}
		}

		if absURI, err := filepath.Abs(uri); err == nil {
			uri = absURI
		}
	}

	return &LeafFile{
		FileAttrs: NewFileAttrs(parent, uri, Standalone, updatedTimePtr, fError),
	}
}

// FindDependencies Returns all the dependencies
func (l *LeafFile) FindDependencies() ([]FileOperator, error) {
	if err := l.FileError(); err != nil {
		return nil, err
	}

	dependencies := make([]FileOperator, 0)
	return dependencies, nil
}

// ReplaceDependencyURIs Replaces all the dependencies uri in the file
func (l *LeafFile) ReplaceDependencyURIs(base, objectKey string, mapper URIMapper) error {
	if err := l.FileError(); err != nil {
		return err
	}
	return nil
}

// To Write the file to a new place
func (l *LeafFile) To(uri string) error {
	if err := l.FileError(); err != nil {
		return err
	}
	if err := utils.CreateDirectory(filepath.Dir(uri)); err != nil {
		return err
	}
	return utils.DownloadFile(l.uri, uri)
}

// ToOBS Write the file to bucket
func (l *LeafFile) ToOBS(bucket provider.Bucket, key string) error {
	if err := l.FileError(); err != nil {
		return err
	}
	if bucket == nil {
		return errors.New("bucket should not be nil")
	}
	_, err := bucket.PutObjectFromFile(filepath.ToSlash(key), l.uri)
	return err
}
