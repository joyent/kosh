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
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/tables"
	"github.com/joyent/kosh/template"
)

type Validations struct {
	*Conch
}

func (c *Conch) Validations() *Validations {
	return &Validations{c}
}

/***/
type ValidationPlan struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

func (v ValidationPlan) String() string {
	if API.JsonOnly {
		return API.AsJSON(v)
	}

	t, err := template.NewTemplate().Parse(validationPlanTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, v); err != nil {
		panic(err)
	}

	return buf.String()

}

type ValidationPlans []ValidationPlan

func (v ValidationPlans) Len() int {
	return len(v)
}

func (v ValidationPlans) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ValidationPlans) Less(i, j int) bool {
	return v[i].Name < v[j].Name
}

func (v ValidationPlans) String() string {
	sort.Sort(v)
	if API.JsonOnly {
		return API.AsJSON(v)
	}
	if len(v) == 0 {
		return ""
	}

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
	tables.TableToMarkdown(table)
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
			p.Name,
			p.Description,
			p.Created.String(),
		})
	}

	table.Render()
	return tableString.String()
}

func (v Validations) GetPlan(id uuid.UUID) (vp ValidationPlan) {
	if (id == uuid.UUID{}) {
		return vp
	}

	uri := fmt.Sprintf("/validation_plan/%s", url.PathEscape(id.String()))
	res := v.Do(v.Sling().New().Get(uri))
	if ok := res.Parse(&vp); !ok {
		panic(res)
	}

	return vp
}

func (v Validations) GetAllPlans() (list ValidationPlans) {

	res := v.Do(v.Sling().New().Get("/validation_plan"))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}

	return
}

func (v Validations) FindPlanID(id string) (bool, uuid.UUID) {
	ids := make([]uuid.UUID, 0)
	for _, vp := range v.GetAllPlans() {
		ids = append(ids, vp.ID)
	}

	return FindUUID(id, ids)
}

func (v Validations) GetPlanByName(name string) (vp ValidationPlan) {

	uri := fmt.Sprintf("/validation_plan/%s", url.PathEscape(name))
	res := v.Do(v.Sling().New().Get(uri))
	if ok := res.Parse(&vp); !ok {
		panic(res)
	}

	return vp
}

/***/

type ValidationResult struct {
	ID                uuid.UUID `json:"id,omitempty"`
	Category          string    `json:"category"`
	Component         string    `json:"component,omitempty"`
	HardwareProductID uuid.UUID `json:"hardware_product_id"`
	Hint              string    `json:"hint,omitempty"`
	Message           string    `json:"message"`
	Order             int       `json:"order"`
	Status            string    `json:"status"`
	ValidationID      uuid.UUID `json:"validation_id"`
}

type ValidationResults []ValidationResult

func (v ValidationResults) Len() int {
	return len(v)
}

func (v ValidationResults) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ValidationResults) Less(i, j int) bool {
	return v[i].Category < v[j].Category
}

func (v ValidationResults) String() string {
	sort.Sort(v)
	if API.JsonOnly {
		return API.AsJSON(v)
	}
	if len(v) == 0 {
		return ""
	}

	tableString := &strings.Builder{}
	table := tables.NewTable(tableString)
	tables.TableToMarkdown(table)
	table.SetRowLine(true)

	table.SetHeader([]string{
		"Status",
		"Category",
		"Component",
		"Message",
	})

	for _, r := range v {
		table.Append([]string{
			r.Status,
			r.Category,
			r.Component,
			r.Message,
		})
	}

	table.Render()
	return tableString.String()
}

/***/

type ValidationStateWithResults struct {
	ID               uuid.UUID `json:"id"`
	Completed        time.Time `json:"completed"`
	Created          time.Time `json:"created"`
	DeviceID         uuid.UUID `json:"device_id"`
	Status           string    `json:"status"`
	ValidationPlanID uuid.UUID `json:"validation_plan_id"`
	DeviceReportID   uuid.UUID `json:"device_report_id"`

	Results ValidationResults `json:"results"`
}

func (v ValidationStateWithResults) String() string {
	if API.JsonOnly {
		return API.AsJSON(v)
	}

	type extendedVsR struct {
		ValidationStateWithResults
		ValidationPlan ValidationPlan `json:"-"`
	}

	out := extendedVsR{
		v,
		API.Validations().GetPlan(v.ValidationPlanID),
	}

	t, err := template.NewTemplate().Parse(validationStateWithResultsTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, out); err != nil {
		panic(err)
	}

	return buf.String()
}

func init() {

	App.Command("validation", "Work with validations", func(cmd *cli.Cmd) {
		v := API.Validations()
		cmd.Command("plans", "Work with validation plans", func(cmd *cli.Cmd) {
			cmd.Command("get ls", "Get a list of all plans", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(v.GetAllPlans())
				}
			})

		})
		cmd.Command("plan", "Work with a specific validation plan", func(cmd *cli.Cmd) {
			var p ValidationPlan
			idArg := cmd.StringArg("UUID", "", "UUID of the Validation Plan, Short IDs accepted")

			cmd.Spec = "UUID"
			cmd.Before = func() {
				ok, planID := v.FindPlanID(*idArg)
				if !ok {
					panic(errors.New("could not find the validation plan"))
				}
				p = v.GetPlan(planID)
			}

			cmd.Command("get", "Get information about a single build by its name", func(cmd *cli.Cmd) {

				cmd.Action = func() {
					fmt.Println(p)
				}
			})

		})
	})

}
