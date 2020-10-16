package tables

import (
	"io"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func NewTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	TableToMarkdown(table)
	return table
}

func TableToMarkdown(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
}

type Tabulable interface {
	Headers() []string
	ForEach(func([]string))
	sort.Interface
}

func Render(list Tabulable) string {
	sort.Sort(list)

	tableString := &strings.Builder{}
	table := NewTable(tableString)

	table.SetHeader(list.Headers())
	list.ForEach(table.Append)

	table.Render()
	return tableString.String()
}
