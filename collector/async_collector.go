package collector

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"

	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/utils"
)

// AsyncCollector Collector which collect file in another Goroutine
type AsyncCollector struct {
	collectableFile       collectable.FileOperator
	base                  string
	objectKey             string
	targetPath            string
	force                 bool
	depCollectorGenerator Generator
	depURIMapper          collectable.URIMapper
	freshValidator        FreshValidator
	mover                 Mover
}

func defaultCollectorConfigs() *Configs {
	configs := &Configs{
		Force:                 false,
		DepCollectorGenerator: LocalCollectorGenerator,
	}
	return configs
}

// NewAsyncCollector Constructor for NewAsyncCollector
func NewAsyncCollector(cf collectable.FileOperator, base, objectKey string, freshValidator FreshValidator, mover Mover, depURIMapper collectable.URIMapper, options ...Option) (*AsyncCollector, error) {
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
		depCollectorGenerator: configs.DepCollectorGenerator,
		force:                 configs.Force,
	}, nil
}

// Collect Collect the collectableFile
func (c *AsyncCollector) Collect(ctx context.Context) <-chan error {
	complete := make(chan error)
	go c.collectFileAsync(ctx, complete)
	return complete
}

func (c *AsyncCollector) collectFileAsync(ctx context.Context, complete chan<- error) {
	cancelCF, err := collectable.WithCancel(ctx, c.collectableFile)
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
		depObjDir := utils.GetTargetResourcesDirPath(c.objectKey)

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
func NewLocalAsyncCollector(cf collectable.FileOperator, base, objectKey string, depURIMapper collectable.URIMapper, options ...Option) (*AsyncCollector, error) {
	collector, err := NewAsyncCollector(
		cf, base, objectKey, LocalFileFreshValidator, LocalMover, depURIMapper, options...)
	if err != nil {
		return nil, err
	}
	return collector, nil
}

// NewOBSAsyncCollector Returns a AsyncCollector which collect the file to OBS
func NewOBSAsyncCollector(bucket provider.Bucket, cf collectable.FileOperator, base, objectKey string, depURIMapper collectable.URIMapper, options ...Option) (*AsyncCollector, error) {
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
