// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Workspaces struct {
	*Conch
}

func (c *Conch) Workspaces() *Workspaces {
	return &Workspaces{c}
}

/****/

var WorkspaceRoleList = []string{"admin", "rw", "ro"}

func prettyWorkspaceRoleList() string {
	return strings.Join(WorkspaceRoleList, ", ")
}

func okWorkspaceRole(role string) bool {
	for _, b := range WorkspaceRoleList {
		if role == b {
			return true
		}
	}
	return false
}

/****/

type WorkspaceAndRole struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ParentID    uuid.UUID `json:"parent_id,omitempty"`
	Role        string    `json:"role"`
	RoleVia     uuid.UUID `json:"role_via"`

	// These are for user friendly variants of those UUIDs
	Parent string `json:"parent"`
	Via    string `json:"via"`
}

func (w WorkspaceAndRole) String() string {
	if API.JsonOnly {
		return API.AsJSON(w)
	}

	t, err := template.New("w").Parse(workspaceTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, w); err != nil {
		panic(err)
	}

	return buf.String()
}

type WorkspaceAndRoles []WorkspaceAndRole

func (w WorkspaceAndRoles) Len() int {
	return len(w)
}

func (w WorkspaceAndRoles) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WorkspaceAndRoles) Less(i, j int) bool {
	return w[i].Name < w[j].Name
}

func (w WorkspaceAndRoles) String() string {
	sort.Sort(w)
	if API.JsonOnly {
		return API.AsJSON(w)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Name",
		"Role",
		"Description",
		"Role Via",
		"Parent",
	})

	for _, ws := range w {
		table.Append([]string{
			ws.Name,
			ws.Role,
			ws.Description,
			ws.Via,
			ws.Parent,
		})
	}

	table.Render()
	return tableString.String()
}

/****/

func (w *Workspaces) GetAll() WorkspaceAndRoles {
	list := make(WorkspaceAndRoles, 0)

	res := w.Do(w.Sling().Get("/workspace"))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}

	ret := make(WorkspaceAndRoles, 0)

	cache := make(map[uuid.UUID]string)

	for _, ws := range list {
		if (ws.ParentID != uuid.UUID{}) {
			if _, ok := cache[ws.ParentID]; !ok {
				cache[ws.ParentID] = w.Get(ws.ParentID).Name
			}

			ws.Parent = cache[ws.ParentID]
		}

		if (ws.RoleVia != uuid.UUID{}) {
			if _, ok := cache[ws.RoleVia]; !ok {
				cache[ws.RoleVia] = w.Get(ws.RoleVia).Name
			}

			ws.Via = cache[ws.RoleVia]
		}

		ret = append(ret, ws)
	}

	return ret
}

func (w *Workspaces) Get(id uuid.UUID) (ws WorkspaceAndRole) {
	res := w.Do(w.Sling().Get("/workspace/" + url.PathEscape(id.String())))
	if ok := res.Parse(&ws); !ok {
		panic(res)
	}
	if (ws.ParentID != uuid.UUID{}) {
		ws.Parent = w.Get(ws.ParentID).Name
	}

	if (ws.RoleVia != uuid.UUID{}) {
		ws.Via = w.Get(ws.ParentID).Name
	}

	return ws
}

func (w *Workspaces) GetByName(name string) (ws WorkspaceAndRole) {
	res := w.Do(w.Sling().Get("/workspace/" + url.PathEscape(name)))
	if ok := res.Parse(&ws); !ok {
		panic(res)
	}

	if (ws.ParentID != uuid.UUID{}) {
		ws.Parent = w.Get(ws.ParentID).Name
	}

	if (ws.RoleVia != uuid.UUID{}) {
		ws.Via = w.Get(ws.ParentID).Name
	}

	return ws
}

func (w *Workspaces) Create(parent string, sub string, desc string) (ws WorkspaceAndRole) {
	uri := fmt.Sprintf(
		"/workspace/%s/child",
		url.PathEscape(parent),
	)

	payload := make(map[string]string)
	payload["name"] = sub
	if desc != "" {
		payload["description"] = desc
	}

	res := w.Do(
		w.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	if ok := res.Parse(&ws); !ok {
		panic(res)
	}

	if (ws.ParentID != uuid.UUID{}) {
		ws.Parent = w.Get(ws.ParentID).Name
	}

	if (ws.RoleVia != uuid.UUID{}) {
		ws.Via = w.Get(ws.ParentID).Name
	}

	return ws
}

/***/

type WorkspaceUser struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	RoleVia uuid.UUID `json:"role_via"`

	Via string `json:"via"`
}

