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

type Table interface {
	Headers(OutputFormat) []string
	Values(fields []string) []string
}

func TableOutput[T Table](w io.Writer, format OutputFormat, arr []T) error {
	if format == OutputJson {
		return json.NewEncoder(w).Encode(arr)
	}

	if len(arr) == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, strings.Join(arr[0].Headers(format), "\t"))

	headers := arr[0].Headers(format)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, item := range arr {
		fmt.Fprintln(tw, strings.Join(item.Values(headers), "\t"))
	}

	return tw.Flush()
}

type FieldValue struct {
	Value string
	Lazy  func() string
}

func (v FieldValue) String() string {
	if v.Lazy != nil {
		return v.Lazy()
	}
	return v.Value
}

func NewFieldValue(value string) FieldValue {
	return FieldValue{Value: value}
}

func NewLazyFieldValue(fn func() string) FieldValue {
	return FieldValue{Lazy: fn}
}

func SelectFields(fields []string, values map[string]FieldValue) []string {
	row := make([]string, 0, len(fields))

	for _, field := range fields {
		row = append(row, values[field].String())
	}

	return row
}
