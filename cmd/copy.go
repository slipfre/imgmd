package cmd

import (
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
	source, destination, err := parseCopyArguments(c)
	if err != nil {
		return err
	}

	if err := copyMDs(c, source, destination); err != nil {
		return err
	}

	return nil
}
