package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

type Cli struct{}

const VERSION = "dev"

var VersionCmd *cli.Command = &cli.Command{
	Name: "version",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println(VERSION)
		return nil
	},
}

func (c Cli) Command() cli.Command {
	return cli.Command{
		Name:                  "tci",
		Usage:                 "Manage internet.tci.ir panel",
		EnableShellCompletion: true,
		// TODO: remove config flag from the other places
		Commands: []*cli.Command{
			VersionCmd,
		},
	}
}

func main() {
	cmd := (Cli{}).Command()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
