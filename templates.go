// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

const relayTemplate = `
ID: {{ .ID }}
Alias: {{ .Alias }}
Version: {{ .Version }}
Created: {{ .Created.Local }}
Last Seen: {{ .LastSeen }}

IP Address: {{ .IpAddr }}
SSH Port: {{ .SshPort }}

Location: 
  AZ: {{ .Location.AZ }}
  Rack Name: {{ .Location.RackName }}
  Rack Unit: {{ .Location.RackUnitStart }}
  Rack ID: {{ .Location.RackID }}
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

Created: {{ .Created.Local }}
Last Login: {{ if $.LastLogin.IsZero }}Never/Unknown{{ else }}{{ .LastLogin.Local }}{{ end }}


Workspaces:
{{ .Workspaces }}
`

const datacenterTemplate = `
ID: {{ .ID }}
Vendor: {{ .Vendor }}
Vendor Name: {{ .VendorName }}
Region: {{ .Region }}
Location: {{ .Location }}

Created: {{ .Created.Local }}
Updated: {{ .Updated.Local }}
`

const roomTemplate = `
Alias: {{ .Alias }}
AZ: {{ .AZ }}
Vendor Name: {{ .VendorName }}
Datacenter ID: {{ .DatacenterID }}

Created: {{ .Created.Local }}
Updated: {{ .Updated.Local }}
`

const rackRoleTemplate = `
Name: {{ .Name }}
Rack Size: {{ .RackSize }}

Created: {{ .Created.Local }}
Updated: {{ .Updated.Local }}
`

const rackTemplate = `
ID: {{ .ID }}
Name: {{ .Name }}
Serial Number: {{ .SerialNumber }}
Asset Tag: {{ .AssetTag }}
Phase: {{ .Phase }}
Role: {{ .Role.Name }}
Room: {{ .Room.Alias }}

Created: {{ .Created.Local }}
Updated: {{ .Updated.Local }}
`
