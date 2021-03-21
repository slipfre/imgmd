package cmd

import (
	"errors"
	"path/filepath"

	"github.com/slipfre/imgmd/cmd/conf"
	"github.com/urfave/cli/v2"
)

var moveFlag = &cli.BoolFlag{
	Name:    "move",
	Aliases: []string{"m"},
	Value:   false,
	Usage:   "Move files rather than copy files. Delete source files after collectation",
}

// CollectCommand cres collect
var CollectCommand = &cli.Command{
	Name:   "collect",
	Usage:  "Collect resources to repository",
	Flags:  []cli.Flag{moveFlag},
	Action: collect,
}

func collect(c *cli.Context) error {
	source, key, err := parseCopyArguments(c)
	if err != nil {
		return err
	}
	confFilePath := c.Path("config")
	if confFilePath == "" {
		return errors.New("config file not specified")
	}
	repoPath, err := conf.GetRepoPathFromConfig(confFilePath)
	if err != nil {
		return err
	}
	dest := filepath.Join(repoPath, key)
	if err := copyMDs(c, source, dest); err != nil {
		return err
	}
	return nil
}
