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

type collectorConfigs struct {
	force bool
}

func defaultCollectorConfigs() *collectorConfigs {
	return &collectorConfigs{
		force: false,
	}
}

// CollectorOption Options for collectors
type CollectorOption func(configs *collectorConfigs)

// WithForce Option config for collectors. If force is true, the file will be
// collected even though the target file exists and up-to-date
func WithForce(force bool) CollectorOption {
	return func(configs *collectorConfigs) {
		configs.force = force
	}
}

// AsyncCollector Collector which collect file in another Goroutine
type AsyncCollector struct {
	collectableFile CollectableFileOperator
	targetPath      string
	force           bool
}

// NewAsyncCollector Constructor for NewAsyncCollector
func NewAsyncCollector(cf CollectableFileOperator, targetPath string, options ...CollectorOption) *AsyncCollector {
	configs := defaultCollectorConfigs()

	for _, option := range options {
		option(configs)
	}

	return &AsyncCollector{
		collectableFile: cf,
		targetPath:      targetPath,
		force:           configs.force,
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

	if IsFileExist(c.targetPath) {
		updatedTime, _ := GetUpdatedTime(c.targetPath)
		isUpdatedSince, err := c.collectableFile.IsUpdatedSince(updatedTime)
		if err != nil {
			complete <- err
			return
		}
		if !isUpdatedSince {
			// No need for collecting
			complete <- nil
			return
		}
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
			forceOption := WithForce(false)
			if c.force {
				forceOption = WithForce(true)
			}
			depTargetPath := filepath.Join(targetResourcesDirPath, filepath.Base(dep.GetURI()))

			collector := NewAsyncCollector(dep, depTargetPath, forceOption)
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
