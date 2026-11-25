package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/thekhanj/tci"
	"github.com/urfave/cli/v3"
)

type Cli struct{}

var VERSION = "dev"

var VersionCmd *cli.Command = &cli.Command{
	Name: "version",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println(VERSION)
		return nil
	},
}

var ConfigCmd *cli.Command = &cli.Command{
	Name: "config",
	Commands: []*cli.Command{
		{
			Name: "schema",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				var buf bytes.Buffer
				if err := json.Indent(&buf, tci.ConfigSchema, "", "  "); err != nil {
					return err
				}

				_, err := fmt.Fprint(os.Stdout, buf.String())
				return err
			},
		},
	},
}

func (c Cli) Command() cli.Command {
	return cli.Command{
		Name:                  "tci",
		Usage:                 "Manage internet.tci.ir panel",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			ConfigCmd,
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
