package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/slipfre/imgmd/cmd/conf"
	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/collector"
	"github.com/urfave/cli/v2"
)

// CopyCommand cres copy
var CopyCommand = &cli.Command{
	Name:      "copy",
	Usage:     "Copy resources to specified place",
	Action:    copy,
	ArgsUsage: "[source] [destination]",
}

func copy(c *cli.Context) error {
	fail := 0
	success := 0

	source, destination, err := parseArguments(c)
	if err != nil {
		return err
	}

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

	var collectorGenerator = collector.LocalCollectorGenerator
	if dep2obsFlag != nil && len(dep2obsFlag) > 0 {
		if bucket, err := conf.GetBucketFromConfigFile(config); err != nil {
			collectorGenerator = collector.GetOBSCollectorGenerator(bucket)
		}
	}

	collectors := []collector.Collector{}
	if recursive {
		if collectors, err = getCollectorsRecursively(source, destination, collectorGenerator); err != nil {
			return err
		}
	} else {
		collectableFile := collectable.NewMarkdownFile("", source)
		c, err := collectorGenerator(
			collectableFile,
			filepath.Dir(destination),
			filepath.Base(destination),
			collector.LocalCollectorGenerator,
		)
		if err != nil {
			return err
		}
		collectors = append(collectors, c)
	}

	cases := make([]reflect.SelectCase, len(collectors))
	for i := 0; i < len(cases); {
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
			fail++
		}
		success++
		cases[chosen].Chan = reflect.ValueOf(nil)
	}

	fmt.Printf("Finished! Total: %d, success: %d, failed: %d\n", fail+success, success, fail)

	return nil
}
