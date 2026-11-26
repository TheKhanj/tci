package cli

import (
	"github.com/urfave/cli/v3"
)

type OutputOptions struct {
	Wide    bool
	Columns []string
	Format  string
}

func OutputFlags(opts *OutputOptions) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "wide",
			Aliases:     []string{"w"},
			Usage:       "Display all available columns",
			Destination: &opts.Wide,
		},
		&cli.StringSliceFlag{
			Name:        "columns",
			Aliases:     []string{"c"},
			Usage:       "Columns to display",
			Destination: &opts.Columns,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "Output format (tsv, json)",
			Value:       "tsv",
			Destination: &opts.Format,
		},
	}
}
