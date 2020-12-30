package collectable

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/utils"
)

var imgRegex *regexp.Regexp
var once sync.Once

// GetMarkdownImgRegex 获取编译后的正则表达式
func GetMarkdownImgRegex() *regexp.Regexp {
	once.Do(func() {
		imgRegex, _ = regexp.Compile(`(!\[[^\]\[]+\])\(([^()]+)( ".*")?\)`)
	})
	return imgRegex
}

// MarkdownFile Collectable files which is markdown format files
type MarkdownFile struct {
	*FileAttrs
	buffer []byte
}

// NewMarkdownFile Create a MarkdownFile object which is a collectable file for
// Markdown file
func NewMarkdownFile(parent, uri string) *MarkdownFile {
	reader, fError := utils.NewFileReader(uri)
	defer reader.Close()

	var data []byte
	if fError == nil {
		data, fError = ioutil.ReadAll(reader)
	}

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

	return &MarkdownFile{
		FileAttrs: NewFileAttrs(parent, uri, Markdown, updatedTimePtr, fError),
		buffer:    data,
	}
}

// FindDependencies Returns dependencies which are uri of imgs in the file
func (m *MarkdownFile) FindDependencies() ([]FileOperator, error) {
	if err := m.FileError(); err != nil {
		return nil, err
	}

	dependencies := make([]FileOperator, 0, 3)

	matchs := GetMarkdownImgRegex().FindAllSubmatch(m.buffer, -1)
	for _, match := range matchs {
		path := string(match[2])
		if !filepath.IsAbs(path) {
			// TODO: what if path is a http or https url
			path = filepath.Join(filepath.Dir(m.uri), path)
		}

		dependencies = append(dependencies, NewLeafFile(m.uri, path))
	}

	return dependencies, nil
}

// ReplaceDependencyURIs Replace denpendency's uri in the file
func (m *MarkdownFile) ReplaceDependencyURIs(base, objectKey string, mapper URIMapper) error {
	if err := m.FileError(); err != nil {
		return err
	}

	m.buffer = GetMarkdownImgRegex().ReplaceAllFunc(
		m.buffer,
		func(match []byte) []byte {
			subMatchs := GetMarkdownImgRegex().FindSubmatch(match)
			newURI := mapper(Leaf, subMatchs[2], base, objectKey)
			return bytes.Replace(subMatchs[0], subMatchs[2], newURI, 1)
		},
	)

	return nil
}

// To Write the buffer to file
func (m *MarkdownFile) To(uri string) error {
	if err := m.FileError(); err != nil {
		return err
	}
	if err := utils.CreateDirectory(filepath.Dir(uri)); err != nil {
		return err
	}
	return ioutil.WriteFile(uri, m.buffer, 666)
}

// ToOBS Write the file to bucket
func (m *MarkdownFile) ToOBS(bucket provider.Bucket, key string) error {
	if err := m.FileError(); err != nil {
		return err
	}
	if bucket == nil {
		return errors.New("bucket should not be nil")
	}
	_, err := bucket.PutObjectFromBytes(filepath.ToSlash(key), m.buffer)
	return err
}
