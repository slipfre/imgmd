package collectable

import (
	"time"

	"github.com/slipfre/imgmd/provider"
)

// AttributesGetter Attributes of a file
type AttributesGetter interface {
	GetFileType() FileType
	GetParent() string
	GetURI() string
	GetUpdatedTime() (*time.Time, error)
	IsUpdatedSince(time *time.Time) (bool, error)
	FileError() error
}

// FileOperator Files can be collected
type FileOperator interface {
	AttributesGetter
	FindDependencies() ([]FileOperator, error)
	ReplaceDependencyURIs(base, objectKey string, mapper URIMapper) error
	To(uri string) error
	ToOBS(bucket provider.Bucket, key string) error
}

// URIMapper Map the uri
type URIMapper func(fileType FileType, originURI []byte, base, objectKey string) []byte

// FileAttrs Attributes of the collectable file
type FileAttrs struct {
	parent      string
	uri         string
	fileType    FileType
	updatedTime *time.Time
	err         error
}

// NewFileAttrs Create and return a FileAttrs object
func NewFileAttrs(parent, uri string, fileType FileType, updatedTime *time.Time, err error) *FileAttrs {
	return &FileAttrs{
		parent:      parent,
		uri:         uri,
		fileType:    fileType,
		updatedTime: updatedTime,
		err:         err,
	}
}

// FileType FileType of collectable files
type FileType string

const (
	// Standalone Stand for collectable files with no dependencies
	Standalone FileType = "standalone"
	// Markdown Stand for Markdown files
	Markdown FileType = "markdown"
	// None Stand for a file which is not exist
	None FileType = "none"
)

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

// GetUpdatedTime Get lastest update time
func (f *FileAttrs) GetUpdatedTime() (*time.Time, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.updatedTime, f.err
}

// IsUpdatedSince For local file, Returns true if the file has updated.
func (f *FileAttrs) IsUpdatedSince(time *time.Time) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	if f.updatedTime == nil {
		return true, nil
	}
	return time.After(*f.updatedTime), nil
}
