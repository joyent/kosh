package types

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/joyent/kosh/tables"
	"github.com/joyent/kosh/template"
)

func (bl Builds) Len() int           { return len(bl) }
func (bl Builds) Swap(i, j int)      { bl[i], bl[j] = bl[j], bl[i] }
func (bl Builds) Less(i, j int) bool { return bl[i].Name < bl[j].Name }
func (bl Builds) Headers() []string {
	return []string{
		"Name",
		"Description",
		"Started",
		"Completed",
	}
}

func (bl Builds) ForEach(do func([]string)) {
	for _, b := range bl {
		do([]string{
			string(b.Name),
			b.Description,
			template.TimeStr(b.Started),
			template.TimeStr(b.Completed),
		})
	}
}

const buildTemplate = `
Build {{ .Name }}
=================

{{ .Description }}

Admins
------

{{ range .Admins }}
* {{ .Name }} - {{ .Email }}
{{ end }}

Links
-----
{{ range .Links }}
* {{ . }}
{{ end }}

---
* Created: {{ TimeStr .Created }}
* Started: {{ TimeStr .Started }}
* Completed: {{ TimeStr .Completed }} by {{ .CompletedUser.Name }}({{ .CompletedUser.Email }})
`

func (b Build) Template() string { return buildTemplate }

func (bu BuildUsers) Len() int           { return len(bu) }
func (bu BuildUsers) Swap(i, j int)      { bu[i], bu[j] = bu[j], bu[i] }
func (bu BuildUsers) Less(i, j int) bool { return bu[i].Name < bu[j].Name }
func (bu BuildUsers) Headers() []string {
	return []string{
		"ID",
		"Name",
		"Email",
		"Role",
	}
}

func (bu BuildUsers) ForEach(do func([]string)) {
	for _, u := range bu {
		do([]string{
			template.CutUUID(u.ID.String()),
			u.Name,
			string(u.Email),
			string(u.Role),
		})
	}
}

func (bo BuildOrganizations) Len() int           { return len(bo) }
func (bo BuildOrganizations) Swap(i, j int)      { bo[i], bo[j] = bo[j], bo[i] }
func (bo BuildOrganizations) Less(i, j int) bool { return bo[i].Name < bo[j].Name }
func (bo BuildOrganizations) Headers() []string {
	return []string{
		"ID",
		"Name",
		"Description",
		"Role",
	}
}

func (bo BuildOrganizations) ForEach(do func([]string)) {
	for _, o := range bo {
		do([]string{
			template.CutUUID(o.ID.String()),
			o.Name,
			o.Description,
			string(o.Role),
		})
	}
}

func (dl Datacenters) Len() int           { return len(dl) }
func (dl Datacenters) Swap(i, j int)      { dl[i], dl[j] = dl[j], dl[i] }
func (dl Datacenters) Less(i, j int) bool { return dl[i].VendorName < dl[j].VendorName }
func (dl Datacenters) Headers() []string {
	return []string{
		"ID",
		"Vendor",
		"Vendor Name",
		"Region",
		"Location",
	}
}

func (dl Datacenters) ForEach(do func([]string)) {
	for _, d := range dl {
		do([]string{
			template.CutUUID(d.ID.String()),
			d.Vendor,
			d.VendorName,
			d.Region,
			d.Location,
		})
	}
}

const datacenterTemplate = `
Datacenter
==========

ID: {{ .ID }}
Vendor: {{ .Vendor }}
Vendor Name: {{ .VendorName }}
Region: {{ .Region }}
Location: {{ .Location }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

func (d Datacenter) Template() string { return datacenterTemplate }

func (ds DeviceSettings) String() string {
	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)

	var keys []string
	for key := range ds {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := ds[key]
		table.Append([]string{key, string(value)})
	}

	table.Render()
	return tableString.String()
}

const detailedDeviceTemplate = `
Device {{ .SerialNumber }}
==========================

ID: {{ .ID }}
Serial: {{ .SerialNumber }}
Asset Tag: {{ .AssetTag }}
Hostname: {{ .Hostname }}
System UUID: {{ .SystemUUID }}

Phase: {{ .Phase }}
Health: {{ .Health }}
Validated: {{ if not $.Validated.IsZero }}{{ .Validated.Local }}{{ end }}

Created:   {{ TimeStr .Created }}
Updated:   {{ TimeStr .Updated }}
Last Seen: {{ TimeStr .LastSeen }}{{ if .Links }}

