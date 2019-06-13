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
