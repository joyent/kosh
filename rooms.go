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
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Rooms struct {
	*Conch
}

func (c *Conch) Rooms() *Rooms {
	return &Rooms{c}
}

/****/

// This is called DatacenterRoomsDetailed in the json schema
type RoomList []Room

func (r RoomList) Len() int {
	return len(r)
}

func (r RoomList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomList) Less(i, j int) bool {
	return r[i].Alias < r[j].Alias
}

func (dr RoomList) String() string {
	sort.Sort(dr)
	if API.JsonOnly {
		return API.AsJSON(dr)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Alias",
		"AZ",
		"Vendor Name",
		"Datacenter ID",
		"Created",
		"Updated",
	})

	for _, r := range dr {
		table.Append([]string{
			r.Alias,
			r.AZ,
			r.VendorName,
			CutUUID(r.DatacenterID.String()),
			TimeStr(r.Created),
			TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()

}

// This is called DatacenterRoomDetailed in the json schema
type Room struct {
	ID           uuid.UUID `json:"id" faker:"uuid"`
	AZ           string    `json:"az"`
	Alias        string    `json:"alias"`
	VendorName   string    `json:"vendor_name,omitempty"`
	DatacenterID uuid.UUID `json:"datacenter_id" faker:"uuid"`
	Created      time.Time `json:"created" faker:"-"`
	Updated      time.Time `json:"updated" faker:"-"`
}

func (r Room) String() string {
	if API.JsonOnly {
		return API.AsJSON(r)
	}

	t, err := NewTemplate().Parse(roomTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}

// Accepting partial UUIDs, full UUIDs, and 'alias'
func (r *Rooms) FindID(id string) (bool, uuid.UUID) {
	ids := make([]uuid.UUID, 0)
	for _, room := range r.GetAll() {
		if room.Alias == id {
			return true, room.ID
		}
		ids = append(ids, room.ID)
	}

	return FindUUID(id, ids)
}

func (r *Rooms) GetAll() RoomList {
	rl := make(RoomList, 0)
	res := r.Do(r.Sling().Get("/room"))
	if ok := res.Parse(&rl); !ok {
		panic(res)
	}
	return rl
}

func (r *Rooms) Get(id uuid.UUID) Room {
	var room Room
	uri := fmt.Sprintf(
		"/room/%s",
		url.PathEscape(id.String()),
	)

	res := r.Do(r.Sling().Get(uri))
	if ok := res.Parse(&room); !ok {
		panic(res)
	}
	return room
}

func (r *Rooms) Create(datacenterID uuid.UUID, az string, alias string, vendorName string) Room {
	payload := make(map[string]string)
	if (datacenterID == uuid.UUID{}) {
		panic(errors.New("'datacenterID' cannot be empty"))
	}
	payload["datacenter_id"] = datacenterID.String()

	if az == "" {
		panic(errors.New("'az' cannot be empty"))
	}
	payload["az"] = az

	if alias == "" {
		panic(errors.New("'alias' cannot be empty"))
	}
	payload["alias"] = alias

	if vendorName != "" {
		payload["vendor_name"] = vendorName
	}

	/**/

	var room Room

	// We get a 303 on success
	res := r.Do(
		r.Sling().New().Post("/room").
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)
	if ok := res.Parse(&room); !ok {
		panic(res)
	}

	return room
}

func (r *Rooms) Update(id uuid.UUID, datacenterID uuid.UUID, az string, alias string, vendorName string) Room {
	payload := make(map[string]string)
	if (datacenterID != uuid.UUID{}) {
		payload["datacenter_id"] = datacenterID.String()
	}

	if az != "" {
		payload["az"] = az
	}

	if alias != "" {
		payload["alias"] = alias
	}

	if vendorName != "" {
		payload["vendor_name"] = vendorName
	}

	/**/

	var room Room
	uri := fmt.Sprintf(
		"/room/%s",
		url.PathEscape(id.String()),
	)

	// We get a 303 on success
	res := r.Do(
		r.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)
	if ok := res.Parse(&room); !ok {
		panic(res)
	}

	return room
}

func (r *Rooms) Delete(id uuid.UUID) {
	uri := fmt.Sprintf("/room/%s", url.PathEscape(id.String()))
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

func (r *Rooms) Racks(id uuid.UUID) RackList {
	uri := fmt.Sprintf("/room/%s/racks", url.PathEscape(id.String()))

	rl := make(RackList, 0)
	res := r.Do(r.Sling().New().Get(uri))
	if ok := res.Parse(&rl); !ok {
		panic(res)
	}
	return rl
}

func init() {
	App.Command("rooms", "Work with datacenter rooms", func(cmd *cli.Cmd) {
		cmd.Before = RequireSysAdmin
		cmd.Command("get", "Get a list of all rooms", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Rooms().GetAll())
			}
		})
		cmd.Command("create", "Create a single room", func(cmd *cli.Cmd) {
			var (
				aliasOpt        = cmd.StringOpt("alias", "", "Alias")
				azOpt           = cmd.StringOpt("az", "", "AZ")
				datacenterIdOpt = cmd.StringOpt("datacenter-id", "", "Datacenter UUID (first segment of UUID accepted)")
				vendorNameOpt   = cmd.StringOpt("vendor-name", "", "Vendor Name")
			)

			cmd.Spec = "--datacenter-id --alias --az [OPTIONS]"
			cmd.Action = func() {
				// The user can be very silly and supply something like
				// '--alias ""' which will pass the cli lib's requirement
				// check but is still crap
				if *aliasOpt == "" {
					panic(errors.New("--alias is required"))
				}
				if *azOpt == "" {
					panic(errors.New("--az is required"))
				}
				if *datacenterIdOpt == "" {
					panic(errors.New("--datacenter-id is required"))
				}

				var datacenterID uuid.UUID
				var ok bool
				if ok, datacenterID = API.Datacenters().FindDatacenterID(*datacenterIdOpt); !ok {
					panic(errors.New("could not find the datacenter"))
				}

				fmt.Println(API.Rooms().Create(
					datacenterID,
					*azOpt,
					*aliasOpt,
					*vendorNameOpt,
				))
			}
		})

	})

	App.Command("room", "Deal with a single datacenter room", func(cmd *cli.Cmd) {
		var roomID uuid.UUID

		aliasArg := cmd.StringArg(
			"ALIAS",
			"",
			"The unique alias of the datacenter room",
		)

		cmd.Spec = "ALIAS"

		cmd.Before = func() {
			RequireSysAdmin()
			var ok bool

			if ok, roomID = API.Rooms().FindID(*aliasArg); !ok {
				panic(errors.New("could not find the room"))
			}
		}

		cmd.Command("get", "Information about a single room", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Rooms().Get(roomID))
			}
		})

		cmd.Command("update", "Update information about a single room", func(cmd *cli.Cmd) {
			var (
				aliasOpt        = cmd.StringOpt("alias", "", "Alias")
				azOpt           = cmd.StringOpt("az", "", "AZ")
				datacenterIdOpt = cmd.StringOpt("datacenter-id", "", "Datacenter UUID (first segment of UUID accepted)")
				vendorNameOpt   = cmd.StringOpt("vendor-name", "", "Vendor Name")
			)

			cmd.Action = func() {
				var datacenterID uuid.UUID

				if *datacenterIdOpt != "" {
					var ok bool
					if ok, datacenterID = API.Datacenters().FindDatacenterID(*datacenterIdOpt); !ok {
						panic(errors.New("could not find the datacenter"))
					}
				}

				fmt.Println(API.Rooms().Update(
					roomID,
					datacenterID,
					*azOpt,
					*aliasOpt,
					*vendorNameOpt,
				))
			}
		})

		cmd.Command("delete", "Delete a single room", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				// Lower layers panic if there's a problem
				API.Rooms().Delete(roomID)

				fmt.Println(API.Rooms().GetAll())
			}
		})

		cmd.Command("racks", "View the racks assigned to a single room", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Rooms().Racks(roomID))
			}
		})

	})
}