Links: {{ range .Links }}
  - {{ $ }}
{{ end }}{{ end }}

Location: {{- if ne .Phase "integration" }} ** Device has left integration. This data is historic and likely not accurate. **{{ end }}
  AZ:  {{ .Location.Az }}
  Datacenter:
    Datacenter: {{ .Location.DatacenterRoom }}
    Rack:   {{ .Location.Rack }}
    RU:   {{ .Location.RackUnitStart }}


Network Interfaces: {{ range .Nics }}
  - {{ .InterfaceName }} - {{ .Mac }}
    Type: {{ .InterfaceType }}
    Vendor: {{ .InterfaceVendor }}{{ if ne .PeerMac "" }}
    Peer: {{ .PeerMac }}{{ end }}{{ if ne .PeerSwitch "" }} - {{ .PeerSwitch }}{{ end }}
{{ end }}

Disks:{{range $name, $slots := .Disks}}
  Enclosure: {{ $name }}{{ range $slots }}
    Slot: {{ .Slot }}
        SN:     {{ .SerialNumber }}
        Type:   {{ .DriveType }}
        Vendor: {{ .Vendor }}
        Model:  {{ .Model }}
        Size:   {{ .Size }}
        Health: {{ .Health }}
        Firmware: {{ .Firmware }}
        Transport: {{ .Transport }}
{{ end }}{{ end }}
`

func (d DetailedDevice) Template() string { return detailedDeviceTemplate }

const deviceTemplate = `
Device {{ .SerialNumber }}
==========================

ID: {{ .ID }}
Serial: {{ .SerialNumber }}
Asset Tag: {{ .AssetTag }}
Hostname: {{ .Hostname }}
System UUID: {{ .SystemUUID }}

Phase: {{ .Phase }}
Health: {{ .Health }}
Validated: {{ if not $.Validated.IsZero }}{{ .Validated.Local }}{{ end }}

Created:   {{ TimeStr .Created }}
Updated:   {{ TimeStr .Updated }}
Last Seen: {{ TimeStr .LastSeen }}{{ if .Links }}

Rack:
  ID:    {{ CutUUID .RackID }}
  Name:  {{ .RackName }}
  RU:    {{ .RackUnitStart }}

Links: {{ range .Links }}
  - {{ $ }}
{{ end }}{{ end }}

`

func (d Device) Template() string { return deviceTemplate }

func (d Devices) Len() int           { return len(d) }
func (d Devices) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Devices) Less(i, j int) bool { return d[i].SerialNumber < d[j].SerialNumber }
func (d Devices) Headers() []string {
	return []string{
		"Serial",
		"Hostname",
		"Asset Tag",
		"Hardware",
		"Phase",
		"Updated",
		"Validated",
	}
}

func (d Devices) ForEach(do func([]string)) {
	for _, device := range d {
		do([]string{
			string(device.SerialNumber),
			device.Hostname,
			string(device.AssetTag),
			device.HardwareProductID.String(),
			string(device.Phase),
			template.TimeStr(device.Updated),
			template.TimeStr(device.Validated),
		})
	}
}

const deviceReportTemplate = ``

func (d DeviceReport) Template() string { return deviceReportTemplate }

const hardwareProductTemplate = `
Hardware Product {{ .Name }}
============================

ID: {{ .ID }}
Name: {{ .Name }}
SKU: {{ .SKU }}

Alias: {{ .Alias }}
GenerationName: {{ .GenerationName }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

func (hp HardwareProduct) Template() string { return hardwareProductTemplate }

// TODO sort interface, tabulable interface
func (h HardwareProducts) String() string {
	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)

	table.SetHeader([]string{
		"ID",
		"SKU",
		"Name",
		"Alias",
		"GenerationName",
		"Created",
		"Updated",
	})

	for _, hp := range h {
		table.Append([]string{
			template.CutUUID(hp.ID.String()),
			string(hp.SKU),
			string(hp.Name),
			string(hp.Alias),
			hp.GenerationName,
			hp.Created.String(),
			hp.Updated.String(),
		})
	}
	table.Render()
	return tableString.String()
}

