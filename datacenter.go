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
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Datacenters struct {
	*Conch
}

func (c *Conch) Datacenters() *Datacenters {
	return &Datacenters{c}
}

type Datacenter struct {
	ID         uuid.UUID `json:"id" faker:"uuid"`
	Vendor     string    `json:"vendor"`
	VendorName string    `json:"vendor_name,omitempty"`
	Region     string    `json:"region"`
	Location   string    `json:"location"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

func (d Datacenter) String() string {
	if API.JsonOnly {
		return API.AsJSON(d)
	}

	t, err := NewTemplate().Parse(datacenterTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, d); err != nil {
		panic(err)
	}

	return buf.String()
}

type DatacenterList []Datacenter

func (dl DatacenterList) String() string {
	if API.JsonOnly {
		return API.AsJSON(dl)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"ID",
		"Vendor",
		"Vendor Name",
		"Region",
		"Location",
	})

	for _, d := range dl {
		table.Append([]string{
			CutUUID(d.ID.String()),
			d.Vendor,
			d.VendorName,
			d.Region,
			d.Location,
		})
	}

	table.Render()
	return tableString.String()
}

/****/

func (d *Datacenters) GetAll() DatacenterList {
	dl := make(DatacenterList, 0)
	res := d.Do(d.Sling().Get("/dc"))
	if ok := res.Parse(&dl); !ok {
		panic(res)
	}
	return dl
}

func (d *Datacenters) FindDatacenterID(id string) (bool, uuid.UUID) {
	ids := make([]uuid.UUID, 0)
	for _, datacenter := range d.GetAll() {
		ids = append(ids, datacenter.ID)
	}

	return FindUUID(id, ids)
}

func (d *Datacenters) Get(id uuid.UUID) Datacenter {
	var dc Datacenter
	uri := fmt.Sprintf(
		"/dc/%s",
		url.PathEscape(id.String()),
	)

	res := d.Do(d.Sling().Get(uri))
	if ok := res.Parse(&dc); !ok {
		panic(res)
	}
	return dc
}

func (d *Datacenters) Update(id uuid.UUID, region string, vendor string, location string, vendorName string) Datacenter {

	var dc Datacenter

	uri := fmt.Sprintf(
		"/dc/%s",
		url.PathEscape(id.String()),
	)

	payload := make(map[string]string)

	if region != "" {
		payload["region"] = region
	}

	if vendor != "" {
		payload["vendor"] = vendor
	}

	if vendorName != "" {
		payload["vendor_name"] = vendorName
	}

	if location != "" {
		payload["location"] = location
	}

	if len(payload) == 0 {
		panic(errors.New("at least one property must be defined: region, vendor, vendor name, location"))
	}

	res := d.Do(
		d.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	if ok := res.Parse(&dc); !ok {
		panic(res)
	}

	return dc
}

func (d *Datacenters) Create(region string, vendor string, location string, vendorName string) Datacenter {
	payload := make(map[string]string)
	if vendor == "" {
		panic(errors.New("'vendor' cannot be empty"))
	}
	payload["vendor"] = vendor

	if region == "" {
		panic(errors.New("'region' cannot be empty"))
	}
	payload["region"] = region

	if location == "" {
		panic(errors.New("'location' cannot be empty"))
	}
	payload["location"] = location

	if vendorName != "" {
		payload["vendor_name"] = vendorName
	}

	/**/

	var dc Datacenter

	// We get a 303 on success
	res := d.Do(
		d.Sling().New().Post("/dc").
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)
	if ok := res.Parse(&dc); !ok {
		panic(res)
	}

	return dc
}

func (d *Datacenters) CreateFromStruct(newDC Datacenter) Datacenter {
	return d.Create(
		newDC.Region,
		newDC.Vendor,
		newDC.Location,
		newDC.VendorName,
	)
}

func (d *Datacenters) Delete(id uuid.UUID) {
	uri := fmt.Sprintf(
		"/dc/%s",
		url.PathEscape(id.String()),
	)

	res := d.Do(d.Sling().New().Delete(uri))

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

func (d *Datacenters) GetRooms(id uuid.UUID) RoomList {
	rooms := make(RoomList, 0)
	uri := fmt.Sprintf(
		"/dc/%s/rooms",
		url.PathEscape(id.String()),
	)

	res := d.Do(d.Sling().Get(uri))
	if ok := res.Parse(&rooms); !ok {
		panic(res)
	}

	return rooms
}

/****/

func init() {

	App.Command("datacenters", "Work with all the datacenters you have access to", func(cmd *cli.Cmd) {
		cmd.Before = RequireSysAdmin
		cmd.Command("get", "Get a list of all datacenters", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Datacenters().GetAll())
			}
		})

		cmd.Command("create", "Create a single datacenter", func(cmd *cli.Cmd) {
			var (
				vendorOpt     = cmd.StringOpt("vendor", "", "Vendor")
				regionOpt     = cmd.StringOpt("region", "", "Region")
				locationOpt   = cmd.StringOpt("location", "", "Location")
				vendorNameOpt = cmd.StringOpt("vendor-name", "", "Vendor Name")
			)

			cmd.Spec = "--vendor --region --location [OPTIONS]"
			cmd.Action = func() {
				// The user can be very silly and supply something like
				// '--vendor ""' which will pass the cli lib's requirement
				// check but is still crap
				if *vendorOpt == "" {
					panic(errors.New("--vendor is required"))
				}
				if *regionOpt == "" {
					panic(errors.New("--region is required"))
				}
				if *locationOpt == "" {
					panic(errors.New("--location is required"))
				}

				fmt.Println(API.Datacenters().Create(
					*regionOpt,
					*vendorOpt,
					*locationOpt,
					*vendorNameOpt,
				))
			}
		})

	})

	App.Command("datacenter", "Deal with a single datacenter", func(cmd *cli.Cmd) {
		var datacenterID uuid.UUID

		idArg := cmd.StringArg(
			"UUID",
			"",
			"The UUID of the datacenter. Short UUIDs (first segment) accepted",
		)

		cmd.Spec = "UUID"

		cmd.Before = func() {
			RequireSysAdmin()
			var ok bool

			if ok, datacenterID = API.Datacenters().FindDatacenterID(*idArg); !ok {
				panic(errors.New("could not find the datacenter"))
			}
		}

		cmd.Command("get", "Information about a single datacenter", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Datacenters().Get(datacenterID))
			}
		})

		cmd.Command("delete", "Delete a single datacenter", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				// Lower layers panic if there's a problem
				API.Datacenters().Delete(datacenterID)

				fmt.Println(API.Datacenters().GetAll())
			}
		})

		cmd.Command("update", "Update a single datacenter", func(cmd *cli.Cmd) {
			regionOpt := cmd.StringOpt(
				"region",
				"",
				"Region identifier",
			)
			vendorOpt := cmd.StringOpt(
				"vendor",
				"",
				"Vendor",
			)
			vendorNameOpt := cmd.StringOpt(
				"vendor-name",
				"",
				"Vendor Name",
			)
			locationOpt := cmd.StringOpt(
				"location",
				"",
				"Location",
			)

			cmd.Action = func() {
				var count int
				if *regionOpt != "" {
					count++
				}
				if *vendorOpt != "" {
					count++
				}
				if *vendorNameOpt != "" {
					count++
				}
				if *locationOpt != "" {
					count++
				}

				if count == 0 {
					panic(errors.New("one option must be provided"))
				}

				fmt.Println(API.Datacenters().Update(
					datacenterID,
					*regionOpt,
					*vendorOpt,
					*locationOpt,
					*vendorNameOpt,
				))
			}
		})

		cmd.Command("rooms", "Get the room list for a single datacenter", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Datacenters().GetRooms(datacenterID))
			}
		})
	})
}
