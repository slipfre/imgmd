package internal

import (
	"fmt"
	"path"
	"strings"
)

// Collector Collect collecatable files
type Collector interface {
	Collect() ([]Collector, error)
}

// Listener Listen certain kind of event
type Listener interface {
	OnEvent(attrs FileAttributesGetter, err error)
}

// CollectorWithListener Collector which contains listener
type CollectorWithListener struct {
	errListener      Listener
	completeListener Listener
	collectableFile  CollectableFileOperator
	targetPath       string
}

// CollectorWithListenerConfig Config for CollectorWithListener
type CollectorWithListenerConfig struct {
	errListener      Listener
	completeListener Listener
}

// CollectorWithListenerOption Option configs for CollectorWithListener
type CollectorWithListenerOption func(config *CollectorWithListenerConfig)

// WithErrListener ErrListener option for CollectorWithListener
func WithErrListener(errListener Listener) CollectorWithListenerOption {
	return func(config *CollectorWithListenerConfig) {
		config.errListener = errListener
	}
}

// WithCompleteListener CompleteListener option for CollectorWithListener
func WithCompleteListener(completeListenr Listener) CollectorWithListenerOption {
	return func(config *CollectorWithListenerConfig) {
		config.completeListener = completeListenr
	}
}

// NewCollectorWithListener Constructor for CollectorWithListener
func NewCollectorWithListener(cf CollectableFileOperator, targetPath string, options ...CollectorWithListenerOption) *CollectorWithListener {
	config := &CollectorWithListenerConfig{}
	for _, option := range options {
		option(config)
	}

	return &CollectorWithListener{
		errListener:      config.errListener,
		completeListener: config.completeListener,
		collectableFile:  cf,
		targetPath:       targetPath,
	}
}

// Collect Collect the collectableFile
func (c *CollectorWithListener) Collect() ([]Collector, error) {
	dependencies, err := c.collectFile()
	if err != nil {
		if c.errListener != nil {
			c.errListener.OnEvent(c.collectableFile, err)
		}
		return nil, err
	}
	if c.completeListener != nil {
		c.completeListener.OnEvent(c.collectableFile, err)
	}
	return dependencies, err
}

func (c *CollectorWithListener) collectFile() ([]Collector, error) {
	deps, err := c.collectableFile.FindDependencies()
	if err != nil {
		return nil, err
	}

	targetResourcesDirPath := c.getTargetResourcesDirPath()
	if err = CreateDirectory(targetResourcesDirPath); err != nil {
		return nil, err
	}

	var depCollectors []Collector
	for _, dep := range deps {
		collector := NewCollectorWithListener(
			dep,
			fmt.Sprintf("%s/%s", targetResourcesDirPath, path.Base(dep.GetURI())),
			WithErrListener(c.errListener),
			WithCompleteListener(c.completeListener),
		)
		depCollectors = append(depCollectors, collector)
	}

	err = c.collectableFile.ReplaceDependencyURIs(func(fileType FileType, uri []byte) []byte {
		dirName := path.Base(targetResourcesDirPath)
		newReferencePath := fmt.Sprintf("%s/%s", dirName, string(uri))
		return []byte(newReferencePath)
	})
	if err != nil {
		return nil, err
	}

	if err = c.collectableFile.To(c.targetPath); err != nil {
		return nil, err
	}

	return depCollectors, err
}

func (c *CollectorWithListener) getTargetResourcesDirPath() string {
	targetURI := c.collectableFile.GetURI()
	directory := path.Dir(targetURI)
	filenameWithSuffix := path.Base(targetURI)
	suffix := path.Ext(filenameWithSuffix)
	filename := strings.TrimSuffix(filenameWithSuffix, suffix)
	return fmt.Sprintf("%s%s_medias", directory, filename)
}
