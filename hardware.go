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
	ID           uuid.UUID `json:"id"`
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
	ID                     uuid.UUID              `json:"id"`
	Name                   string                 `json:"name"`
	Alias                  string                 `json:"alias"`
	Prefix                 string                 `json:"prefix,omitempty"`
	HardwareVendorID       uuid.UUID              `json:"hardware_vendor_id"`
	GenerationName         string                 `json:"generate_name,omitempty"`
	LegacyProductName      string                 `json:"legacy_product_name,omitempty"`
	SKU                    string                 `json:"sku,omitempty"`
	Specification          string                 `json:"specification,omitempty"`
	RackUnitSize           int                    `json:"rack_unit_size"`
	HardwareProductProfile HardwareProductProfile `json:"hardware_product_profile,omitempty"`
	Created                time.Time              `json:"created"`
	Updated                time.Time              `json:"updated"`
}

func (h *Hardware) GetProduct(id uuid.UUID) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/%s", url.PathEscape(id.String()))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

// There are three string identifiers currently accepted by the API: name,
// alias, sku. The calls all look exactly the same where we stick the string in
// the url and hope for the best.
func (h *Hardware) GetProductByString(wat string) (hp HardwareProduct) {
	uri := fmt.Sprintf("/hardware_product/%s", url.PathEscape(wat))
	res := h.Do(h.Sling().New().Get(uri))
	if ok := res.Parse(&hp); !ok {
		panic(res)
	}

	return hp
}

func (h *Hardware) GetProductByName(name string) HardwareProduct {
	return h.GetProductByString(name)
}

func (h *Hardware) GetProductByAlias(alias string) HardwareProduct {
	return h.GetProductByString(alias)
}

func (h *Hardware) GetProductBySku(sku string) HardwareProduct {
	return h.GetProductByString(sku)
}
