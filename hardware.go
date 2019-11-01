// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	// "bytes"
	// "errors"
	"fmt"
	"net/url"

	// "sort"
	// "strconv"
	// "strings"
	"time"

	"github.com/gofrs/uuid"
	// "github.com/jawher/mow.cli"
	// "github.com/olekukonko/tablewriter"
)

type Hardware struct {
	*Conch
}

func (c *Conch) Hardware() *Hardware {
	return &Hardware{c}
}

type HardwareProductProfile struct {
	ID           uuid.UUID `json:"id" faker:"uuid"`
	BiosFirmware string    `json:"bios_firmware"`
	CpuNum       int       `json:"cpu_num"`
	CpuType      string    `json:"cpu_type"`
	DimmsNum     int       `json:"dimms_num"`
	HbaFirmware  string    `json:"hba_firmware,omitempty"`
	NicsNum      int       `json:"nics_num"`
	Purpose      string    `json:"purpose"`
	RamTotal     int       `json:"ram_total"`
	SasHddSlots  string    `json:"sas_hdd_slots,omitempty"`
	SataHddSlots string    `json:"sata_hdd_slots,omitempty"`
	SataSsdSlots string    `json:"sata_ssd_slots,omitempty"`
	SasSsdSlots  string    `json:"sas_ssd_slots,omitempty"`
	NvmeSsdSlots string    `json:"nvme_ssd_slots,omitempty"`
	UsbNum       int       `json:"usb_num"`

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
}

type HardwareProducts []HardwareProduct
type HardwareProduct struct {
	ID                     uuid.UUID              `json:"id" faker:"uuid"`
	Name                   string                 `json:"name"`
	Alias                  string                 `json:"alias"`
	Prefix                 string                 `json:"prefix,omitempty"`
	HardwareVendorID       uuid.UUID              `json:"hardware_vendor_id" faker:"uuid"`
	GenerationName         string                 `json:"generation_name,omitempty"`
	LegacyProductName      string                 `json:"legacy_product_name,omitempty"`
	SKU                    string                 `json:"sku,omitempty"`
	Specification          string                 `json:"specification,omitempty"`
	RackUnitSize           int                    `json:"rack_unit_size" faker:"rack_unit_size"`
	HardwareProductProfile HardwareProductProfile `json:"hardware_product_profile,omitempty"`
	Created                time.Time              `json:"created" faker:"-"`
	Updated                time.Time              `json:"updated" faker:"-"`
	ValidationPlanID       uuid.UUID              `json:"validation_plan_id,omitempty" faker:"-"`
}

func (h *Hardware) GetProduct(id uuid.UUID) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/%s", url.PathEscape(id.String()))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

func (h *Hardware) GetProductByName(name string) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/name=%s", url.PathEscape(name))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

func (h *Hardware) GetProductByAlias(alias string) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/alias=%s", url.PathEscape(alias))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

func (h *Hardware) GetProductBySku(sku string) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/sku=%s", url.PathEscape(sku))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

type HardwareProductProfileUpdate struct {
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
}

func (h *Hardware) Create(
	name, alias string,
	vendorID uuid.UUID,
	SKU string,
	rackUnitSize int,
	validationPlanID uuid.UUID,
	hardwareProductProfile HardwareProductProfileUpdate,
) (hp HardwareProduct) {
	payload := make(map[string]interface{})
	payload["name"] = name
	payload["alias"] = alias
	payload["hardware_vendor_id"] = vendorID
	payload["sku"] = SKU
	payload["rack_unit_size"] = rackUnitSize
	payload["validation_plan_id"] = validationPlanID
	payload["hardware_product_profile"] = hardwareProductProfile

	res := h.Do(h.Sling().New().Post("/hardware_product").
		Set("Content-Type", "application/json").
		BodyJSON(payload),
	)

	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return
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

func (h *Hardware) DeleteVendor(name string) {
	uri := fmt.Sprintf("/hardware_vendor/%s", url.PathEscape(name))
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
