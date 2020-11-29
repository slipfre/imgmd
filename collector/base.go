package collector

import (
	"context"

	"github.com/slipfre/imgmd/collectable"
)

// Collector Collect collecatable files
type Collector interface {
	Collect(ctx context.Context) <-chan error
}

// FreshValidator Validate whether the local file is up to date
type FreshValidator func(cf collectable.FileOperator, base, objectKey string) (bool, error)

// Mover Make files to specified place
type Mover func(cf collectable.FileOperator, base, objectKey string) error

// Generator Generate collectors
type Generator func(cf collectable.FileOperator, base, objectKey string, depURIMapper collectable.URIMapper, options ...Option) (Collector, error)

// Configs Configurations for collector
type Configs struct {
	Force                 bool
	DepCollectorGenerator Generator
}

// Option Options for collectors
type Option func(configs *Configs)

// WithForce Option config for collectors. If force is true, the file will be
// collected even though the target file exists and up-to-date
func WithForce(force bool) Option {
	return func(configs *Configs) {
		configs.Force = force
	}
}

// WithDependencyGenerator Option config for collectors
func WithDependencyGenerator(generator Generator) Option {
	return func(configs *Configs) {
		configs.DepCollectorGenerator = generator
	}
}
