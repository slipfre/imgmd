package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// MoveCommand cres move
var MoveCommand = &cli.Command{
	Name:  "move",
	Usage: "Move resources to specified place",
	Action: func(c *cli.Context) error {
		fmt.Printf("Hello %q", c.Args().Get(0))
		return nil
	},
}
