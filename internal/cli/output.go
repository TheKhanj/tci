package cli

import (
	"github.com/thekhanj/tci"
	"github.com/urfave/cli/v3"
)

type OutputOptions struct {
	Wide            bool
	SelectedColumns []string
	Format          string
}

func (o OutputOptions) Headers(headersFn tci.HeadersFn) []string {
	if len(o.SelectedColumns) != 0 {
		return o.SelectedColumns
	}
	if o.Wide {
		return headersFn(tci.ColumnsWide)
	}

	return headersFn(tci.ColumnsDefault)
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
			Usage:       "Columns to display",
			Destination: &opts.SelectedColumns,
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
