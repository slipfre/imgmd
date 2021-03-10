package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/collector"
	"github.com/urfave/cli/v2"
)

func parseArguments(c *cli.Context) (source, destOrKey string, err error) {
	source = c.Args().Get(0)
	destOrKey = c.Args().Get(1)
	if source == "" {
		err = errors.New("source must be specified")
	}
	if destOrKey == "" {
		err = errors.New("destination/key must be specified")
	}
	return
}

func parseGlobalFlags(c *cli.Context) (types []string, recursive bool, config string, dep2obs []string) {
	types = c.StringSlice("type")
	recursive = c.Bool("recursive")
	config = c.Path("config")
	dep2obs = c.StringSlice("dep2obs")
	return
}

func getCollectableFileRecursively(dirname string) (collectableFiles []collectable.FileOperator, err error) {
	collectableFiles = []collectable.FileOperator{}
	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		collectableFiles = append(collectableFiles, collectable.NewMarkdownFile("", path))
		return nil
	})
	return collectableFiles, err
}

func getCollectorsRecursively(source, destination string, generator collector.Generator) (collectors []collector.Collector, err error) {
	collectors = []collector.Collector{}
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		collectableFile := collectable.NewMarkdownFile("", path)
		key := strings.TrimPrefix(path, source)
		c, err := generator(collectableFile, destination, key, collector.LocalCollectorGenerator)
		if err != nil {
			return err
		}
		collectors = append(collectors, c)
		return nil
	})
	return collectors, err
}

func validateDir(path string) string {
	s, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}

	if !s.IsDir() {
		return fmt.Sprintf("'%s' is not a directory", path)
	}

	return ""
}

func validateFile(path string) string {
	s, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}

	if s.IsDir() {
		return fmt.Sprintf("'%s' is a directory", path)
	}

	return ""
}
