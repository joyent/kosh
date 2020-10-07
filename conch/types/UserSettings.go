package types

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/joyent/kosh/tables"
)

func (u UserSettings) JSON() ([]byte, error) {
	return AsJSON(u)
}

func (u UserSettings) String() string {
	builder := &strings.Builder{}
	u.RenderTable(builder)
	return builder.String()
}

func (u UserSettings) RenderTable(writer io.Writer) {
	keys := make([]string, 0)
	for setting := range u {
		keys = append(keys, setting)
	}
	sort.Strings(keys)

	table := tables.NewTable(writer)
	tables.TableToMarkdown(table)
	table.SetHeader([]string{
		"Key",
		"Value",
	})

	for _, key := range keys {
		table.Append([]string{
			key,
			fmt.Sprintf("%v", u[key]),
		})
	}

	table.Render()
}

func (u UserSetting) JSON() ([]byte, error) {
	return AsJSON(u)
}

func (u UserSetting) String() string {
	return string(u)
}
