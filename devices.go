// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/olekukonko/tablewriter"
)

type Devices struct {
	*Conch
}

func (c *Conch) Devices() *Devices {
	return &Devices{c}
}

/***/

type Device struct {
	// ID uuid.UUID `json:"id"`
	// Serial string `json:"serial_number"`
	Serial   string    `json:"id"`
	AssetTag string    `json:"asset_tag,omitempty"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	LastSeen time.Time `json:"last_seen"`

	HardwareProductID uuid.UUID `json:"hardware_product_id"`
	Health            string    `json:"health"`
	Hostname          string    `json:"hostname,omitempty"`
	SystemUUID        uuid.UUID `json:"system_uuid"`
	UptimeSince       time.Time `json:"uptime_since,omitempty"`
	Validated         time.Time `json:"validated,omitempty"`
	Phase             string    `json:"phase"`
	RackID            uuid.UUID `json:"rack_id,omitempty"`
	RackUnitStart     int       `json:"rack_unit_start,omitempty"`
}

type DeviceList []Device

func (d DeviceList) Len() int {
	return len(d)
}

func (d DeviceList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DeviceList) Less(i, j int) bool {
	return d[i].Serial < d[j].Serial
}

func (d DeviceList) String() string {
	sort.Sort(d)
	if API.JsonOnly {
		return API.AsJSON(d)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	// TODO(sungo): rack, hardware product
	table.SetHeader([]string{
		"Serial",
		"Hostname",
		"Asset Tag",
		"Phase",
		"Updated",
		"Validated",
	})

	for _, device := range d {
		table.Append([]string{
			device.Serial,
			device.Hostname,
			device.AssetTag,
			device.Phase,
			TimeStr(device.Updated),
			TimeStr(device.Validated),
		})
	}

	table.Render()
	return tableString.String()
}

/***/

var HealthList = []string{"error", "fail", "unknown", "pass"}

func prettyDeviceHealthList() string {
	return strings.Join(HealthList, ", ")
}

func okHealth(health string) bool {
	for _, b := range HealthList {
		if health == b {
			return true
		}
	}
	return false
}

/***/
