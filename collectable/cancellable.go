package collectable

import (
	"context"
	"errors"

	"github.com/slipfre/imgmd/provider"
)

// CancellableFile It's a kind of collectable file which can be cancelled by
// context
type CancellableFile struct {
	*FileAttrs
	collectableFile FileOperator
	cancelCtx       context.Context
}

// WithCancel Returns a CollectableFileOperator which maybe be cancelled due to the context
func WithCancel(ctx context.Context, cf FileOperator) (FileOperator, error) {
	if ctx == nil {
		err := errors.New("ctx should not be nil")
		return nil, err
	}

	updatedTime, _ := cf.GetUpdatedTime()

	return &CancellableFile{
		FileAttrs: NewFileAttrs(
			cf.GetParent(),
			cf.GetURI(),
			cf.GetFileType(),
			updatedTime,
			cf.FileError(),
		),
		collectableFile: cf,
		cancelCtx:       ctx,
	}, nil
}

// FindDependencies Returns all the dependencies
func (c *CancellableFile) FindDependencies() ([]FileOperator, error) {
	if yes, err := c.cancelled(); yes {
		return nil, err
	}
	return c.collectableFile.FindDependencies()
}

// ReplaceDependencyURIs Replaces all the dependencies uri in the file
func (c *CancellableFile) ReplaceDependencyURIs(base, objectKey string, mapper URIMapper) error {
	if yes, err := c.cancelled(); yes {
		return err
	}
	return c.collectableFile.ReplaceDependencyURIs(base, objectKey, mapper)
}

// To Write the file to a new place
func (c *CancellableFile) To(uri string) error {
	if yes, err := c.cancelled(); yes {
		return err
	}
	return c.collectableFile.To(uri)
}

// ToOBS Write the file to bucket
func (c *CancellableFile) ToOBS(bucket provider.Bucket, key string) error {
	if yes, err := c.cancelled(); yes {
		return err
	}
	return c.collectableFile.ToOBS(bucket, key)
}

func (c *CancellableFile) cancelled() (bool, error) {
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
