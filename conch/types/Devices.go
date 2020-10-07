package types

import (
	"sort"
	"strings"

	"github.com/joyent/kosh/tables"
	"github.com/joyent/kosh/template"
)

func (d Device) JSON() ([]byte, error) {
	return AsJSON(d)
}

func (d Device) String() string {
	return "[ comming soon ]"
}

func (d Devices) JSON() ([]byte, error) {
	return AsJSON(d)
}

func (d Devices) Len() int {
	return len(d)
}

func (d Devices) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Devices) Less(i, j int) bool {
	return d[i].SerialNumber < d[j].SerialNumber
}

func (d Devices) String() string {
	sort.Sort(d)

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
	tables.TableToMarkdown(table)

	table.SetHeader([]string{
		"Serial",
		"Hostname",
		"Asset Tag",
		"Hardware",
		"Phase",
		"Updated",
		"Validated",
	})

	for _, device := range d {
		table.Append([]string{
			string(device.SerialNumber),
			device.Hostname,
			device.AssetTag,
			device.HardwareProductID.String(),
			string(device.Phase),
			template.TimeStr(device.Updated),
			device.Validated,
		})
	}

	table.Render()
	return tableString.String()
}

func (d DevicePhase) JSON() ([]byte, error) { return AsJSON(d) }
func (d DevicePhase) String() string        { return string(d) }
