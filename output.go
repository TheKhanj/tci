package tci

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type ColumnSet int

const (
	ColumnsDefault ColumnSet = iota
	ColumnsWide
)

type Column struct {
	Header string
	Set    ColumnSet
}

type Columns []Column

func (cs Columns) Headers(set ColumnSet) []string {
	ret := make([]string, 0)
	for _, c := range cs {
		if c.Set <= set {
			ret = append(ret, c.Header)
		}
	}
	return ret
}

func NewColumn(name string, set ColumnSet) Column {
	return Column{name, set}
}

type OutputFormat string

const (
	OutputTsv  OutputFormat = "tsv"
	OutputJson              = "json"
)

type Tabular interface {
	Columns() Columns
	Values(names []string) []string
}

func TableWriteJson[T Tabular](
	w io.Writer, arr []T,
) error {
	return json.NewEncoder(w).Encode(arr)
}

func TableWriteTsv[T Tabular](
	w io.Writer, names []string, arr []T,
) error {
	if len(arr) == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, strings.Join(names, "\t"))
	for _, item := range arr {
		fmt.Fprintln(tw, strings.Join(item.Values(names), "\t"))
	}

	return tw.Flush()
}

type Cell struct {
	Value string
	Lazy  func() string
}

func (v Cell) String() string {
	if v.Lazy != nil {
		return v.Lazy()
	}
	return v.Value
}

func NewCell(value string) Cell {
	return Cell{Value: value}
}

func NewLazyCell(fn func() string) Cell {
	return Cell{Lazy: fn}
}

func SelectCell(fields []string, values map[string]Cell) []string {
	row := make([]string, 0, len(fields))

	for _, field := range fields {
		row = append(row, values[field].String())
	}

	return row
}
