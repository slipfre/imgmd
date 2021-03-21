package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/slipfre/imgmd/cmd/conf"
	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/collector"
	"github.com/urfave/cli/v2"
)

func parseCopyArguments(c *cli.Context) (source, destOrKey string, err error) {
	source = c.Args().Get(0)
	destOrKey = c.Args().Get(1)
	if source == "" {
		err = errors.New("source must be specified")
	}
	if destOrKey == "" {
		err = errors.New("destination or key must be specified")
	}
	return
}

func parseCollectArguments(c *cli.Context) (key string, err error) {
	key = c.Args().Get(0)
	if key == "" {
		err = errors.New("key must be specified")
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

func getCollectorsRecursively(source, destination string, generator collector.Generator, uriMapper collectable.URIMapper) (collectors []collector.Collector, err error) {
	collectors = []collector.Collector{}
	sourceAbsolute, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		collectableFile := collectable.NewMarkdownFile("", path)
		pathAbsolute, err := filepath.Abs(path)
		if err != nil {
			return nil
		}
		key := strings.TrimPrefix(pathAbsolute, sourceAbsolute)
		key = strings.TrimPrefix(key, "\\")
		key = strings.TrimPrefix(key, "/")
		c, err := collector.GetLocalCollectorGenerator(uriMapper)(
			collectableFile,
			destination,
			key,
			generator,
		)
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

func copyMDs(c *cli.Context, source, destination string) (err error) {
	fail := 0
	success := 0
	_, recursive, config, dep2obsFlag := parseGlobalFlags(c)

	if recursive {
		if errStr := validateDir(source); errStr != "" {
			return errors.New(errStr)
		}
	} else {
		if errStr := validateFile(source); errStr != "" {
			return errors.New(errStr)
		}
	}

	var depCollectorGenerator = collector.LocalCollectorGenerator
	var depURIMapper = collectable.LocalURIMapper
	if dep2obsFlag != nil && len(dep2obsFlag) > 0 {
		if bucket, err := conf.GetBucketFromConfigFile(config); err == nil {
			depCollectorGenerator = collector.GetOBSCollectorGenerator(bucket)
			if depURIMapper, err = collectable.GetOBSURIMapper(bucket); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	collectors := []collector.Collector{}
	if recursive {
		if collectors, err = getCollectorsRecursively(source, destination, depCollectorGenerator, depURIMapper); err != nil {
			return err
		}
	} else {
		collectableFile := collectable.NewMarkdownFile("", source)
		c, err := collector.GetLocalCollectorGenerator(depURIMapper)(
			collectableFile,
			filepath.Dir(destination),
			filepath.Base(destination),
			depCollectorGenerator,
		)
		if err != nil {
			return err
		}
		collectors = append(collectors, c)
	}

	cases := make([]reflect.SelectCase, len(collectors))
	for i := 0; i < len(cases); i++ {
		complete := collectors[i].Collect(context.Background())
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(complete),
		}
	}

	for remaining := len(cases); remaining > 0; remaining-- {
		chosen, value, _ := reflect.Select(cases)
		if value.Interface() != nil {
			err = value.Interface().(error)
		}
		if err != nil {
			log.Printf(err.Error())
			fail++
		} else {
			success++
		}
		cases[chosen].Chan = reflect.ValueOf(nil)
	}

	log.Printf("Finished! Total: %d, success: %d, failed: %d\n", fail+success, success, fail)

	return nil
}
