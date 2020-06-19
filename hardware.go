// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
	// "github.com/olekukonko/tablewriter"
)

type Hardware struct {
	*Conch
}

func (c *Conch) Hardware() *Hardware {
	return &Hardware{c}
}

type HardwareProductSummaries []HardwareProductSummary

func (hps HardwareProductSummaries) String() string {
	if API.JsonOnly {
		return API.AsJSON(hps)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"ID",
		"SKU",
		"Name",
		"Alias",
		"GenerationName",
		"Created",
		"Updated",
	})

	for _, hp := range hps {
		table.Append([]string{
			CutUUID(hp.ID.String()),
			hp.SKU,
			hp.Name,
			hp.Alias,
			hp.GenerationName,
			hp.Created.String(),
			hp.Updated.String(),
		})
	}
	table.Render()
	return tableString.String()
}

func (hp HardwareProduct) HardwareVendor() HardwareVendor {
	return API.Hardware().GetVendor(hp.HardwareVendorID.String())
}

func (hp HardwareProduct) ValidationPlan() ValidationPlan {
	return API.Validations().GetPlan(hp.ValidationPlanID)
}

type HardwareProductSummary struct {
	ID                uuid.UUID `json:"id" faker:"uuid"`
	Name              string    `json:"name"`
	Alias             string    `json:"alias"`
	GenerationName    string    `json:"generation_name,omitempty"`
	SKU               string    `json:"sku"`
	Created time.Time `json:"created" faker:"-"`
	Updated time.Time `json:"updated" faker:"-"`
}

type HardwareProduct struct {
	ID                uuid.UUID `json:"id" faker:"uuid"`
	Name              string    `json:"name"`
	Alias             string    `json:"alias"`
	Prefix            string    `json:"prefix,omitempty"`
	GenerationName    string    `json:"generation_name,omitempty"`
	LegacyProductName string    `json:"legacy_product_name,omitempty"`
	SKU               string    `json:"sku"`
	Specification     string    `json:"specification,omitempty"`
	RackUnitSize      int       `json:"rack_unit_size" faker:"rack_unit_size"`

	HardwareVendorID uuid.UUID `json:"hardware_vendor_id" faker:"uuid"`
	ValidationPlanID uuid.UUID `json:"validation_plan_id,omitempty" faker:"-"`

	BiosFirmware string `json:"bios_firmware"`
	CpuNum       int    `json:"cpu_num"`
	CpuType      string `json:"cpu_type"`
	DimmsNum     int    `json:"dimms_num"`
	HbaFirmware  string `json:"hba_firmware,omitempty"`
	NicsNum      int    `json:"nics_num"`
	Purpose      string `json:"purpose"`
	RamTotal     int    `json:"ram_total"`
	SasHddSlots  string `json:"sas_hdd_slots,omitempty"`
	SataHddSlots string `json:"sata_hdd_slots,omitempty"`
	SataSsdSlots string `json:"sata_ssd_slots,omitempty"`
	SasSsdSlots  string `json:"sas_ssd_slots,omitempty"`
	NvmeSsdSlots string `json:"nvme_ssd_slots,omitempty"`
	UsbNum       int    `json:"usb_num"`

	// NOTE the pointers. 0 is a valid value so zero values aren't
	PsuTotal   *int `json:"psu_total,omitempty"`
	RaidLunNum *int `json:"raid_lun_num,omitempty"`

	SasHddNum  *int `json:"sas_hdd_num,omitempty"`
	SasHddSize *int `json:"sas_hdd_size,omitempty"`

	SataHddNum  *int `json:"sata_hdd_num,omitempty"`
	SataHddSize *int `json:"sata_hdd_size,omitempty"`

	SataSsdNum  *int `json:"sata_ssd_num,omitempty"`
	SataSsdSize *int `json:"sata_ssd_size,omitempty"`

	SasSsdNum  *int `json:"sas_ssd_num,omitempty"`
	SasSsdSize *int `json:"sas_ssd_size,omitempty"`

	NvmeSsdNum  *int `json:"nvme_ssd_num,omitempty"`
	NvmeSsdSize *int `json:"nvme_ssd_size,omitempty"`

	Created time.Time `json:"created" faker:"-"`
	Updated time.Time `json:"updated" faker:"-"`
}

