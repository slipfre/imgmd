package cmd

import (
	"github.com/slipfre/imgmd/cmd/conf"
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
	source, destination, err := parseArguments(c)
	if err != nil {
		return err
	}

	types, recursive, config, dep2obsFlag := parseGlobalFlags(c)

	bucket, err := conf.ParseConfigFile(config)
	if err != nil {
		return err
	}
	// TODO: 若设置了 dep2obs，解析配置文件，获取 bucket
	// TODO: 根据传入参数创建 collector 列表
	// TODO: collect
	return nil
}