const hardwareVendorTemplate = `
Hardware Vendor {{ .Name }}
===========================

Name: {{ .Name }}
ID: {{ .ID }}
Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

func (h HardwareVendor) Template() string { return hardwareVendorTemplate }

func (h HardwareVendors) Len() int           { return len(h) }
func (h HardwareVendors) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h HardwareVendors) Less(i, j int) bool { return h[i].Name < h[j].Name }
func (h HardwareVendors) Headers() []string {
	return []string{
		"Name",
		"ID",
		"Created",
		"Updated",
	}
}

func (h HardwareVendors) ForEach(do func([]string)) {
	for _, v := range h {
		do([]string{
			string(v.Name),
			template.CutUUID(v.ID.String()),
			template.TimeStr(v.Created),
			template.TimeStr(v.Updated),
		})
	}
}

const organizationTemplate = `
Organization {{ .Name }}
========================
ID: {{ .ID }}
Description: {{ .Description }}
`

func (o Organization) Template() string { return organizationTemplate }

func (o Organizations) Len() int           { return len(o) }
func (o Organizations) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Organizations) Less(i, j int) bool { return o[i].Name < o[j].Name }
func (o Organizations) Headers() []string {
	return []string{
		"Name",
		"Role",
		"Description",
	}
}

func (o Organizations) ForEach(do func([]string)) {
	for _, org := range o {
		do([]string{
			string(org.Name),
			string(org.Role),
			org.Description,
		})
	}
}
func (o Organizations) String() { tables.Render(o) }

const rackTemplate = `
Rack {{ .Name }}
================

ID: {{ .ID }}
Name: {{ .Name }}
Serial Number: {{ .SerialNumber }}
Asset Tag: {{ .AssetTag }}
Phase: {{ .Phase }}
Role: {{ .RackRoleName }}
Room: {{ .DatacenterRoomAlias }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

func (r Rack) Template() string { return rackTemplate }

func (r Racks) Len() int           { return len(r) }
func (r Racks) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Racks) Less(i, j int) bool { return r[i].SerialNumber > r[j].SerialNumber }
func (rl Racks) Headers() []string {
	return []string{
		"ID",
		"Name",
		"Room",
		"Role",
		"Serial Number",
		"Asset Tag",
		"Phase",
		"Created",
		"Updated",
	}
}

func (rl Racks) ForEach(do func([]string)) {
	for _, r := range rl {
		do([]string{
			r.ID.String(),
			string(r.Name),
			template.CutUUID(r.DatacenterRoomID.String()),
			string(r.RackRoleName),
			string(r.SerialNumber),
			string(r.AssetTag),
			string(r.Phase),
			template.TimeStr(r.Created),
			template.TimeStr(r.Updated),
		})
	}
}

func (r RackLayouts) Len() int           { return len(r) }
func (r RackLayouts) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RackLayouts) Less(i, j int) bool { return r[i].RackUnitStart > r[j].RackUnitStart }

func (rl RackLayouts) Headers() []string {
	return []string{
		"Rack Unit Start",
		"Rack Unit Size",
		"ID",
		"Hardware Product",
		"Created",
		"Updated",
	}
}

func (rl RackLayouts) ForEach(do func([]string)) {
	for _, r := range rl {
		do([]string{
			strconv.Itoa(int(r.RackUnitStart)),
			strconv.Itoa(int(r.RackUnitSize)),
			template.CutUUID(r.ID.String()),
			r.HardwareProductID.String(),
			template.TimeStr(r.Created),
			template.TimeStr(r.Updated),
		})
	}
}

func (ra RackAssignments) Len() int           { return len(ra) }
func (ra RackAssignments) Swap(i, j int)      { ra[i], ra[j] = ra[j], ra[i] }
func (ra RackAssignments) Less(i, j int) bool { return ra[i].RackUnitStart > ra[j].RackUnitStart }

func (ra RackAssignments) Headers() []string {
	return []string{
		"Device Serial",
		"Device Asset Tag",
		"Hardware Product",
		"Rack Unit Start",
		"Rack Unit Size",
	}
}

func (ra RackAssignments) ForEach(do func([]string)) {
	for _, r := range ra {
		do([]string{
			string(r.DeviceSerialNumber),
			string(r.DeviceAssetTag),
			string(r.HardwareProductName),
			strconv.Itoa(int(r.RackUnitStart)),
			strconv.Itoa(int(r.RackUnitSize)),
		})
	}
}

