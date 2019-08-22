// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Devices struct {
	*Conch
}

func (c *Conch) Devices() *Devices {
	return &Devices{c}
}

/***/

type DeviceSettings map[string]interface{}

func (ds DeviceSettings) String() string {
	if API.JsonOnly {
		return API.AsJSON(ds)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	var keys []string
	for key := range ds {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := ds[key]
		table.Append([]string{key, value.(string)})
	}

	table.Render()
	return tableString.String()
}

func (d Devices) Setting(id string, key string) interface{} {
	uri := fmt.Sprintf(
		"/device/%s/settings/%s",
		url.PathEscape(id),
		url.PathEscape(key),
	)

	// The json schema for a DeviceSetting is basically "A DeviceSettings but with only one key"
	var settings DeviceSettings

	res := d.Do(d.Sling().New().Get(uri))
	if ok := res.Parse(&settings); !ok {
		panic(res)
	}
	return settings[key]
}

func (d Devices) Settings(id string) (settings DeviceSettings) {
	uri := fmt.Sprintf("/device/%s/settings", url.PathEscape(id))
	res := d.Do(d.Sling().New().Get(uri))
	if ok := res.Parse(&settings); !ok {
		panic(res)
	}
	return settings
}

/***/

type DeviceReport map[string]interface{}

/***/

type DeviceLocation struct {
	Datacenter            Datacenter `json:"datacenter"`
	Room                  Room       `json:"datacenter_room"`
	Rack                  Rack       `json:"rack"`
	RackUnitStart         int        `json:"rack_unit_start"`
	TargetHardwareProduct struct {
		ID     uuid.UUID `json:"id"`
		Name   string    `json:"name"`
		Alias  string    `json:"alias"`
		Vendor string    `json:"vendor"`
	} `json:"target_hardware_product"`
}

type deviceCore struct {
	ID       uuid.UUID `json:"id"`
	Serial   string    `json:"serial_number"`
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
}

type Disk struct {
	ID           uuid.UUID `json:"id"`
	SerialNumber string    `json:"serial_number"`
	Slot         int       `json:"slot,omitempty"`
	Size         int       `json:"size,omitempty"`
	Vendor       string    `json:"vendor,omitempty"`
	Model        string    `json:"model,omitempty"`
	Firmware     string    `json:"firmware,omitempty"`
	Transport    string    `json:"transport,omitempty"`
	Health       string    `json:"health,omitempty"`
	DriveType    string    `json:"drive_type,omitempty"`
	Enclosure    int       `json:"enclosure,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}
type Disks []Disk

type DetailedDevice struct {
	deviceCore
	Links    []string       `json:"links"`
	Location DeviceLocation `json:"location,omitempty"`
	Nics     []struct {
		Mac             string `json:"mac"`
		InterfaceName   string `json:"iface_name"`
		InterfaceVendor string `json:"iface_vendor"`
		InterfaceType   string `json:"iface_type"`
		PeerMac         string `json:"peer_mac,omitempty"`
		PeerSwitch      string `json:"peer_switch,omitempty`
	} `json:"nics"`
	Disks        Disks        `json:"disks"`
	LatestReport DeviceReport `json:"latest_report,omitempty"`
}

func (d DetailedDevice) String() string {
	if API.JsonOnly {
		return API.AsJSON(d)
	}

	enclosures := make(map[int]map[int]Disk)
	for _, disk := range d.Disks {
		enclosure, ok := enclosures[disk.Enclosure]
		if !ok {
			enclosure = make(map[int]Disk)
		}

		if _, ok := enclosure[disk.Slot]; !ok {
			enclosure[disk.Slot] = disk
		}

		enclosures[disk.Enclosure] = enclosure
	}

	var rackRole RackRole
	if (d.Location.Rack.RoleID != uuid.UUID{}) {
		rackRole = API.RackRoles().Get(d.Location.Rack.RoleID)
	}

	var hp HardwareProduct
	if (d.HardwareProductID != uuid.UUID{}) {
		hp = API.Hardware().GetProduct(d.HardwareProductID)
	}

	extended := struct {
		DetailedDevice
		RackRole        RackRole
		HardwareProduct HardwareProduct
		Enclosures      map[int]map[int]Disk
	}{d, rackRole, hp, enclosures}

	t, err := template.New("d").Parse(deviceTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, extended); err != nil {
		panic(err)
	}

	return buf.String()
}

/***/

type Device struct {
	deviceCore
	RackID        uuid.UUID `json:"rack_id,omitempty"`
	RackUnitStart int       `json:"rack_unit_start,omitempty"`
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

// id is a string because the API accepts both a UUID and a serial number
func (ds *Devices) Get(id string) (d DetailedDevice) {
	uri := fmt.Sprintf("/device/%s", url.PathEscape(id))
	res := ds.Do(ds.Sling().New().Get(uri))
	if ok := res.Parse(&d); !ok {
		panic(res)
	}
	return d
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

func init() {
	App.Command("device", "Perform actions against a single device", func(cmd *cli.Cmd) {
		idArg := cmd.StringArg(
			"DEVICE",
			"",
			"UUID or serial number of the device",
		)

		cmd.Spec = "DEVICE"

		cmd.Command("get", "Get information about a single device", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Devices().Get(*idArg)) }
		})

		cmd.Command("settings", "See all settings for a device", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Devices().Settings(*idArg)) }
		})

		cmd.Command("setting", "See a single setting for a device", func(cmd *cli.Cmd) {
			keyArg := cmd.StringArg(
				"NAME",
				"",
				"Name of the setting",
			)

			cmd.Spec = "NAME"

			cmd.Action = func() {
				fmt.Println(API.Devices().Setting(*idArg, *keyArg))
			}
		})
	})

}
