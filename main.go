package main

import (
	"fmt"
	"os"

	"github.com/slipfre/imgmd/cmd"
	"github.com/slipfre/imgmd/collectable"
	"github.com/urfave/cli/v2"
)

var typeFlag = &cli.StringSliceFlag{
	Name:    "type",
	Aliases: []string{"t"},
	Value:   cli.NewStringSlice(string(collectable.Markdown)),
	Usage:   "Types of file you want to collect",
}

var recursiveFlag = &cli.BoolFlag{
	Name:    "recursive",
	Aliases: []string{"r"},
	Value:   false,
	Usage:   "Copy or move all eligible files under specified directory and its sub-directories",
}

var configFlag = &cli.PathFlag{
	Name:    "config",
	Aliases: []string{"c"},
	Value:   "~/.cresrc.yml",
	Usage:   "Path of config file",
}

var dep2obsFlag = &cli.StringSliceFlag{
	Name:    "dep2obs",
	Aliases: []string{"d2o"},
	Value:   cli.NewStringSlice(),
	Usage:   "Types of dependency files which want to put in obs",
}

func main() {
	app := &cli.App{
		Name:  "cres",
		Usage: "Collect resources such as markdown documents which contain links to other resources and independent resources such as .png format images",
		Flags: []cli.Flag{
			typeFlag,
			recursiveFlag,
			configFlag,
			dep2obsFlag,
		},
		Commands: []*cli.Command{
			cmd.MoveCommand,
			cmd.CopyCommand,
			cmd.CollectCommand,
		},
		// Action: func(c *cli.Context) error {
		// 	fmt.Printf("Hello %q", c.Args().Get(0))
		// 	return nil
		// },
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}
