package tci

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type OutputFormat string

const (
	OutputDefault OutputFormat = ""
	OutputWide    OutputFormat = "wide"
	OutputJson    OutputFormat = "json"
)

type Outputter interface {
	OutputHeader(OutputFormat) []string
	OutputRow(OutputFormat) []string
}

func Output[T Outputter](w io.Writer, format OutputFormat, arr []T) error {
	if format == OutputJson {
		return json.NewEncoder(w).Encode(arr)
	}

	if len(arr) == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, strings.Join(arr[0].OutputHeader(format), "\t"))

	for _, item := range arr {
		fmt.Fprintln(tw, strings.Join(item.OutputRow(format), "\t"))
	}

	return tw.Flush()
}
