// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/tables"
	"github.com/joyent/kosh/template"
)

type Users struct {
	*Conch
}

func (c *Conch) Users() *Users {
	return &Users{c}
}

type UserAndRole struct {
	ID    uuid.UUID `json:"id" faker:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email" faker:"email"`
	Role  string    `json:"role"`
}

type UserAndRoles []UserAndRole

func (u UserAndRoles) Len() int {
	return len(u)
}

func (u UserAndRoles) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UserAndRoles) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}

func (ur UserAndRoles) String() string {
	sort.Sort(ur)
	if API.JsonOnly {
		return API.AsJSON(ur)
	}

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
	tables.TableToMarkdown(table)

	table.SetHeader([]string{
		"ID",
		"Name",
		"Email",
		"Role",
	})

	for _, u := range ur {
		table.Append([]string{
			u.ID.String(),
			u.Name,
			u.Email,
			u.Role,
		})
	}

	table.Render()
	return tableString.String()
}

/*****/

// In the json schema, DetailedUser is UserDetailed and DetailedUsers is UsersDetailed

type DetailedUser struct {
	ID                  uuid.UUID         `json:"id" faker:"uuid"`
	Name                string            `json:"name"`
	Email               string            `json:"email"`
	Created             time.Time         `json:"created"`
	LastLogin           time.Time         `json:"last_login,omitempty"`
	LastSeen            time.Time         `json:"last_seen,omitempty"`
	RefuseSessionAuth   bool              `json:"refuse_session_auth"`
	ForcePasswordChange bool              `json:"force_password_change"`
	IsAdmin             bool              `json:"is_admin"`
	Workspaces          WorkspaceAndRoles `json:"workspaces"`
	Organizations       OrgAndRoles       `json:"organizations"`
	Builds              interface{}       `json:"builds"` // TODO build support
}

func (u DetailedUser) String() string {
	if API.JsonOnly {
		return API.AsJSON(u)
	}

	t, err := template.NewTemplate().Parse(detailedUserTemplate)
	if err != nil {
		panic(err)
	}

	buf := &strings.Builder{}

	if err := t.Execute(buf, u); err != nil {
		panic(err)
	}

	return buf.String()
}

type DetailedUsers []DetailedUser

func (u DetailedUsers) Len() int {
	return len(u)
}

func (u DetailedUsers) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u DetailedUsers) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}

func (u *Users) Me() (user DetailedUser) {
	res := u.Do(u.Sling().Get("/user/me"))
	if ok := res.Parse(&user); !ok {
		panic(res)
	}
	ret := make(WorkspaceAndRoles, 0)
	cache := make(map[uuid.UUID]string)

	for _, ws := range user.Workspaces {
		if (ws.ParentID != uuid.UUID{}) {
			if _, ok := cache[ws.ParentID]; !ok {
				cache[ws.ParentID] = API.Workspaces().Get(ws.ParentID).Name
			}

			ws.Parent = cache[ws.ParentID]
		}

		if (ws.RoleVia != uuid.UUID{}) {
			if _, ok := cache[ws.RoleVia]; !ok {
				cache[ws.RoleVia] = API.Workspaces().Get(ws.RoleVia).Name
			}

			ws.Via = cache[ws.RoleVia]
		}

		ret = append(ret, ws)
	}

	user.Workspaces = ret

	return user
}

type UserSettings map[string]interface{}

func (u UserSettings) String() string {
	if API.JsonOnly {
		return API.AsJSON(u)
	}

	keys := make([]string, 0)
	for setting := range u {
		keys = append(keys, setting)
	}
	sort.Strings(keys)

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
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
	return tableString.String()
}

func (u *Users) MySettings() UserSettings {
	settings := make(UserSettings)

	res := u.Do(u.Sling().Get("/user/me/settings"))
	if ok := res.Parse(&settings); !ok {
		panic(res)
	}

	return settings
}

func (u *Users) GetMySetting(name string) interface{} {
	uri := fmt.Sprintf(
		"/user/me/settings/%s",
		url.PathEscape(name),
	)

	data := make(map[string]interface{})

	res := u.DoBadly(u.Sling().Get(uri))
	if res.StatusCode() == 404 {
		return ""
	}
	if res.IsError() {
		panic(res)
	}

	if ok := res.Parse(&data); !ok {
		panic(res)
	}

	return data[name]
}

func (u *Users) SetMySetting(name string, value string) interface{} {
	var userData interface{}

	if err := json.Unmarshal([]byte(value), &userData); err != nil {
		// If the value doesn't parse properly as JSON, we assume it's
		// literal. This catches the single-value case where we want
		// { "foo": "bar" } by just letting the user pass in a name of
		// "foo" and a value of "bar"

		// The perhaps surprising side effect is that crappy JSON will
		// enter the database as a string.
		userData = value
	}

	data := make(map[string]interface{})
	data[name] = userData

	uri := fmt.Sprintf(
		"/user/me/settings/%s",
		url.PathEscape(name),
	)

	// This endpoint either returns errors or a 204. We catch errors elsewhere.
	_ = u.Do(
		u.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(data),
	)

	// Pivot into a get so the caller can make sure the data was stored properly
	return u.GetMySetting(name)
}

func (u *Users) DeleteMySetting(name string) {
	uri := fmt.Sprintf(
		"/user/me/settings/%s",
		url.PathEscape(name),
	)

	res := u.DoBadly(u.Sling().Delete(uri))
	if res.StatusCode() == 404 {
		return
	}
	if res.IsError() {
		panic(res)
	}
}

/*****/

func init() {
	App.Command("whoami", "Display details of the current user", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Println(API.Users().Me())
		}
	})

	App.Command("user", "Commands for dealing with the current user (you)", func(cmd *cli.Cmd) {
		cmd.Command("profile", "View your Conch profile", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Users().Me())
			}
		})

		cmd.Command("settings", "Get the settings for the current user", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Users().MySettings())
			}
		})

		cmd.Command("setting", "Commands for dealing with a single setting for the current user", func(cmd *cli.Cmd) {
			settingNameArg := cmd.StringArg("NAME", "", "The string name of a setting")

			cmd.Spec = "NAME"

			cmd.Command("get", "Get a setting for the current user", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					setting := API.Users().GetMySetting(*settingNameArg)
					if API.JsonOnly {
						out := make(map[string]interface{})
						out[*settingNameArg] = setting
						API.PrintJSON(out)
						return
					}

					fmt.Printf("%s\n", setting)
				}
			})

			cmd.Command("set", "Set a setting for the current user", func(cmd *cli.Cmd) {
				valueArg := cmd.StringArg("VALUE", "", "The new value of the setting")

				cmd.Spec = "VALUE"

				cmd.Action = func() {
					setting := API.Users().SetMySetting(
						*settingNameArg,
						*valueArg,
					)

					if API.JsonOnly {
						out := make(map[string]interface{})
						out[*settingNameArg] = setting
						API.PrintJSON(out)
						return
					}

					fmt.Printf("%s\n", setting)
				}
			})

			cmd.Command("delete", "Delete a setting for the current user", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					API.Users().DeleteMySetting(*settingNameArg)
				}
			})
		})
	})
}