type WorkspaceUsers []WorkspaceUser

func (w WorkspaceUsers) String() string {
	if API.JsonOnly {
		return API.AsJSON(w)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Name",
		"Email",
		"Role",
		"Role Via",
	})

	sort.Sort(w)

	for _, user := range w {
		table.Append([]string{
			user.Name,
			user.Email,
			user.Role,
			user.Via,
		})
	}

	table.Render()

	return tableString.String()
}

func (w WorkspaceUsers) Len() int {
	return len(w)
}

func (w WorkspaceUsers) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WorkspaceUsers) Less(i, j int) bool {
	return w[i].Name < w[j].Name
}

/***/

func (w *Workspaces) GetUsers(name string) WorkspaceUsers {
	list := make(WorkspaceUsers, 0)
	users := make(WorkspaceUsers, 0)

	url := fmt.Sprintf(
		"/workspace/%s/user",
		url.PathEscape(name),
	)

	res := w.Do(w.Sling().Get(url))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}

	for _, u := range list {
		if (u.RoleVia != uuid.UUID{}) {
			u.Via = w.Get(u.RoleVia).Name
		}
		users = append(users, u)
	}

	return users
}

func (w *Workspaces) AddOrModifyUser(workspace string, email string, role string, sendEmail bool) bool {
	payload := make(map[string]string)

	if email == "" {
		panic(errors.New("email address is required"))
	} else {
		payload["email"] = email
	}

	if role == "" {
		payload["role"] = "ro"
	} else {
		payload["role"] = role
	}

	params := make(map[string]int)
	if sendEmail {
		params["send_email"] = 1
	} else {
		params["send_email"] = 0
	}

	uri := fmt.Sprintf(
		"/workspace/%s/user",
		url.PathEscape(workspace),
	)

	p := struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}{payload["email"], payload["role"]}

	q := struct {
		SendEmail int `url:"send_mail"`
	}{params["send_email"]}

	res := w.Do(
		w.Sling().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(p).
			QueryStruct(q),
	)

	// NOTE: at time of writing, the only possible success response is a 204.
	// Everything else is a 400 or above. Here, that translates to a panic if
	// we didn't get a 20x. if we got a 20x, it's a 204. In the end, this will
	// return true or panic.
	return res.StatusCode() == 204
}

func (w *Workspaces) RemoveUser(workspace string, email string, sendEmail bool) bool {

	params := make(map[string]int)
	if sendEmail {
		params["send_email"] = 1
	} else {
		params["send_email"] = 0
	}
	q := struct {
		SendEmail int `url:"send_mail"`
	}{params["send_email"]}

	uri := fmt.Sprintf(
		"/workspace/%s/user/%s",
		url.PathEscape(workspace),
		url.PathEscape(email),
	)
	res := w.Do(w.Sling().Delete(uri).QueryStruct(q))

	// NOTE: at time of writing, the only possible success response is a 204.
	// Everything else is a 400 or above. Here, that translates to a panic if
	// we didn't get a 200. if we got a 200, it's a 204. In the end, this will
	// return true or panic.
	return res.StatusCode() == 204

}

/***/

func (w *Workspaces) GetDevices(name string, health string, validated *bool) DeviceList {
	var opts interface{}

	valid := 0
	if (validated != nil) && *validated {
		valid = 1
	}

	if (health != "") && (validated != nil) {
		opts = struct {
			Health    string `url:"health"`
			Validated int    `url:"validated"`
		}{url.PathEscape(health), valid}
	} else if health != "" {
		opts = struct {
			Health string `url:"health"`
		}{url.PathEscape(health)}
	} else if validated != nil {
		opts = struct {
			Validated int `url:"validated"`
		}{valid}
	}

	url := fmt.Sprintf(
		"/workspace/%s/device",
		url.PathEscape(name),
	)
	devices := make(DeviceList, 0)

	res := w.Do(w.Sling().Get(url).QueryStruct(opts))
	if ok := res.Parse(&devices); !ok {
		panic(res)
	}

	return devices
}

/***/