func (hp HardwareProduct) String() string {
	if API.JsonOnly {
		return API.AsJSON(hp)
	}
	t, err := NewTemplate().Parse(hardwareProductTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, hp); err != nil {
		panic(err)
	}

	return buf.String()
}

func (h *Hardware) GetAllProducts() (hps HardwareProductSummaries) {
	res := h.Do(h.Sling().New().Get("/hardware_product"))
	if ok := res.Parse(&hps); !ok {
		panic(res)
	}
	return hps
}

func (h *Hardware) GetProduct(id uuid.UUID) (hp HardwareProduct) {
	return h.GetProductByString(id.String())
}

func (h *Hardware) GetProductByString(s string) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/%s", url.PathEscape(s))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

func (h Hardware) GetProductByName(name string) HardwareProduct {
	return h.GetProductByString(name)
}

func (h Hardware) GetProductBySku(sku string) HardwareProduct {
	return h.GetProductByString(sku)
}

func (h Hardware) GetProductByAlias(alias string) HardwareProduct {
	return h.GetProductByString(alias)
}

func (h *Hardware) Create(
	name, alias string,
	vendorID uuid.UUID,
	SKU string,
	rackUnitSize int,
	validationPlanID uuid.UUID,
	Purpose string,
	BiosFirmware string,
	CpuType string,
) (hp HardwareProduct) {
	payload := make(map[string]interface{})
	payload["name"] = name
	payload["alias"] = alias
	payload["hardware_vendor_id"] = vendorID
	payload["sku"] = SKU
	payload["rack_unit_size"] = rackUnitSize
	payload["validation_plan_id"] = validationPlanID
	payload["purpose"] = Purpose
	payload["bios_firmware"] = BiosFirmware
	payload["cpu_type"] = CpuType

	res := h.Do(h.Sling().New().Post("/hardware_product").
		Set("Content-Type", "application/json").
		BodyJSON(payload),
	)

	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

func (h *Hardware) Delete(ID uuid.UUID) {
	uri := fmt.Sprintf("/hardware_product/%s", url.PathEscape(ID.String()))
	res := h.Do(h.Sling().New().Delete(uri))

	if res.StatusCode() != 204 {
		// I know this is weird. Like in other places, it should be impossible
		// to reach here unless the status code is 204. The API returns 204
		// (which gets us here) or 409 (which will explode before it gets here).
		// If we got here via some other code, then there's some new behavior
		// that we need to know about.
		panic(res)
	}

}

type HardwareVendor struct {
	ID      uuid.UUID `json:"id" faker:"uuid"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (h *Hardware) GetAllVendors() (hvs []HardwareVendor) {
	res := h.Do(h.Sling().Get("/hardware_vendor/"))
	if ok := res.Parse(&hvs); !ok {
		panic(res)
	}

	return
}

func (h *Hardware) GetVendor(name string) (hv HardwareVendor) {
	uri := fmt.Sprintf("/hardware_vendor/%s", url.PathEscape(name))

	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hv); !ok {
		panic(res)
	}

	return
}

func (h *Hardware) CreateVendor(name string) (hv HardwareVendor) {
	uri := fmt.Sprintf("/hardware_vendor/%s", url.PathEscape(name))

	_ = h.Do(h.Sling().New().Post(uri))

	return h.GetVendor(name)
}

func (h *Hardware) FindOrCreateVendor(name string) (hv HardwareVendor) {
	// if we fail to Find the vendor we panic,
	// really we  need to calm down and just create a new Vendor
	defer func() {
		if r := recover(); r != nil {
			hv = h.CreateVendor(name)
		}
	}()
	hv = h.GetVendor(name)
	return
}

func (h *Hardware) DeleteVendor(id uuid.UUID) {
	uri := fmt.Sprintf("/hardware_vendor/%s", url.PathEscape(id.String()))
	res := h.Do(h.Sling().New().Delete(uri))

	if res.StatusCode() != 204 {
		// I know this is weird. Like in other places, it should be impossible
		// to reach here unless the status code is 204. The API returns 204
		// (which gets us here) or 409 (which will explode before it gets here).
		// If we got here via some other code, then there's some new behavior
		// that we need to know about.
		panic(res)
	}

}

/*
	kosh hardware products create \
		--name Foo \
		--alias Bar \
		--vendor MyBigVendor \
		--SKU FooBle-001 \
		--rack-unit-size 4 \
		--validation-plan ServerPlan \
		--purpose Server \
		--biosFirmware firmware23.dat \
		--cpuType Intel
*/

func cmdCreateProduct(cmd *cli.Cmd) {
	var (
		name              = cmd.StringOpt("name", "", "Name of the hardware product")
		alias             = cmd.StringOpt("alias", "", "Alias for the hardware product")
		vendor            = cmd.StringOpt("vendor", "", "Vendor of the hardware product")
		SKU               = cmd.StringOpt("sku", "", "SKU for the hardware product")
		rackUnitSize      = cmd.IntOpt("rack-unit-size", 2, "RU size of the hardware product")
		validationPlanOpt = cmd.StringOpt("validation-plan", "", "Name of the plan to validate the product against")
		purpose           = cmd.StringOpt("purpose", "", "Purpose of the product")
		biosFirmware      = cmd.StringOpt("bios-firmware", "", "BIOS firmware version for the product")
		cpuType           = cmd.StringOpt("cpu-type", "", "CPU type for the product")
	)

	cmd.Spec = "--sku --name --alias --vendor --validation-plan --purpose --bios-firmware --cpu-type [OPTIONS]"
	cmd.Action = func() {
		validationPlan := API.Validations().GetPlanByName(*validationPlanOpt)
		vendor := API.Hardware().GetVendor(*vendor)
		fmt.Println(API.Hardware().Create(
			*name,
			*alias,
			vendor.ID,
			*SKU,
			*rackUnitSize,
			validationPlan.ID,
			*purpose,
			*biosFirmware,
			*cpuType,
		))
	}
}

// kosh hardware products get
func cmdListProducts(cmd *cli.Cmd) {
	cmd.Action = func() {
		fmt.Println(API.Hardware().GetAllProducts())
	}
}

func initHardwareCli(cmd *cli.Cli) {

	cmd.Command("hardware", "Work with hardware profiles and vendors", func(cmd *cli.Cmd) {
		cmd.Command("products", "Work with hardware products", func(cmd *cli.Cmd) {
			cmd.Command("create", "Create a hardware product", cmdCreateProduct)
			cmd.Command("get ls", "Get a list of all hardware products", cmdListProducts)
		})

		cmd.Command("product", "Work with a hardware product", func(cmd *cli.Cmd) {
			var hp HardwareProduct
			idArg := cmd.StringArg("PRODUCT", "", "The SKU, UUID, alias, or name of the hardware product.")
			cmd.Before = func() {
				hp = API.Hardware().GetProductByString(*idArg)
				if (hp == HardwareProduct{}) {
					panic("Hardware Product not found for " + *idArg)
				}
			}
			cmd.Command("get", "Show a hardware vendor's details", func(cmd *cli.Cmd) {
				cmd.Action = func() { fmt.Println(hp) }
			})
			cmd.Command("delete rm", "Remove a hardware product", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					API.Hardware().Delete(hp.ID)
					fmt.Println(API.Hardware().GetAllProducts())
				}
			})
		})

		cmd.Command("vendors", "Work with hardware vendors", func(cmd *cli.Cmd) {
			cmd.Command("get ls", "Get a list of all hardware vendors", func(cmd *cli.Cmd) {
				API.Hardware().GetAllVendors()
			})
			cmd.Command("create", "Create a hardware vendor", func(cmd *cli.Cmd) {
				name := cmd.StringArg("NAME", "", "The name of the hardware vendor.")
				cmd.Action = func() { API.Hardware().FindOrCreateVendor(*name) }
			})
		})

		cmd.Command("vendor", "Work a specific hardware vendor", func(cmd *cli.Cmd) {
			var hv HardwareVendor
			idArg := cmd.StringArg("NAME", "", "The name, or UUID of the hardware vendor.")

			// grab the Vendor for the given ID
			cmd.Before = func() {
				hv = API.Hardware().GetVendor(*idArg)
				if (hv == HardwareVendor{}) {
					panic("Hardware Vendor not found for " + *idArg)
				}
			}

			cmd.Command("get", "Show a hardware vendor's details", func(cmd *cli.Cmd) {
				cmd.Action = func() { fmt.Println(hv) }
			})
			cmd.Command("delete rm", "Remove a hardware vendor", func(cmd *cli.Cmd) {
				cmd.Action = func() { API.Hardware().DeleteVendor(hv.ID) }
			})
		})
	})
}

func init() { initHardwareCli(App) }
