package cmd

import (
	"errors"

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
