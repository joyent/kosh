// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type RackRoles struct {
	*Conch
}

func (c *Conch) RackRoles() *RackRoles {
	return &RackRoles{c}
}

/****/

type RackRoleList []RackRole

func (r RackRoleList) Len() int {
	return len(r)
}

func (r RackRoleList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RackRoleList) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

/****/

type RackRole struct {
	ID       uuid.UUID `json:"id" faker:"uuid"`
	Name     string    `json:"name"`
	RackSize int       `json:"rack_size"`
	Created  time.Time `json:"created" faker:"-"`
	Updated  time.Time `json:"updated" faker:"-"`
}

func (rl RackRoleList) String() string {
	sort.Sort(rl)
	if API.JsonOnly {
		return API.AsJSON(rl)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Name",
		"RackSize",
		"Created",
		"Updated",
	})

	for _, r := range rl {
		table.Append([]string{
			r.Name,
			strconv.Itoa(r.RackSize),
			TimeStr(r.Created),
			TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()

}

func (r RackRole) String() string {
	if API.JsonOnly {
		return API.AsJSON(r)
	}

	t, err := NewTemplate().Parse(rackRoleTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}

/****/

func (r *RackRoles) GetAll() RackRoleList {
	rl := make(RackRoleList, 0)
	res := r.Do(r.Sling().Get("/rack_role"))
	if ok := res.Parse(&rl); !ok {
		panic(res)
	}
	return rl
}

func (r *RackRoles) Get(id uuid.UUID) RackRole {
	var role RackRole

	uri := fmt.Sprintf(
		"/rack_role/%s",
		url.PathEscape(id.String()),
	)

	res := r.Do(r.Sling().Get(uri))
	if ok := res.Parse(&role); !ok {
		panic(res)
	}
	return role
}

func (r *RackRoles) GetByName(name string) RackRole {
	var role RackRole

	uri := fmt.Sprintf(
		"/rack_role/name=%s",
		url.PathEscape(name),
	)

	res := r.Do(r.Sling().Get(uri))
	if ok := res.Parse(&role); !ok {
		panic(res)
	}
	return role
}

func (r *RackRoles) FindID(name string) (bool, uuid.UUID) {
	var role RackRole

	uri := fmt.Sprintf(
		"/rack_role/name=%s",
		url.PathEscape(name),
	)

	res := r.DoBadly(r.Sling().Get(uri))
	if res.IsError() {
		return false, uuid.UUID{}
	}

	return res.Parse(&role), role.ID
}

func (r *RackRoles) Create(name string, rackSize int) RackRole {
	if name == "" {
		panic(errors.New("'name' is required"))
	}

	if rackSize == 0 {
		panic(errors.New("'rackSize' is required and cannot be 0"))
	}

	payload := make(map[string]interface{})
	payload["name"] = name
	payload["rack_size"] = rackSize

	var role RackRole

	// We get a 303 on success
	res := r.Do(
		r.Sling().New().Post("/rack_role").
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	if ok := res.Parse(&role); !ok {
		panic(res)
	}

	return role
}

func (r *RackRoles) Update(id uuid.UUID, newName string, rackSize int) RackRole {
	payload := make(map[string]interface{})
	if newName != "" {
		payload["name"] = newName
	}
	if rackSize > 0 {
		payload["rack_size"] = rackSize
	}

	var role RackRole

	uri := fmt.Sprintf(
		"/rack_role/%s",
		url.PathEscape(id.String()),
	)

	// We get a 303 on success
	res := r.Do(
		r.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	if ok := res.Parse(&role); !ok {
		panic(res)
	}

	return role
}

func (r *RackRoles) Delete(id uuid.UUID) {
	uri := fmt.Sprintf("/rack_role/%s", url.PathEscape(id.String()))
	res := r.Do(r.Sling().New().Delete(uri))

	if res.StatusCode() != 204 {
		// I know this is weird. Like in other places, it should be impossible
		// to reach here unless the status code is 204. The API returns 204
		// (which gets us here) or 409 (which will explode before it gets here).
		// If we got here via some other code, then there's some new behavior
		// that we need to know about.
		panic(res)
	}
}

/****/

func init() {
	App.Command("roles", "Work with datacenter rack roles", func(cmd *cli.Cmd) {
		cmd.Before = RequireSysAdmin
		cmd.Command("get", "Get a list of all rack roles", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.RackRoles().GetAll()) }
		})

		cmd.Command("create", "Create a new rack role", func(cmd *cli.Cmd) {
			var (
				nameOpt     = cmd.StringOpt("name", "", "The name of the role")
				rackSizeOpt = cmd.IntOpt("rack-size", 0, "Size of the rack necessary for this role")
			)

			cmd.Spec = "--name --rack-size"
			cmd.Action = func() {
				if *nameOpt == "" {
					panic(errors.New("--name is required"))
				}

				if *rackSizeOpt == 0 {
					panic(errors.New("--rack-size is required and cannot be 0"))
				}

				fmt.Println(API.RackRoles().Create(*nameOpt, *rackSizeOpt))
			}
		})
	})

	App.Command("role", "Work with a single rack role", func(cmd *cli.Cmd) {
		var roleID uuid.UUID

		nameArg := cmd.StringArg(
			"NAME",
			"",
			"The name of the rack role",
		)

		cmd.Spec = "NAME"

		cmd.Before = func() {
			RequireSysAdmin()
			var ok bool

			if ok, roleID = API.RackRoles().FindID(*nameArg); !ok {
				panic(errors.New("could not find the role"))
			}
		}

		cmd.Command("get", "Get information about a single rack role", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.RackRoles().Get(roleID)) }
		})

		cmd.Command("update", "Update information about a single rack role", func(cmd *cli.Cmd) {
			var (
				nameOpt     = cmd.StringOpt("name", "", "The name of the role")
				rackSizeOpt = cmd.IntOpt("rack-size", 0, "Size of the rack necessary for this role")
			)

			cmd.Action = func() {
				fmt.Println(API.RackRoles().Update(roleID, *nameOpt, *rackSizeOpt))
			}
		})

		cmd.Command("delete", "Delete a single rack role", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				API.RackRoles().Delete(roleID)
				fmt.Println(API.RackRoles().GetAll())
			}
		})
	})
}
