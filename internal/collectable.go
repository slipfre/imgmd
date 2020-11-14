package internal

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sync"
)

// FileAttributesGetter Attributes of a file
type FileAttributesGetter interface {
	GetFileType() FileType
	GetParent() string
	GetURI() string
	FileError() error
}

// CollectableFileOperator Files can be collected
type CollectableFileOperator interface {
	FileAttributesGetter
	FindDependencies() ([]CollectableFileOperator, error)
	ReplaceDependencyURIs(mapper URIMapper) error
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
	// TODO: 通过链表的形式链接起来
	parent   string
	uri      string
	fileType FileType
	err      error
}

// NewFileAttrs Create and return a FileAttrs object
func NewFileAttrs(parent, uri string, fileType FileType, err error) *FileAttrs {
	return &FileAttrs{
		parent:   parent,
		uri:      uri,
		fileType: fileType,
		err:      err,
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

// FileError Get file error if exist
func (f *FileAttrs) FileError() error {
	return f.err
}

// LeafFile Collectable file which has no dependencies
type LeafFile struct {
	*FileAttrs
}

// NewLeafFile 创建一个 StandaloneCollector，它可以 collect 没有依赖项的文件
func NewLeafFile(parent, uri string) *LeafFile {
	reader, fError := NewFileReader(uri)
	defer reader.Close()

	if absURI, err := filepath.Abs(uri); err == nil {
		uri = absURI
	}

	return &LeafFile{
		FileAttrs: NewFileAttrs(parent, uri, Standalone, fError),
	}
}

// FindDependencies Returns all the dependencies
func (l *LeafFile) FindDependencies() ([]CollectableFileOperator, error) {
	if err := l.FileError(); err != nil {
		return nil, err
	}

	dependencies := make([]CollectableFileOperator, 0)
	return dependencies, nil
}

// ReplaceDependencyURIs Replaces all the dependencies uri in the file
func (l *LeafFile) ReplaceDependencyURIs(mapper URIMapper) error {
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
	return DownloadFile(l.uri, uri)
}

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
	reader, err := NewFileReader(uri)
	defer reader.Close()

	var data []byte
	if err == nil {
		data, err = ioutil.ReadAll(reader)
	}

	if absURI, err := filepath.Abs(uri); err == nil {
		uri = absURI
	}

	return &MarkdownFile{
		FileAttrs: NewFileAttrs(parent, uri, Markdown, err),
		buffer:    data,
	}
}

// FindDependencies Returns dependencies which are uri of imgs in the file
func (m *MarkdownFile) FindDependencies() ([]CollectableFileOperator, error) {
	if err := m.FileError(); err != nil {
		return nil, err
	}

	dependencies := make([]CollectableFileOperator, 0, 3)

	matchs := GetMarkdownImgRegex().FindAllSubmatch(m.buffer, -1)
	for _, match := range matchs {
		path := string(match[2])
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(m.uri), path)
		}

		dependencies = append(dependencies, NewLeafFile(m.uri, path))
	}

	return dependencies, nil
}

// ReplaceDependencyURIs Replace denpendency's uri in the file
func (m *MarkdownFile) ReplaceDependencyURIs(mapper URIMapper) error {
	if err := m.FileError(); err != nil {
		return err
	}

	m.buffer = GetMarkdownImgRegex().ReplaceAllFunc(
		m.buffer,
		func(match []byte) []byte {
			subMatchs := GetMarkdownImgRegex().FindSubmatch(match)
			newURI := mapper(Standalone, subMatchs[2])
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

	return ioutil.WriteFile(uri, m.buffer, 666)
}

type cancelCollectableFile struct {
	*FileAttrs
	collectableFile CollectableFileOperator
	cancelCtx       context.Context
}

// WithCancel Returns a CollectableFileOperator which maybe be cancelled due to the context
func WithCancel(ctx context.Context, cf CollectableFileOperator) (CollectableFileOperator, error) {
	if ctx == nil {
		err := errors.New("ctx should not be nil")
		return nil, err
	}

	return &cancelCollectableFile{
		FileAttrs: NewFileAttrs(
			cf.GetParent(),
			cf.GetURI(),
			cf.GetFileType(),
			cf.FileError(),
		),
		collectableFile: cf,
		cancelCtx:       ctx,
	}, nil
}

// FindDependencies Returns all the dependencies
func (c *cancelCollectableFile) FindDependencies() ([]CollectableFileOperator, error) {
	if yes, err := c.cancelled(); yes {
		return nil, err
	}
	return c.collectableFile.FindDependencies()
}

// ReplaceDependencyURIs Replaces all the dependencies uri in the file
func (c *cancelCollectableFile) ReplaceDependencyURIs(mapper URIMapper) error {
	if yes, err := c.cancelled(); yes {
		return err
	}
	return c.collectableFile.ReplaceDependencyURIs(mapper)
}

// To Write the file to a new place
func (c *cancelCollectableFile) To(uri string) error {
	if yes, err := c.cancelled(); yes {
		return err
	}
	return c.collectableFile.To(uri)
}

func (c *cancelCollectableFile) cancelled() (bool, error) {
	done := c.cancelCtx.Done()
	if done == nil {
		return false, nil
	}

	select {
	case <-done:
		err := c.cancelCtx.Err()
		return true, err
	default:
		return false, nil
	}
}