func (w *Workspaces) GetDirectChildren(name string) WorkspaceAndRoles {
	children := w.GetChildren(name)
	directs := make(WorkspaceAndRoles, 0)

	for _, child := range children {
		if child.Parent == name {
			directs = append(directs, child)
		}
	}
	return directs
}

func (w *Workspaces) GetChildren(name string) WorkspaceAndRoles {
	list := make(WorkspaceAndRoles, 0)

	uri := fmt.Sprintf("/workspace/%s/child", url.PathEscape(name))
	res := w.Do(w.Sling().Get(uri))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}

	ret := make(WorkspaceAndRoles, 0)

	cache := make(map[uuid.UUID]string)

	for _, ws := range list {
		if (ws.ParentID != uuid.UUID{}) {
			if _, ok := cache[ws.ParentID]; !ok {
				cache[ws.ParentID] = w.Get(ws.ParentID).Name
			}

			ws.Parent = cache[ws.ParentID]
		}

		if (ws.RoleVia != uuid.UUID{}) {
			if _, ok := cache[ws.RoleVia]; !ok {
				cache[ws.RoleVia] = w.Get(ws.RoleVia).Name
			}

			ws.Via = cache[ws.RoleVia]
		}

		ret = append(ret, ws)
	}

	return ret
}

/***/

// The DeviceProgress element is a map where the string is 'valid' or of the
// 'device_health' type containing 'error', 'fail', 'unknown', 'pass'
type WorkspaceRackSummary struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Phase          string         `json:"phase"`
	RoleName       string         `json:"role_name"`
	RackSize       int            `json:"rack_size"`
	DeviceProgress map[string]int `json:"device_progress"`
}

type WorkspaceRackSummaries map[string][]WorkspaceRackSummary

// The AZ gets lost in this conversion
func (summaries WorkspaceRackSummaries) Slice() []WorkspaceRackSummary {
	s := make([]WorkspaceRackSummary, 0)

	for _, values := range summaries {
		s = append(s, values...)
	}
	return s
}

func (summaries WorkspaceRackSummaries) String() string {
	if API.JsonOnly {
		return API.AsJSON(summaries)
	}

	var output string

	keys := make([]string, 0)
	for az := range summaries {
		keys = append(keys, az)
	}

	sort.Strings(keys)

	for _, az := range keys {
		for _, summary := range summaries[az] {
			type status struct {
				Status string
				Count  int
			}
			statusii := make([]*status, 0)

			statusStrs := make([]string, 0)
			for str := range summary.DeviceProgress {
				statusStrs = append(statusStrs, str)
			}
			sort.Strings(statusStrs)

			for _, statusStr := range statusStrs {
				statusii = append(statusii, &status{
					Status: statusStr,
					Count:  summary.DeviceProgress[statusStr],
				})
			}

			s := struct {
				WorkspaceRackSummary
				AZ       string
				Statuses []*status
			}{summary, az, statusii}

			t, err := template.New("r").Parse(rackSummaryTemplate)
			if err != nil {
				panic(err)
			}

			buf := new(bytes.Buffer)

			if err := t.Execute(buf, s); err != nil {
				panic(err)
			}

			output = output + buf.String()
		}
	}
	return output
}

/****/

func (w *Workspaces) GetRackSummaries(name string) WorkspaceRackSummaries {
	summaries := make(WorkspaceRackSummaries)

	url := fmt.Sprintf(
		"/workspace/%s/rack",
		url.PathEscape(name),
	)

	res := w.Do(w.Sling().Get(url))

	if ok := res.Parse(&summaries); !ok {
		panic(res)
	}

	return summaries
}

/***/

