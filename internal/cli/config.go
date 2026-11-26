package cli

import (
	"github.com/thekhanj/tci/internal/config"
	"github.com/urfave/cli/v3"
)

func ConfigFlag(path *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Path to configuration file",
		Value:       config.DefaultConfigPath(),
		Destination: path,
	}
}
