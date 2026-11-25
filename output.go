package tci

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type ColumnSet string

const (
	ColumnsDefault ColumnSet = "default"
	ColumnsWide              = "wide"
)

type OutputFormat string

const (
	OutputTsv  OutputFormat = "tsv"
	OutputJson              = "json"
)

type Tabular interface {
	Columns(columns ColumnSet) []string
	Values(columns []string) []string
}

func TableWriteJson[T Tabular](
	w io.Writer, arr []T,
) error {
	return json.NewEncoder(w).Encode(arr)
}

func TableWriteTsv[T Tabular](
	w io.Writer, columns []string, arr []T,
) error {
	if len(arr) == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, strings.Join(columns, "\t"))
	for _, item := range arr {
		fmt.Fprintln(tw, strings.Join(item.Values(columns), "\t"))
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
