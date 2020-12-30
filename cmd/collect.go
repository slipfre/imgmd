package cmd

import (
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
	return nil
}
