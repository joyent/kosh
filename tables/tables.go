package tables

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func NewTable(writer io.Writer) *tablewriter.Table {
	return tablewriter.NewWriter(writer)
}

func TableToMarkdown(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
}
