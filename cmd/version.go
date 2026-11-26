package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var VERSION = "dev"

var VersionCmd *cli.Command = &cli.Command{
	Name: "version",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println(VERSION)
		return nil
	},
}
