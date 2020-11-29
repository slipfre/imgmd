package collector

import (
	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/provider"
)

// LocalCollectorGenerator Generate local collectors
func LocalCollectorGenerator(cf collectable.FileOperator, base, objectKey string, depURIMapper collectable.URIMapper, options ...Option) (Collector, error) {
	return NewLocalAsyncCollector(cf, base, objectKey, depURIMapper, options...)
}

// GetOBSCollectorGenerator Return a OBS collector generator which generate obs collectors
func GetOBSCollectorGenerator(bucket provider.Bucket) Generator {
	return func(cf collectable.FileOperator, base, objectKey string, depURIMapper collectable.URIMapper, options ...Option) (Collector, error) {
		return NewOBSAsyncCollector(bucket, cf, base, objectKey, depURIMapper, options...)
	}
}
