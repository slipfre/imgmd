package collector

import (
	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/provider"
)

// LocalCollectorGenerator Generate local collectors
func LocalCollectorGenerator(cf collectable.FileOperator, base, objectKey string, depGenerator Generator, options ...Option) (Collector, error) {
	collector, err := NewAsyncCollector(
		cf, base, objectKey, LocalFileFreshValidator, LocalMover, collectable.LocalURIMapper, depGenerator, options...)
	if err != nil {
		return nil, err
	}
	return collector, nil
}

// GetOBSCollectorGenerator Return a OBS collector generator which generate obs collectors
func GetOBSCollectorGenerator(bucket provider.Bucket) Generator {
	return func(cf collectable.FileOperator, base, objectKey string, generator Generator, options ...Option) (Collector, error) {
		validator, err := GetOBSFileFreshValidator(bucket)
		if err != nil {
			return nil, err
		}
		mover, err := GetOBSMover(bucket)
		if err != nil {
			return nil, err
		}
		mapper, err := collectable.GetOBSURIMapper(bucket)
		if err != nil {
			return nil, err
		}
		collector, err := NewAsyncCollector(
			cf, base, objectKey, validator, mover, mapper, generator, options...)
		if err != nil {
			return nil, err
		}
		return collector, nil
	}
}

// GetPartOBSCollectorGenerator Returns a collector generator which generate obs
// collector for specific file types
func GetPartOBSCollectorGenerator(bucket provider.Bucket, types map[collectable.FileType]struct{}) Generator {
	return func(cf collectable.FileOperator, base, objectKey string, generator Generator, options ...Option) (Collector, error) {
		if types == nil {
			return LocalCollectorGenerator(cf, base, objectKey, generator, options...)
		}
		if _, ok := types[cf.GetFileType()]; ok {
			return GetOBSCollectorGenerator(bucket)(cf, base, objectKey, generator, options...)
		}
		return LocalCollectorGenerator(cf, base, objectKey, generator, options...)
	}
}

// GetLocalCollectorGenerator Returns a collector generator which collect to
// local and collect dependencies to obs
func GetLocalCollectorGenerator(urimapper collectable.URIMapper) Generator {
	return func(cf collectable.FileOperator, base, objectKey string, generator Generator, options ...Option) (Collector, error) {
		collector, err := NewAsyncCollector(
			cf, base, objectKey, LocalFileFreshValidator, LocalMover, urimapper, generator, options...)
		if err != nil {
			return nil, err
		}
		return collector, nil
	}
}
