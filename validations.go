// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	// "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
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

// type ValidationPlans []ValidationPlan

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
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)
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

type ValidationStatesWithResults []ValidationStateWithResults

func (v ValidationStatesWithResults) Len() int {
	return len(v)
}

func (v ValidationStatesWithResults) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ValidationStatesWithResults) Less(i, j int) bool {
	return v[i].Created.Before(v[j].Created)
}

func (vs ValidationStatesWithResults) String() string {
	sort.Sort(vs)
	if API.JsonOnly {
		return API.AsJSON(vs)
	}

	type extendedVsR struct {
		ValidationStateWithResults
		ValidationPlan ValidationPlan `json:"-"`
	}

	out := make([]extendedVsR, 0)
	for _, v := range vs {
		out = append(out, extendedVsR{
			v,
			API.Validations().GetPlan(v.ValidationPlanID),
		})
	}

	t, err := NewTemplate().Parse(validationStatesWithResultsTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, out); err != nil {
		panic(err)
	}

	return buf.String()
}