type WorkspaceRelay struct {
	ID         string    `json:"id"`
	Name       string    `json:"name,omitempty"`
	Version    string    `json:"version,omitempty"`
	IpAddr     string    `json:"ipaddr,omitempty"`
	SshPort    int       `json:"ssh_port,omitempty"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	LastSeen   time.Time `json:"last_seen"`
	NumDevices int       `json:"num_devices"`
	Location   struct {
		RackID        uuid.UUID `json:"rack_id"`
		RackName      string    `json:"rack_name"`
		RackUnitStart int       `json:"rack_unit_start"`
		RoleName      string    `json:"role_name"`
		AZ            string    `json:"az"`
	} `json:"location"`
}

func (w WorkspaceRelay) String() string {
	if API.JsonOnly {
		return API.AsJSON(w)
	}

	t, err := template.New("r").Parse(workspaceRelayTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, w); err != nil {
		panic(err)
	}

	return buf.String()
}

type WorkspaceRelays []WorkspaceRelay

func (w WorkspaceRelays) Len() int {
	return len(w)
}

func (w WorkspaceRelays) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WorkspaceRelays) Less(i, j int) bool {
	return w[i].LastSeen.Unix() < w[j].LastSeen.Unix()
}

func (w WorkspaceRelays) String() string {
	sort.Sort(w)
	if API.JsonOnly {
		return API.AsJSON(w)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"ID",
		"Version",
		"IP",
		"Port",
		"Last Seen",
		"Location",
	})

	for _, r := range w {
		rackID := ""
		if (r.Location.RackID != uuid.UUID{}) {
			rackID = fmt.Sprintf(
				"[%s]",
				CutUUID(r.Location.RackID.String()),
			)
		}
		location := fmt.Sprintf(
			"%s %s - RU %d %s",
			r.Location.AZ,
			r.Location.RackName,
			r.Location.RackUnitStart,
			rackID,
		)
		table.Append([]string{
			r.ID,
			r.Version,
			r.IpAddr,
			strconv.Itoa(r.SshPort),
			TimeStr(r.LastSeen),
			location,
		})
	}

	table.Render()
	return tableString.String()
}

func (w *Workspaces) GetRelays(name string) WorkspaceRelays {
	relays := make(WorkspaceRelays, 0)

	uri := fmt.Sprintf(
		"/workspace/%s/relay",
		url.PathEscape(name),
	)

	res := w.Do(w.Sling().Get(uri))
	if ok := res.Parse(&relays); !ok {
		panic(res)
	}

	return relays
}

func (w *Workspaces) GetRelayDevices(workspace string, relay string) DeviceList {
	devices := make(DeviceList, 0)

	uri := fmt.Sprintf(
		"/workspace/%s/relay/%s/device",
		url.PathEscape(workspace),
		url.PathEscape(relay),
	)

	res := w.Do(w.Sling().Get(uri))
	if ok := res.Parse(&devices); !ok {
		panic(res)
	}

	return devices
}

/******/

func (w *Workspaces) AddRack(workspace string, rackID uuid.UUID) WorkspaceRackSummaries {
	// I really explictly am not supporting the full functionality of this
	// endpoint. Tehnically, it supports updating the rack's serial number and
	// asset tag at the same time you add it to the workspace. That's.. yeah.
	// If you want to change those fields, use the function for updating a
	// rack's data, not this one.
	uri := fmt.Sprintf(
		"/workspace/%s/rack",
		url.PathEscape(workspace),
	)

	payload := struct {
		ID string `json:"id"`
	}{rackID.String()}

	// Ignoring the return because we're about to pivot into a different call.
	// If this call errors out, it'll panic on its own
	_ = w.Do(
		w.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	return w.GetRackSummaries(workspace)
}

func (w *Workspaces) RemoveRack(workspace string, rackID uuid.UUID) WorkspaceRackSummaries {
	uri := fmt.Sprintf(
		"/workspace/%s/rack/%s",
		url.PathEscape(workspace),
		url.PathEscape(rackID.String()),
	)

	// Ignoring the return because we're about to pivot into a different call.
	// If this call errors out, it'll panic on its own
	_ = w.Do(w.Sling().New().Delete(uri))

	return w.GetRackSummaries(workspace)
}

func (w *Workspaces) FindRackID(workspace string, id string) (bool, uuid.UUID) {

	summaries := make([]WorkspaceRackSummary, 0)

	summaries = append(
		summaries,
		w.GetRackSummaries(workspace).Slice()...,
	)

	ws := w.GetByName(workspace)
	if ws.Parent != "" {
		summaries = append(
			summaries,
			w.GetRackSummaries(ws.Parent).Slice()...,
		)
	}

	ids := make([]uuid.UUID, 0)
	for _, s := range summaries {
		ids = append(ids, s.ID)
	}

	return FindUUID(id, ids)
}

/******/

func init() {
	App.Command("workspaces", "Get a list of all workspaces you have access to", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Println(API.Workspaces().GetAll())
		}
	})

	App.Command("workspace", "Deal with a single workspace", func(cmd *cli.Cmd) {
		var workspaceName string
		workspaceNameArg := cmd.StringArg(
			"NAME",
			"",
			"The string name of the workspace")

		cmd.Spec = "NAME"
		cmd.Before = func() {
			workspaceName = *workspaceNameArg
			// TODO(sungo): should we verify that the workspace exists?
		}

		cmd.Command("get", "Get information about a single workspace, using its name", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Workspaces().GetByName(workspaceName))
			}
		})

		cmd.Command("create", "Create a new subworkspace", func(cmd *cli.Cmd) {
			nameArg := cmd.StringArg(
				"SUB",
				"",
				"Name of the new subworkspace",
			)

			descOpt := cmd.StringOpt(
				"description desc",
				"",
				"A description of the workspace",
			)

			cmd.Spec = "SUB [OPTIONS]"
			cmd.Action = func() {
				fmt.Println(API.Workspaces().Create(
					workspaceName,
					*nameArg,
					*descOpt,
				))
			}
		})

		cmd.Command("add", "Add various structures to a single workspace", func(cmd *cli.Cmd) {
			cmd.Command("user", "Add a user to a workspace", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email of the user to add to the workspace. Does *not* create the user",
				)

				roleOpt := cmd.StringOpt(
					"role",
					"ro",
					"The role for the user. One of: "+prettyWorkspaceRoleList(),
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target user, notifying them of the change",
				)

				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					if !okWorkspaceRole(*roleOpt) {
						panic(fmt.Errorf(
							"'role' value must be one of: %s",
							prettyWorkspaceRoleList(),
						))
					}

					if ok := API.Workspaces().AddOrModifyUser(
						workspaceName,
						*userEmailArg,
						*roleOpt,
						*sendEmailOpt,
					); ok {

						fmt.Println(API.Workspaces().GetUsers(workspaceName))

					} else {
						// It should be impossible to reach this
						// code as the lower code panics in all
						// known failure conditions.
						panic(errors.New("failure"))
					}
				}
			})

			cmd.Command("rack", "Add a rack to a workspace", func(cmd *cli.Cmd) {
				idArg := cmd.StringArg(
					"UUID",
					"",
					"The UUID of the rack to add. Short UUIDs (first segment) accepted",
				)

				cmd.Spec = "UUID"
				cmd.Action = func() {

					var rackID uuid.UUID
					var ok bool

					if ok, rackID = API.Workspaces().FindRackID(workspaceName, *idArg); !ok {
						panic(errors.New("could not locate the rack in either this workspace or its parent"))
					}

					fmt.Println(API.Workspaces().AddRack(
						workspaceName,
						rackID,
					))
				}
			})
		})

		cmd.Command("update", "Update various structures in a single workspace", func(cmd *cli.Cmd) {
			cmd.Command("user", "Update a user in a workspace", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email of the user to modify",
				)

				roleOpt := cmd.StringOpt(
					"role",
					"ro",
					"The role for the user. One of: "+prettyWorkspaceRoleList(),
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target user, notifying them of the change",
				)

				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					if !okWorkspaceRole(*roleOpt) {
						panic(fmt.Errorf(
							"'role' value must be one of: %s",
							prettyWorkspaceRoleList(),
						))
					}

					if ok := API.Workspaces().AddOrModifyUser(
						workspaceName,
						*userEmailArg,
						*roleOpt,
						*sendEmailOpt,
					); ok {

						fmt.Println(API.Workspaces().GetUsers(workspaceName))
					} else {
						// It should be impossible to reach this
						// code as the lower code panics in all
						// known failure conditions.
						panic(errors.New("failure"))
					}
				}
			})
		})

		cmd.Command("remove rm", "Remove various structures from a single workspace", func(cmd *cli.Cmd) {
			cmd.Command("user", "Remove a user from a workspace", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email of the user to modify",
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target user, notifying them of the change",
				)

				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					if ok := API.Workspaces().RemoveUser(
						workspaceName,
						*userEmailArg,
						*sendEmailOpt,
					); ok {
						fmt.Println(API.Workspaces().GetUsers(workspaceName))
					} else {
						// It should be impossible to reach this
						// code as the lower code panics in all
						// known failure conditions.
						panic(errors.New("failure"))
					}
				}
			})

			cmd.Command("rack", "Remove a rack from a workspace", func(cmd *cli.Cmd) {
				idArg := cmd.StringArg(
					"UUID",
					"",
					"The UUID of the rack to remove. Short UUIDs (first segment) accepted",
				)

				cmd.Spec = "UUID"
				cmd.Action = func() {

					var id uuid.UUID
					var err error
					if id, err = uuid.FromString(*idArg); err != nil {
						if ok, rackID := API.Workspaces().FindRackID(workspaceName, *idArg); ok {
							id = rackID
						} else {
							panic(errors.New("could not locate the rack in either this workspace or its parent"))
						}
					}

					if (id == uuid.UUID{}) {
						panic(errors.New("could not locate the rack in either this workspace or its parent"))
					}

					fmt.Println(API.Workspaces().RemoveRack(workspaceName, id))

				}
			})

		})

		cmd.Command("relays", "Get a list of relays assigned to a single workspace", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Workspaces().GetRelays(workspaceName))
			}
		})

		cmd.Command("relay", "Deal with a single relay", func(cmd *cli.Cmd) {
			cmd.Command("get", "Get information about a single relay", func(cmd *cli.Cmd) {
				relayArg := cmd.StringArg(
					"RELAY",
					"",
					"ID of the relay",
				)

				cmd.Spec = "RELAY"
				cmd.Action = func() {
					relays := API.Workspaces().GetRelays(workspaceName)

					for _, relay := range relays {
						if relay.ID != *relayArg {
							continue
						}

						fmt.Println(relay)
						return
					}
					panic(errors.New("relay not found"))
				}
			})

			cmd.Command("devices", "Get the device list for a relay", func(cmd *cli.Cmd) {
				relayArg := cmd.StringArg(
					"RELAY",
					"",
					"ID of the relay",
				)

				cmd.Spec = "RELAY"
				cmd.Action = func() {
					fmt.Println(API.Workspaces().GetRelayDevices(
						workspaceName,
						*relayArg,
					))
				}
			})
		})

		cmd.Command("children subs", "Get a list of a workspace's children", func(cmd *cli.Cmd) {

			allOpt := cmd.BoolOpt(
				"all",
				false,
				"Retrieve all children, not just the direct lineage",
			)

			// TODO(sungo): tree mode?
			cmd.Action = func() {
				if *allOpt {
					fmt.Println(API.Workspaces().GetChildren(workspaceName))
				} else {
					fmt.Println(API.Workspaces().GetDirectChildren(workspaceName))
				}
			}
		})

		cmd.Command("devices", "Get a list of devices in a single workspace, by name", func(cmd *cli.Cmd) {
			healthOpt := cmd.StringOpt(
				"health",
				"",
				"Filter by the 'health' field. Value must be one of: "+prettyDeviceHealthList(),
			)

			var validatedSetByUser bool
			validatedOpt := cmd.Bool(cli.BoolOpt{
				Name:      "validated",
				Value:     false,
				Desc:      "Filter by the 'validated' field",
				SetByUser: &validatedSetByUser,
			})

			cmd.Action = func() {
				if !validatedSetByUser {
					validatedOpt = nil
				}

				if *healthOpt != "" {
					if !okHealth(*healthOpt) {
						panic(fmt.Errorf("'health' value must be one of: %s", prettyDeviceHealthList()))

					}
				}
				fmt.Println(API.Workspaces().GetDevices(
					workspaceName,
					*healthOpt,
					validatedOpt,
				))
			}
		})

		cmd.Command("users", "Operate on the users assigned to a workspace", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Workspaces().GetUsers(workspaceName))
			}
		})

		cmd.Command("racks", "Get a progress summary for each rack", func(cmd *cli.Cmd) {
			phaseOpt := cmd.StringOpt(
				"phase",
				"",
				"Filter based on phase name",
			)

			roleOpt := cmd.StringOpt(
				"role",
				"",
				"Filter on role name",
			)

			cmd.Action = func() {
				ret := API.Workspaces().GetRackSummaries(workspaceName)

				summaries := make(WorkspaceRackSummaries)

				for az, summary := range ret {
					for _, s := range summary {
						save := true

						if *phaseOpt != "" {
							if s.Phase != *phaseOpt {
								save = false
							}
						}

						if *roleOpt != "" {
							if s.RoleName != *roleOpt {
								save = false
							}
						}

						if save {
							if _, ok := summaries[az]; !ok {
								summaries[az] = make([]WorkspaceRackSummary, 0)
							}

							summaries[az] = append(summaries[az], s)
						}
					}
				}

				fmt.Println(summaries)
			}
		})
	})
}
