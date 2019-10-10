// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

const validationStatesWithResultsTemplate = `{{ range . }}
- ID: {{ .ID }}
  Created: {{ TimeStr .Created }}
  Completed: {{ TimeStr .Completed }}
  Status: {{ .Status }}
  Validation Plan: {{ .ValidationPlan.Name }}{{ if len .Results }}

  Results:
{{ .Results }}
{{ end }}{{ end }}
`

const deviceTemplate = `
ID: {{ .ID }}
Serial: {{ .Serial }}
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

Hardware:
  Name: {{ .HardwareProduct.Name }}
  Legacy Name: {{ .HardwareProduct.LegacyProductName }}
  Alias: {{ .HardwareProduct.Alias }}
  Prefix: {{ .HardwareProduct.Prefix }}
  SKU: {{ .HardwareProduct.SKU }}
  Generation Name: {{ .HardwareProduct.GenerationName }}

Location: {{- if ne .Phase "integration" }} ** Device has left integration. This data is historic and likely not accurate. **{{ end }}
  AZ:  {{ .Location.Room.AZ }}
  Datacenter:
    ID: {{ .Location.Datacenter.ID }}
    Vendor:   {{ .Location.Datacenter.Vendor }} / {{ .Location.Datacenter.VendorName }}
    Region:   {{ .Location.Datacenter.Region }}
    Location: {{ .Location.Datacenter.Location }}

  Room:
    ID: {{ .Location.Room.ID }}
    Alias: {{ .Location.Room.Alias }}
    Vendor Name: {{ .Location.Room.VendorName }}

  Rack:
    ID:    {{ .Location.Rack.ID }}
    Name:  {{ .Location.Rack.Name }}{{ if ne .RackRole.Name "" }}
    Role:  {{ .RackRole.Name }}{{ end }}
    Phase: {{ .Location.Rack.Phase }}
    RU:    {{ .Location.RackUnitStart }}


Network Interfaces: {{ range .Nics }}
  - {{ .InterfaceName }} - {{ .Mac }}
    Type: {{ .InterfaceType }}
    Vendor: {{ .InterfaceVendor }}{{ if ne .PeerMac "" }}
    Peer: {{ .PeerMac }}{{ end }}{{ if ne .PeerSwitch "" }} - {{ .PeerSwitch }}{{ end }}
{{ end }}
Disks:{{range $name, $slots := .Enclosures}}
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

Validations:
{{ .Validations }}
`

const workspaceRelayTemplate = `
ID: {{ .ID }}
Name: {{ .Alias }}
Version: {{ .Version }}
Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}

Last Seen: {{ TimeStr .LastSeen }}

IP Address: {{ .IpAddr }}
SSH Port: {{ .SshPort }}

Location:
  AZ: {{ .Location.AZ }}
  Rack Name: {{ .Location.RackName }}
  Rack Unit: {{ .Location.RackUnitStart }}
  Rack ID: {{ .Location.RackID }}
`

const relayTemplate = `
ID: {{ .ID }}
Serial Number: {{ .SerialNumber }}
Name: {{ .Name }}
Version: {{ .Version }}
Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}

IP Address: {{ .IpAddr }}
SSH Port: {{ .SshPort }}
`

const rackSummaryTemplate = `
Name: {{ .AZ }} {{ .Name }}
ID: {{ .ID }}
Size: {{ .RackSize }}
Phase: {{ .Phase }}
Device Progress: {{ range .Statuses }}
  * {{ .Status }}: {{ .Count -}}
{{end}}
`

const workspaceTemplate = `
Name: {{ .Name }}
ID: {{ .ID }}
Description: {{ .Description }}
Your Role: {{ .Role }}
Your Role Was Derived From: {{ if eq "" $.Via }}[Direct Assignment]{{ else }}{{ .Via }}{{ end }}
`

const detailedUserTemplate = `
ID: {{ .ID }}
Name: {{ .Name }}
Email: {{ .Email }}
System Admin: {{ if $.IsAdmin }}Yes{{ else }}No{{ end }}

Created: {{ TimeStr .Created }}
Last Login: {{ if $.LastLogin.IsZero }}Never/Unknown{{ else }}{{ TimeStr .LastLogin }}{{ end }}


Workspaces:
{{ .Workspaces }}

Organizations:
{{ .Organizations }}
`

const organizationTemplate = `
Name: {{ .Name }}
ID: {{ .ID }}
Description: {{ .Description }}
`

const datacenterTemplate = `
ID: {{ .ID }}
Vendor: {{ .Vendor }}
Vendor Name: {{ .VendorName }}
Region: {{ .Region }}
Location: {{ .Location }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

const roomTemplate = `
Alias: {{ .Alias }}
AZ: {{ .AZ }}
Vendor Name: {{ .VendorName }}
Datacenter ID: {{ .DatacenterID }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

const rackRoleTemplate = `
Name: {{ .Name }}
Rack Size: {{ .RackSize }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

const rackTemplate = `
ID: {{ .ID }}
Name: {{ .Name }}
Serial Number: {{ .SerialNumber }}
Asset Tag: {{ .AssetTag }}
Phase: {{ .Phase }}
Role: {{ .Role.Name }}
Room: {{ .Room.Alias }}

Created: {{ TimeStr .Created }}
Updated: {{ TimeStr .Updated }}
`

const deviceLocationTemplate = `
Datacenter:
  ID: {{ .Datacenter.ID }}
  Vendor: {{ .Datacenter.Vendor }}
  Vendor Name: {{ .Datacenter.VendorName }}
  Region: {{ .Datacenter.Region }}
  Location: {{ .Datacenter.Location }}

Room:
  ID: {{ .Room.ID }}
  Alias: {{ .Room.Alias }}
  AZ: {{ .Room.AZ }}
  Vendor Name: {{ .Room.VendorName }}

Rack:
  ID: {{ .Rack.ID }}
  Name: {{ .Rack.Name }}
  Serial Number: {{ .Rack.SerialNumber }}
  Asset Tag: {{ .Rack.AssetTag }}
  Phase: {{ .Rack.Phase }}
  Role: {{ .Rack.Role.Name }}

Rack Unit Start: {{ .RackUnitStart }}
`

const deviceNicTemplate = `
Name: {{ .InterfaceName }}
Vendor: {{ .InterfaceVendor }}
Driver: {{ .InterfaceDriver }}

IP Address: {{ .IpAddress }}
MAC: {{ .MAC }}
MTU: {{ .MTU }}
State: {{ .State }}

Device ID: {{ .DeviceID }}
`

const buildTemplate = `
Name: {{ .Name }}
Description: {{ .Description }}
Created: {{ .Created }}
Started: {{ .Started }}
Completed: {{ .Completed }}
Marked Complete By: {{ .CompletedUser.Email }}
`
