/*
Package tables renders tabular data.
*/
package tables

import (
	"io"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// NewTable returns a pointer to a configured instance of tableWriter.Table
func NewTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	TableToMarkdown(table)
	return table
}

// TableToMarkdown sets the table up to render in a markdown compatible format
func TableToMarkdown(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
}

// Tabulable is an interface for any structure that can be rendered into a
// table form
type Tabulable interface {
	Headers() []string
	ForEach(func([]string))
	sort.Interface
}

// Render takess a Tabulable struct and renders it into markdown compatible
// string
func Render(list Tabulable) string {
	sort.Sort(list)

	tableString := &strings.Builder{}
	table := NewTable(tableString)

	table.SetHeader(list.Headers())
	list.ForEach(table.Append)

	table.Render()
	return tableString.String()
}
