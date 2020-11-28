package internal

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"

	"github.com/slipfre/imgmd/provider"
)

// Collector Collect collecatable files
type Collector interface {
	Collect(ctx context.Context) <-chan error
}

// CollectorGenerator Generate collectors
type CollectorGenerator func(cf CollectableFileOperator, base, objectKey string, depURIMapper URIMapper, options ...CollectorOption) (Collector, error)

// LocalCollectorGenerator Generate local collectors
func LocalCollectorGenerator(cf CollectableFileOperator, base, objectKey string, depURIMapper URIMapper, options ...CollectorOption) (Collector, error) {
	return NewLocalAsyncCollector(cf, base, objectKey, depURIMapper, options...)
}

// GetOBSCollectorGenerator Return a OBS collector generator which generate obs collectors
func GetOBSCollectorGenerator(bucket provider.Bucket) CollectorGenerator {
	return func(cf CollectableFileOperator, base, objectKey string, depURIMapper URIMapper, options ...CollectorOption) (Collector, error) {
		return NewOBSAsyncCollector(bucket, cf, base, objectKey, depURIMapper, options...)
	}
}

// FreshValidator Validate whether the local file is up to date
type FreshValidator func(cf CollectableFileOperator, base, objectKey string) (bool, error)

// LocalFileFreshValidator Validate whether the local file is up to date
func LocalFileFreshValidator(cf CollectableFileOperator, base, objectKey string) (bool, error) {
	targetPath := filepath.Join(base, objectKey)
	if IsFileExist(targetPath) {
		updatedTime, _ := GetUpdatedTime(targetPath)
		isUpdatedSince, err := cf.IsUpdatedSince(updatedTime)
		if err != nil {
			return false, err
		}
		if !isUpdatedSince {
			// No need for collecting
			return false, nil
		}
	}
	return true, nil
}

// GetOBSFileFreshValidator Get a FreshValidator which validates whether the obs file is up to date
func GetOBSFileFreshValidator(bucket provider.Bucket) (FreshValidator, error) {
	if bucket == nil {
		return nil, errors.New("bucket should not be nil")
	}
	return func(cf CollectableFileOperator, base, objectKey string) (bool, error) {
		exist, err := bucket.IsObjectExist(objectKey)
		if err != nil {
			return false, err
		}
		if exist {
			updatedTime, err := bucket.GetObjectLastModified(objectKey)
			if err != nil {
				return false, err
			}
			isUpdatedSince, err := cf.IsUpdatedSince(updatedTime)
			if err != nil {
				return false, err
			}
			if !isUpdatedSince {
				return false, nil
			}
		}
		return true, nil
	}, nil
}

// Mover Make files to specified place
type Mover func(cf CollectableFileOperator, base, objectKey string) error

// LocalMover Make files to specified local place
func LocalMover(cf CollectableFileOperator, base, objectKey string) error {
	return cf.To(filepath.Join(base, objectKey))
}

// GetOBSMover Return a OBSMover which make files to OBS with specified object key
func GetOBSMover(bucket provider.Bucket) (Mover, error) {
	if bucket == nil {
		return nil, errors.New("bucket should not be nil")
	}
	return func(cf CollectableFileOperator, base, objectKey string) error {
		return cf.ToOBS(bucket, objectKey)
	}, nil
}

type collectorConfigs struct {
	force                 bool
	depCollectorGenerator CollectorGenerator
	depURIMapper          URIMapper
}

func defaultCollectorConfigs() *collectorConfigs {
	configs := &collectorConfigs{
		force:                 false,
		depCollectorGenerator: LocalCollectorGenerator,
	}
	return configs
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

// WithDependencyCollectorGenerator Option config for collectors
func WithDependencyCollectorGenerator(generator CollectorGenerator) CollectorOption {
	return func(configs *collectorConfigs) {
		configs.depCollectorGenerator = generator
	}
}

// AsyncCollector Collector which collect file in another Goroutine
type AsyncCollector struct {
	collectableFile       CollectableFileOperator
	base                  string
	objectKey             string
	targetPath            string
	depCollectorGenerator CollectorGenerator
	depURIMapper          URIMapper
	freshValidator        FreshValidator
	mover                 Mover
	force                 bool
}

// NewAsyncCollector Constructor for NewAsyncCollector
func NewAsyncCollector(cf CollectableFileOperator, base, objectKey string, freshValidator FreshValidator, mover Mover, depURIMapper URIMapper, options ...CollectorOption) (*AsyncCollector, error) {
	if freshValidator == nil {
		return nil, errors.New("'IsNeedCollectValidator' should not be nil")
	}

	if mover == nil {
		return nil, errors.New("'Mover' should not be nil")
	}

	configs := defaultCollectorConfigs()

	for _, option := range options {
		option(configs)
	}

	return &AsyncCollector{
		collectableFile:       cf,
		base:                  base,
		objectKey:             objectKey,
		targetPath:            filepath.Join(base, objectKey),
		freshValidator:        freshValidator,
		depURIMapper:          depURIMapper,
		mover:                 mover,
		depCollectorGenerator: configs.depCollectorGenerator,
		force:                 configs.force,
	}, nil
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

	if !c.force {
		needCollect, err := c.freshValidator(c.collectableFile, c.base, c.objectKey)
		if err != nil {
			complete <- err
			return
		}
		if !needCollect {
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
		depObjDir := GetTargetResourcesDirPath(c.objectKey)

		err = cancelCF.ReplaceDependencyURIs(c.base, c.objectKey, c.depURIMapper)
		if err != nil {
			complete <- err
			return
		}

		subCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		cases := make([]reflect.SelectCase, len(deps))
		for i, dep := range deps {
			depObjKey := filepath.Join(depObjDir, filepath.Base(dep.GetURI()))
			collector, err := c.depCollectorGenerator(dep, c.base, depObjKey, c.depURIMapper, WithForce(c.force))
			if err != nil {
				cancel()
				complete <- err
				return
			}
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

	if err = c.mover(cancelCF, c.base, c.objectKey); err != nil {
		complete <- err
		return
	}

	complete <- nil
}

// NewLocalAsyncCollector Return a AsyncCollector which collect the file to
// local place
func NewLocalAsyncCollector(cf CollectableFileOperator, base, objectKey string, depURIMapper URIMapper, options ...CollectorOption) (*AsyncCollector, error) {
	collector, err := NewAsyncCollector(
		cf, base, objectKey, LocalFileFreshValidator, LocalMover, depURIMapper, options...)
	if err != nil {
		return nil, err
	}
	return collector, nil
}

// NewOBSAsyncCollector Returns a AsyncCollector which collect the file to OBS
func NewOBSAsyncCollector(bucket provider.Bucket, cf CollectableFileOperator, base, objectKey string, depURIMapper URIMapper, options ...CollectorOption) (*AsyncCollector, error) {
	validator, err := GetOBSFileFreshValidator(bucket)
	if err != nil {
		return nil, err
	}
	mover, err := GetOBSMover(bucket)
	if err != nil {
		return nil, err
	}
	collector, err := NewAsyncCollector(cf, base, objectKey, validator, mover, depURIMapper, options...)
	if err != nil {
		return nil, err
	}
	return collector, nil
}