func (r RackRoles) Len() int           { return len(r) }
func (r RackRoles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RackRoles) Less(i, j int) bool { return r[i].Name < r[j].Name }
func (rl RackRoles) String() string {
	sort.Sort(rl)

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)

	table.SetHeader([]string{
		"Name",
		"RackSize",
		"Created",
		"Updated",
	})

	for _, r := range rl {
		table.Append([]string{
			string(r.Name),
			strconv.Itoa(int(r.RackSize)),
			template.TimeStr(r.Created),
			template.TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()
}

const rackRoleTemplate = `
Rack Role {{ .Name }}
=====================

Name: {{ .Name }}
Rack Size: {{ .RackSize }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

func (r RackRole) Template() string { return rackRoleTemplate }

const relayTemplate = `
Relay {{ .Name }}
=================

ID: {{ .ID }}
Serial Number: {{ .SerialNumber }}
Name: {{ .Name }}
Version: {{ .Version }}
Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}

IP Address: {{ .IpAddr }}
SSH Port: {{ .SshPort }}
`

func (r Relay) Template() string { return relayTemplate }

func (r Relays) Len() int           { return len(r) }
func (r Relays) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Relays) Less(i, j int) bool { return r[i].Updated.Before(r[j].Updated) }

func (rl Relays) String() string {
	sort.Sort(rl)

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)

	table.SetHeader([]string{
		"Serial",
		"Name",
		"Version",
		"IP",
		"Updated",
	})

	for _, r := range rl {
		table.Append([]string{
			string(r.SerialNumber),
			string(r.Name),
			string(r.Version),
			string(r.Ipaddr),
			template.TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()
}

const roomTemplate = `
Room {{ .Alias }}
=================

ID: {{ .ID }}
Alias: {{ .Alias }}
AZ: {{ .AZ }}
Vendor Name: {{ .VendorName }}
Datacenter ID: {{ .DatacenterID }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

func (r DatacenterRoomDetailed) Template() string { return roomTemplate }

func (dr DatacenterRoomsDetailed) Len() int           { return len(dr) }
func (dr DatacenterRoomsDetailed) Swap(i, j int)      { dr[i], dr[j] = dr[j], dr[i] }
func (dr DatacenterRoomsDetailed) Less(i, j int) bool { return dr[i].Alias < dr[j].Alias }

func (dr DatacenterRoomsDetailed) Headers() []string {
	return []string{
		"ID",
		"Alias",
		"AZ",
		"Vendor Name",
		"Datacenter ID",
		"Created",
		"Updated",
	}
}

func (dr DatacenterRoomsDetailed) ForEach(do func([]string)) {
	for _, r := range dr {
		do([]string{
			template.CutUUID(r.ID.String()),
			string(r.Alias),
			string(r.AZ),
			string(r.VendorName),
			template.CutUUID(r.DatacenterID.String()),
			template.TimeStr(r.Created),
			template.TimeStr(r.Updated),
		})
	}
}

func (u UserSettings) Headers() []string {
	return []string{
		"Key",
		"Value",
	}
}

func (u UserSettings) ForEach(do func([]string)) {
	keys := make([]string, 0)
	for setting := range u {
		keys = append(keys, setting)
	}
	sort.Strings(keys)

	for _, key := range keys {
		do([]string{
			key,
			fmt.Sprintf("%v", u[key]),
		})
	}
}

const detailedUserTemplate = `
User {{ .Name }}
================

* ID: {{ .ID }}
* Email: {{ .Email }}
* System Admin: {{ if $.IsAdmin }}Yes{{ else }}No{{ end }}

Created: {{ TimeStr .Created }}

Last Login: {{ if $.LastLogin.IsZero }}Never/Unknown{{ else }}{{ TimeStr .LastLogin }}{{ end }}

Organizations
-------------

{{ Table .Organizations }}
`

func (u UserDetailed) Template() string { return detailedUserTemplate }

const validationPlanTemplate = `
Validation Plan {{ .Name }}
===========================

ID: {{ .ID }}
Name: {{ .Name }}
Description: {{ .Description }}
Created: {{ .Created }}
`

func (v ValidationPlan) Template() string { return validationPlanTemplate }

func (v ValidationPlans) Len() int           { return len(v) }
func (v ValidationPlans) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ValidationPlans) Less(i, j int) bool { return v[i].Name < v[j].Name }

func (v ValidationPlans) String() string {
	sort.Sort(v)

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
	table.SetRowLine(true)

	table.SetHeader([]string{
		"ID",
		"Name",
		"Description",
		"Created",
	})

	for _, p := range v {
		table.Append([]string{
			template.CutUUID(p.ID.String()),
			string(p.Name),
			p.Description,
			p.Created.String(),
		})
	}

	table.Render()
	return tableString.String()
}

const validationStateWithResultsTemplate = `
Validation State
================

ID: {{ .ID }}
Device: {{ CutUUID .DeviceID.String }}
Hardware Product: {{ CutUUID .HardwareProductID.String }}
Created: {{ TimeStr .Created }}
Status: {{ .Status }}

Results:
{{ .Results }}
`

func (v ValidationStateWithResults) Template() string { return validationStateWithResultsTemplate }

func (v ValidationResults) Len() int           { return len(v) }
func (v ValidationResults) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ValidationResults) Less(i, j int) bool { return v[i].Category < v[j].Category }

