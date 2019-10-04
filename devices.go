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
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
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

func (d Devices) Setting(id, key string) interface{} {
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

func (d Devices) Settings(id string) DeviceSettings {
	uri := fmt.Sprintf("/device/%s/settings", url.PathEscape(id))
	res := d.Do(d.Sling().New().Get(uri))

	settings := make(DeviceSettings)
	if ok := res.Parse(&settings); !ok {
		panic(res)
	}

	out := make(DeviceSettings)
	re := regexp.MustCompile(`^tag\.`)
	for key, value := range settings {
		if !re.MatchString(key) {
			out[key] = value
		}
	}

	return out
}

func (ds *Devices) SetSetting(id, key, value string) {
	uri := fmt.Sprintf(
		"/device/%s/settings/%s",
		url.PathEscape(id),
		url.PathEscape(key),
	)

	settings := make(DeviceSettings)
	settings[key] = value

	ds.Do(
		ds.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(settings),
	)
}

func (ds *Devices) DeleteSetting(id, key string) {
	uri := fmt.Sprintf(
		"/device/%s/settings/%s",
		url.PathEscape(id),
		url.PathEscape(key),
	)
	ds.Do(ds.Sling().New().Delete(uri))
}

func (d Devices) Tags(id string) DeviceSettings {
	uri := fmt.Sprintf("/device/%s/settings", url.PathEscape(id))
	res := d.Do(d.Sling().New().Get(uri))

	settings := make(DeviceSettings)
	if ok := res.Parse(&settings); !ok {
		panic(res)
	}

	tags := make(DeviceSettings)

	re := regexp.MustCompile(`^tag\.`)
	for key, value := range settings {
		if re.MatchString(key) {
			tag := strings.TrimPrefix(key, "tag.")
			tags[tag] = value
		}
	}

	return tags
}

func (d Devices) Tag(id, key string) interface{} {
	re := regexp.MustCompile(`^tag\.`)
	if !re.MatchString(key) {
		key = "tag." + key
	}
	return d.Setting(id, key)
}

func (d Devices) SetTag(id, key, value string) {
	re := regexp.MustCompile(`^tag\.`)
	if !re.MatchString(key) {
		key = "tag." + key
	}
	d.SetSetting(id, key, value)
}

func (d Devices) DeleteTag(id, key string) {
	re := regexp.MustCompile(`^tag\.`)
	if !re.MatchString(key) {
		key = "tag." + key
	}
	d.DeleteSetting(id, key)
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

func (dl DeviceLocation) String() string {
	if API.JsonOnly {
		return API.AsJSON(dl)
	}

	t, err := NewTemplate().Parse(deviceLocationTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, dl); err != nil {
		panic(err)
	}

	return buf.String()
}

func (ds *Devices) GetLocation(id string) (l DeviceLocation) {
	uri := fmt.Sprintf("/device/%s/location", url.PathEscape(id))
	res := ds.Do(ds.Sling().New().Get(uri))
	if ok := res.Parse(&l); !ok {
		panic(res)
	}
	return l
}

/***/

type DeviceNic struct {
	DeviceID        uuid.UUID `json:"device_id"`
	MAC             string    `json:"mac"`
	InterfaceName   string    `json:"iface_name"`
	InterfaceVendor string    `json:"iface_vendor"`
	InterfaceDriver string    `json:"iface_driver"`
	State           string    `json:"state,omitempty"`
	IpAddress       string    `json:"ipaddr,omitempty"`
	MTU             int       `json:"mtu,omitempty"`

	InterfaceType string `json:"iface_type"`
}

func (dn DeviceNic) String() string {
	if API.JsonOnly {
		return API.AsJSON(dn)
	}

	t, err := NewTemplate().Parse(deviceNicTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, dn); err != nil {
		panic(err)
	}

	return buf.String()
}

// type DeviceNics []DeviceNic

func (ds *Devices) GetInterface(id, name string) (n DeviceNic) {
	uri := fmt.Sprintf(
		"/device/%s/interface/%s",
		url.PathEscape(id),
		url.PathEscape(name),
	)

	res := ds.Do(ds.Sling().New().Get(uri))
	if ok := res.Parse(&n); !ok {
		panic(res)
	}
	return n
}

func (ds *Devices) GetIPMI(id string) string {
	return ds.GetInterface(id, "ipmi1").IpAddress
}

/***/

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

	BuildID uuid.UUID `json:"build_id"`
}

type Disk struct {
	ID           uuid.UUID   `json:"id"`
	SerialNumber string      `json:"serial_number"`
	Slot         int         `json:"slot,omitempty"`
	Size         int         `json:"size,omitempty"`
	Vendor       string      `json:"vendor,omitempty"`
	Model        string      `json:"model,omitempty"`
	Firmware     string      `json:"firmware,omitempty"`
	Transport    string      `json:"transport,omitempty"`
	Health       string      `json:"health,omitempty"`
	DriveType    string      `json:"drive_type,omitempty"`
	Enclosure    int         `json:"enclosure,omitempty"`
	Created      time.Time   `json:"created"`
	Updated      time.Time   `json:"updated"`
	HBA          interface{} `json:"hba"` // TODO figure out where this belongs
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
		PeerSwitch      string `json:"peer_switch,omitempty"`
		PeerPort        string `json:"peer_port,omitempty"`
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

	validations := API.Devices().ValidationState(d.ID.String())

	extended := struct {
		DetailedDevice
		RackRole        RackRole
		HardwareProduct HardwareProduct
		Enclosures      map[int]map[int]Disk
		Validations     ValidationStatesWithResults
	}{d, rackRole, hp, enclosures, validations}

	t, err := NewTemplate().Parse(deviceTemplate)
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

	table.SetHeader([]string{
		"Serial",
		"Hostname",
		"Asset Tag",
		"Hardware",
		"Phase",
		"Updated",
		"Validated",
	})

	hpCache := make(map[uuid.UUID]HardwareProduct)

	for _, device := range d {
		if _, ok := hpCache[device.HardwareProductID]; !ok {
			hpCache[device.HardwareProductID] = API.Hardware().GetProduct(device.HardwareProductID)
		}

		table.Append([]string{
			device.Serial,
			device.Hostname,
			device.AssetTag,
			hpCache[device.HardwareProductID].Name,
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
		panic(fmt.Sprintf("%v", res))
	}
	return d
}

func (ds *Devices) FindByField(key, value string) DeviceList {
	uri := fmt.Sprintf(
		"/device?%s=%s",
		url.PathEscape(key),
		url.PathEscape(value),
	)
	d := make(DeviceList, 0)

	res := ds.Do(ds.Sling().New().Get(uri))
	if ok := res.Parse(&d); !ok {
		panic(res)
	}
	return d
}

func (ds *Devices) FindBySetting(key, value string) DeviceList {
	return ds.FindByField(key, value)
}

func (ds *Devices) FindByTag(key, value string) DeviceList {
	return ds.FindByField("tag."+key, value)
}

/***/

// id is a string because the API accepts both a UUID and a serial number
func (ds *Devices) ValidationState(id string) (v ValidationStatesWithResults) {
	uri := fmt.Sprintf("/device/%s/validation_state", url.PathEscape(id))
	res := ds.Do(ds.Sling().New().Get(uri))
	if ok := res.Parse(&v); !ok {
		panic(res)
	}
	return v
}

/***/

func (ds *Devices) GetPhase(id string) string {
	data := struct {
		ID    uuid.UUID `json:"id"`
		Phase string    `json:"phase"`
	}{}

	uri := fmt.Sprintf("/device/%s/phase", url.PathEscape(id))
	res := ds.Do(ds.Sling().New().Get(uri))
	if ok := res.Parse(&data); !ok {
		panic(res)
	}
	return data.Phase
}

func (ds *Devices) SetPhase(id, phase string) string {
	uri := fmt.Sprintf("/device/%s/phase", url.PathEscape(id))

	payload := make(map[string]string)
	payload["id"] = id
	payload["phase"] = phase

	ds.Do(
		ds.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	return ds.GetPhase(id)
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
var PhasesList = []string{"integration", "installation", "production", "diagnostics", "decommissioned"}

func prettyPhasesList() string {
	return strings.Join(PhasesList, ", ")
}

func okPhase(phase string) bool {
	for _, b := range PhasesList {
		if phase == b {
			return true
		}
	}
	return false
}

/***/

func init() {
	App.Command("devices ds", "Commands for dealing with multiple devices", func(cmd *cli.Cmd) {
		cmd.Command("search s", "Search for devices", func(cmd *cli.Cmd) {

			cmd.Command("setting", "Search for devices by exact setting value", func(cmd *cli.Cmd) {
				var (
					keyArg   = cmd.StringArg("KEY", "", "Setting name")
					valueArg = cmd.StringArg("VALUE", "", "Setting Value")
				)
				cmd.Spec = "KEY VALUE"

				cmd.Action = func() {
					fmt.Println(API.Devices().FindBySetting(*keyArg, *valueArg))
				}
			})

			cmd.Command("tag", "Search for devices by exact tag value", func(cmd *cli.Cmd) {
				var (
					keyArg   = cmd.StringArg("KEY", "", "Tag name")
					valueArg = cmd.StringArg("VALUE", "", "Tag Value")
				)
				cmd.Spec = "KEY VALUE"

				cmd.Action = func() {
					fmt.Println(API.Devices().FindByTag(*keyArg, *valueArg))
				}
			})

			cmd.Command("hostname", "Search for devices by exact hostname", func(cmd *cli.Cmd) {
				var (
					hostnameArg = cmd.StringArg("HOSTNAME", "", "hostname")
				)
				cmd.Spec = "HOSTNAME"

				cmd.Action = func() {
					fmt.Println(API.Devices().FindByField("hostname", *hostnameArg))
				}
			})
		})
	},
	)

	App.Command("device", "Perform actions against a single device", func(cmd *cli.Cmd) {
		idArg := cmd.StringArg(
			"DEVICE",
			"",
			"UUID or serial number of the device. Short UUIDs are *not* accepted",
		)

		cmd.Spec = "DEVICE"

		cmd.Command("get", "Get information about a single device", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Devices().Get(*idArg)) }
		})

		cmd.Command("validations", "Get the most recent validation results for a single device", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Devices().ValidationState(*idArg)) }
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

			cmd.Command("get", "Get a particular device setting", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Devices().Setting(*idArg, *keyArg))
				}
			})

			cmd.Command("set", "Set a particular device setting", func(cmd *cli.Cmd) {
				valueArg := cmd.StringArg("VALUE", "", "Value of the setting")
				cmd.Spec = "VALUE"

				cmd.Action = func() {
					API.Devices().SetSetting(*idArg, *keyArg, *valueArg)
					fmt.Println(API.Devices().Settings(*idArg))
				}
			})

			cmd.Command("delete rm", "Delete a particular device setting", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					API.Devices().DeleteSetting(*idArg, *keyArg)
					fmt.Println(API.Devices().Settings(*idArg))
				}
			})
		})

		cmd.Command("tags", "See all tags for a device", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Devices().Tags(*idArg)) }
		})
		cmd.Command("tag", "See a single tag for a device", func(cmd *cli.Cmd) {
			nameArg := cmd.StringArg("NAME", "", "Name of the tag")

			cmd.Spec = "NAME"

			cmd.Action = func() {
				fmt.Println(API.Devices().Tag(*idArg, *nameArg))
			}

			cmd.Command("get", "Get a particular device tag", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Devices().Tag(*idArg, *nameArg))
				}
			})

			cmd.Command("set", "Set a particular device tag", func(cmd *cli.Cmd) {
				valueArg := cmd.StringArg("VALUE", "", "Value of the tag")
				cmd.Spec = "VALUE"

				cmd.Action = func() {
					API.Devices().SetTag(*idArg, *nameArg, *valueArg)
					fmt.Println(API.Devices().Tags(*idArg))
				}
			})

			cmd.Command("delete rm", "Delete a particular device tag", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					API.Devices().DeleteTag(*idArg, *nameArg)
					fmt.Println(API.Devices().Tags(*idArg))
				}
			})
		})

		cmd.Command("interface", "Information about a single interface", func(cmd *cli.Cmd) {
			nameArg := cmd.StringArg("NAME", "", "Name of the interface")
			cmd.Spec = "NAME"
			cmd.Action = func() { fmt.Println(API.Devices().GetInterface(*idArg, *nameArg)) }
		})

		cmd.Command("preflight", "Data that is only accurate inside preflight", func(cmd *cli.Cmd) {
			cmd.Before = func() {
				if API.Devices().GetPhase(*idArg) != "integration" {
					os.Stderr.WriteString("Warning: This device is no longer in the 'integration' phase. This data is likely to be inaccurate\n")
				}
			}

			cmd.Command("location", "The location of a device in preflight", func(cmd *cli.Cmd) {
				cmd.Action = func() { fmt.Println(API.Devices().GetLocation(*idArg)) }
			})

			cmd.Command("ipmi", "IPMI address for a device in preflight", func(cmd *cli.Cmd) {
				cmd.Action = func() { fmt.Println(API.Devices().GetIPMI(*idArg)) }
			})
		})

		cmd.Command("phase", "Actions on the lifecycle phase of the device", func(cmd *cli.Cmd) {
			cmd.Command("get", "Get the phase of the device", func(cmd *cli.Cmd) {
				cmd.Action = func() { fmt.Println(API.Devices().GetPhase(*idArg)) }
			})

			cmd.Command("set", "Set the phase of the device [one of: "+prettyPhasesList()+"]", func(cmd *cli.Cmd) {
				phaseArg := cmd.StringArg("PHASE", "", "Name of the phase [one of: "+prettyPhasesList()+"]")
				cmd.Spec = "PHASE"
				cmd.Action = func() {
					if !okPhase(*phaseArg) {
						panic("Phase must be one of: " + prettyPhasesList())
					}

					fmt.Println(API.Devices().SetPhase(*idArg, *phaseArg))
				}
			})
		})

		cmd.Command("validations", "Information about the latest validation runs", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Devices().ValidationState(*idArg)) }
		})

		cmd.Command("report", "Get the most recently recorded report for this device", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				d := API.Devices().Get(*idArg)
				if d.LatestReport == nil {
					fmt.Println("{}")
					return
				}
				API.PrintJSON(d.LatestReport)
			}
		})
	})
}
