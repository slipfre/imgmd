package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// FileAttributesGetter Attributes of a file
type FileAttributesGetter interface {
	GetFileType() FileType
	GetParent() string
	GetURI() string
}

// CollectableFileOperator Files can be collected
type CollectableFileOperator interface {
	FileAttributesGetter
	FindDependencies() ([]CollectableFileOperator, error)
	ReplaceDependencyURIs(mapper URIMapper)
	To(uri string) error
}

// URIMapper Map the uri
type URIMapper func(fileType FileType, uri []byte) []byte

// FileType Type of collectable files
type FileType string

const (
	// Standalone Stand for collectable files with no dependencies
	Standalone FileType = "standalone"
	// Markdown Stand for Markdown files
	Markdown FileType = "markdown"
	// None Stand for a file which is not exist
	None FileType = "none"
)

// FileAttrs Attributes of the collectable file
type FileAttrs struct {
	parent   string
	uri      string
	fileType FileType
}

// NewFileAttrs Create and return a FileAttrs object
func NewFileAttrs(parent, uri string, fileType FileType) *FileAttrs {
	return &FileAttrs{
		parent:   parent,
		uri:      uri,
		fileType: fileType,
	}
}

// GetParent Get parent of dependency
func (f *FileAttrs) GetParent() string {
	return f.parent
}

// GetURI Get uri of dependency
func (f *FileAttrs) GetURI() string {
	return f.uri
}

// GetFileType Get file type of dependency
func (f *FileAttrs) GetFileType() FileType {
	return f.fileType
}

// LeafFile Collectable file which has no dependencies
type LeafFile struct {
	*FileAttrs
}

// NewLeafFile 创建一个 StandaloneCollector，它可以 collect 没有依赖项的文件
func NewLeafFile(parent, uri string) (*LeafFile, error) {
	var err error

	if _, err := os.Stat(uri); err != nil {
		return nil, err
	}

	if absURI, err := filepath.Abs(uri); err == nil {
		return &LeafFile{
			FileAttrs: NewFileAttrs(parent, absURI, Standalone),
		}, nil
	}

	return nil, err
}

// FindDependencies Returns all the dependencies
func (l *LeafFile) FindDependencies() ([]CollectableFileOperator, error) {
	dependencies := make([]CollectableFileOperator, 0)
	return dependencies, nil
}

// ReplaceDependencyURIs Replaces all the dependencies uri in the file
func (l *LeafFile) ReplaceDependencyURIs(mapper URIMapper) {
	return
}

// To Write the file to a new place
func (l *LeafFile) To(uri string) error {
	return DownloadFile(l.uri, uri)
}

var imgRegex *regexp.Regexp
var once sync.Once

// GetMarkdownImgRegex 获取编译后的正则表达式
func GetMarkdownImgRegex() *regexp.Regexp {
	var err error
	once.Do(func() {
		imgRegex, err = regexp.Compile(`(!\[[^\]\[]+\])\(([^()]+)( ".*")?\)`)
		fmt.Print(err)
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
func NewMarkdownFile(parent, uri string) (*MarkdownFile, error) {
	data, err := ioutil.ReadFile(uri)
	if err != nil {
		return nil, err
	}

	if absURI, err := filepath.Abs(uri); err == nil {
		return &MarkdownFile{
			FileAttrs: NewFileAttrs(parent, absURI, Markdown),
			buffer:    data,
		}, nil
	}

	return nil, err
}

// FindDependencies Returns dependencies which are uri of imgs in the file
func (m *MarkdownFile) FindDependencies() ([]CollectableFileOperator, error) {
	dependencies := make([]CollectableFileOperator, 0, 3)

	matchs := GetMarkdownImgRegex().FindAllSubmatch(m.buffer, -1)
	for _, match := range matchs {
		path := string(match[2])
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(m.uri), path)
		}

		if dependency, err := NewLeafFile(m.uri, path); err == nil {
			dependencies = append(dependencies, dependency)
		} else {
			return nil, err
		}
	}
	return dependencies, nil
}

// ReplaceDependencyURIs Replace denpendency's uri in the file
func (m *MarkdownFile) ReplaceDependencyURIs(mapper URIMapper) {
	m.buffer = GetMarkdownImgRegex().ReplaceAllFunc(
		m.buffer,
		func(match []byte) []byte {
			subMatchs := GetMarkdownImgRegex().FindSubmatch(match)
			newURI := mapper(Standalone, subMatchs[2])
			return bytes.Replace(subMatchs[0], subMatchs[2], newURI, 1)
		},
	)
}

// To Write the buffer to file
func (m *MarkdownFile) To(uri string) error {
	return ioutil.WriteFile(uri, m.buffer, 666)
}
