package main

import (
	"context"
	"fmt"
	"os"

	tcicli "github.com/thekhanj/tci/internal/cli"
	"github.com/thekhanj/tci/internal/config"
	"github.com/thekhanj/tci/panel"
	"github.com/urfave/cli/v3"
)

type Cli struct{}

func (c Cli) Cmd() cli.Command {
	return cli.Command{
		Name:                  "tci",
		Usage:                 "Manage internet.tci.ir panel",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			c.StatusCmd(),
			ConfigCmd,
			VersionCmd,
		},
	}
}

func (c Cli) StatusCmd() *cli.Command {
	var (
		configPath string
		output     tcicli.OutputOptions
	)

	return &cli.Command{
		Name:  "status",
		Usage: "Display panel status",
		Flags: append(
			tcicli.OutputFlags(&output),
			tcicli.ConfigFlag(&configPath),
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return tcicli.Status(ctx, configPath, output)
		},
	}
}

func main() {
	cmd := (Cli{}).Cmd()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
