package cli

import (
	"errors"
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func datacentersCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{})

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()
	}
	cmd.Command("get", "Get a list of all datacenters", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			display(conch.GetAllDatacenters())
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
				fmt.Println("--vendor is required")
				cli.Exit(1)
			}
			if *regionOpt == "" {
				fmt.Println("--region is required")
				cli.Exit(1)
			}
			if *locationOpt == "" {
				fmt.Println("--location is required")
				cli.Exit(1)
			}

			conch.CreateDatacenter(types.DatacenterCreate{
				Location:   types.NonEmptyString(*locationOpt),
				Region:     types.NonEmptyString(*regionOpt),
				Vendor:     types.NonEmptyString(*vendorOpt),
				VendorName: types.NonEmptyString(*vendorNameOpt),
			})
		}
	})
}

func datacenterCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{})
	var dc types.Datacenter

	idArg := cmd.StringArg(
		"UUID",
		"",
		"The UUID of the datacenter. Short UUIDs (first segment) accepted",
	)
	cmd.Spec = "UUID"

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()

		dc = conch.GetDatacenterByName(*idArg)
		if (dc == types.Datacenter{}) {
			fatal(errors.New("couldn't find datacenter"))
		}
	}

	cmd.Command("get", "Information about a single datacenter", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(dc) }
	})

	cmd.Command("delete", "Delete a single datacenter", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch.DeleteDatacenter(dc.ID)
			display(conch.GetAllDatacenters())
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
				fatal(errors.New("one option must be provided"))
			}
			conch.UpdateDatacenter(dc.ID, types.DatacenterUpdate{
				Location:   types.NonEmptyString(*locationOpt),
				Region:     types.NonEmptyString(*regionOpt),
				Vendor:     types.NonEmptyString(*vendorOpt),
				VendorName: types.NonEmptyString(*vendorNameOpt),
			})
		}
	})

	cmd.Command("rooms", "Get the room list for a single datacenter", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			display(conch.GetAllDatacenterRooms(dc.ID))
		}
	})
}
