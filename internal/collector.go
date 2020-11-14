package internal

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"
)

// Collector Collect collecatable files
type Collector interface {
	Collect(ctx context.Context) <-chan error
}

// AsyncCollector Collector which collect file in another Goroutine
type AsyncCollector struct {
	collectableFile CollectableFileOperator
	targetPath      string
}

// NewAsyncCollector Constructor for NewAsyncCollector
func NewAsyncCollector(cf CollectableFileOperator, targetPath string) *AsyncCollector {
	return &AsyncCollector{
		collectableFile: cf,
		targetPath:      targetPath,
	}
}

// Collect Collect the collectableFile
func (c *AsyncCollector) Collect(ctx context.Context) <-chan error {
	complete := make(chan error)
	go c.collectFileAsync(ctx, complete)
	return complete
}

func (c *AsyncCollector) collectFileAsync(ctx context.Context, complete chan<- error) {
	cancelCF, err := WithCancel(ctx, c.collectableFile)
	if err != nil {
		complete <- err
		return
	}

	deps, err := cancelCF.FindDependencies()
	if err != nil {
		complete <- err
		return
	}

	if deps != nil && len(deps) > 0 {
		targetResourcesDirPath := c.getTargetResourcesDirPath()
		if err = CreateDirectory(targetResourcesDirPath); err != nil {
			complete <- err
			return
		}

		err = cancelCF.ReplaceDependencyURIs(func(fileType FileType, uri []byte) []byte {
			dirName := filepath.Base(targetResourcesDirPath)
			fileName := filepath.Base(string(uri))
			newReferencePath := filepath.Join(dirName, fileName)
			return []byte(newReferencePath)
		})
		if err != nil {
			complete <- err
			return
		}

		subCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		cases := make([]reflect.SelectCase, len(deps))
		for i, dep := range deps {
			collector := NewAsyncCollector(
				dep,
				filepath.Join(targetResourcesDirPath, filepath.Base(dep.GetURI())),
			)
			subComplete := collector.Collect(subCtx)
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(subComplete),
			}
		}

		for remaining := len(deps); remaining > 0; remaining-- {
			chosen, value, _ := reflect.Select(cases)
			if value.Interface() != nil {
				err = value.Interface().(error)
			}
			if err != nil {
				cancel()
				complete <- err.(error)
				return
			}
			cases[chosen].Chan = reflect.ValueOf(nil)
			continue
		}
	}

	if err = cancelCF.To(c.targetPath); err != nil {
		complete <- err
		return
	}

	complete <- nil
}

func (c *AsyncCollector) getTargetResourcesDirPath() string {
	targetURI := c.targetPath
	directory := filepath.Dir(targetURI)
	filenameWithSuffix := filepath.Base(targetURI)
	suffix := filepath.Ext(filenameWithSuffix)
	filename := strings.TrimSuffix(filenameWithSuffix, suffix)
	return filepath.Join(directory, filename) + "_medias"
}
