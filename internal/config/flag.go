package config

import "github.com/urfave/cli/v3"

func ConfigFlag(path *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Path to configuration file",
		Value:       DefaultConfigPath(),
		Destination: path,
	}
}