func (v ValidationResults) String() string {
	sort.Sort(v)

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
	table.SetRowLine(true)

	table.SetHeader([]string{
		"Status",
		"Category",
		"Component",
		"Message",
	})

	for _, r := range v {
		table.Append([]string{
			string(r.Status),
			r.Category,
			r.Component,
			r.Message,
		})
	}

	table.Render()
	return tableString.String()
}

const deviceNicTemplate = `
Nic {{ .IfaceName }}
====================

Name: {{ .IfaceName }}
Vendor: {{ .IfaceVendor }}
Type: {{ .IfaceType }}

IP Address: {{ .Ipaddr }}
MAC: {{ .MAC }}
MTU: {{ .MTU }}
State: {{ .State }}

Device ID: {{ .DeviceID }}
`

func (dn DeviceNic) Template() string { return deviceNicTemplate }

const deviceLocationTemplate = `
Location
========

Rack {{ .Rack }}
Rack Unit Start: {{ .RackUnitStart }}
DatacenterRoom: {{ .DatacenterRoom }}
AZ: {{ .Az }}
`

func (dl DeviceLocation) Template() string { return deviceLocationTemplate }

func (ul UsersTerse) Len() int           { return len(ul) }
func (ul UsersTerse) Swap(i, j int)      { ul[i], ul[j] = ul[j], ul[i] }
func (ul UsersTerse) Less(i, j int) bool { return ul[i].Name < ul[j].Name }
func (ul UsersTerse) Headers() []string {
	return []string{
		"Name",
		"Email",
	}
}

func (ul UsersTerse) ForEach(do func([]string)) {
	for _, u := range ul {
		do([]string{
			string(u.Name),
			string(u.Email),
		})
	}
}

const userTokenTemplate = `
Token {{ .Name }}
<<<<<<< HEAD
=================
=======
>>>>>>> 3c81f5d (render output for User Tokens)

* Created: {{ TimeStr .Created }}

* Last IP: {{ .LastIpaddr }}
* Last Used: {{ TimeStr .LastUsed }}

* Expires: {{ TimeStr .Expires }}
`

func (ut UserToken) Template() string { return userTokenTemplate }

const newUserTokenTemplate = `
Token {{ .Name }}
=================

* Created: {{ TimeStr .Created }}

* Last IP: {{ .LastIpaddr }}
* Last Used: {{ TimeStr .LastUsed }}

* Expires: {{ TimeStr .Expires }}

Token
-----

THIS TOKEN CANNOT BE RECOVERED FROM THE SERVER.
THIS IS THE ONLY TIME IT WILL BE PRINTED, PLEASE RECORD IT NOW

{{ .Token }}

`

func (ut NewUserTokenResponse) Template() string { return newUserTokenTemplate }

func (ul UserTokens) Len() int           { return len(ul) }
func (ul UserTokens) Swap(i, j int)      { ul[i], ul[j] = ul[j], ul[i] }
func (ul UserTokens) Less(i, j int) bool { return ul[i].Name < ul[j].Name }
func (ul UserTokens) Headers() []string {
	return []string{
		"Name",
		"Created",
		"Last IP",
		"Last Used",
		"Expires",
	}
}

func (ul UserTokens) ForEach(do func([]string)) {
	for _, u := range ul {
		do([]string{
			string(u.Name),
			template.TimeStr(u.Created),
			string(u.LastIpaddr),
			template.TimeStr(u.LastUsed),
			template.TimeStr(u.Expires),
		})
	}
}
