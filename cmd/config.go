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
