package cli

import (
	"errors"
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

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
		conch := config.ConchClient()
		display := config.Renderer()

		validationPlan, e := conch.GetValidationPlanByName(*validationPlanOpt)
		if e != nil {
			fatal(e)
		}

		vendor, e := conch.GetHardwareVendorByName(*vendor)
		if e != nil {
			fatal(e)
		}
		create := types.HardwareProductCreate{
			Name:             types.MojoStandardPlaceholder(*name),
			Alias:            types.MojoStandardPlaceholder(*alias),
			HardwareVendorID: vendor.ID,
			Sku:              types.MojoStandardPlaceholder(*SKU),
			RackUnitSize:     types.PositiveInteger(*rackUnitSize),
			ValidationPlanID: validationPlan.ID,
			Purpose:          *purpose,
			BiosFirmware:     *biosFirmware,
			CPUType:          *cpuType,
		}
		conch.CreateHardwareProduct(create)
		display(conch.GetHardwareProductByID(*name))
	}
}

func cmdListProducts(cmd *cli.Cmd) {
	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()

		display(conch.GetHardwareProducts())
	}
}

func cmdImportProduct(cmd *cli.Cmd) {
	filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file that defines the new hardware product. '-' indicates STDIN")
	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()

		in, err := getInputReader(*filePathArg)
		if err != nil {
			fatal(err)
		}

		p := conch.ReadHardwareProduct(in)
		conch.CreateHardwareProduct(p)
		display(conch.GetHardwareProducts())
	}
}

func hardwareCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()
	}

	cmd.Command("products ps", "Work with hardware products", func(cmd *cli.Cmd) {
		cmd.Command("create", "Create a hardware product", cmdCreateProduct)
		cmd.Command("import", "Import a hardware product as a JSON file", cmdImportProduct)
		cmd.Command("get ls", "Get a list of all hardware products", cmdListProducts)
	})

	cmd.Command("product p", "Work with a hardware product", func(cmd *cli.Cmd) {
		var hp types.HardwareProduct
		idArg := cmd.StringArg("PRODUCT", "", "The SKU, UUID, alias, or name of the hardware product.")
		cmd.Before = func() {
			var e error
			hp, e = conch.GetHardwareProductByID(*idArg)
			if e != nil {
				fatal(e)
			}
			if (hp == types.HardwareProduct{}) {
				fatal(errors.New("Hardware Product not found for " + *idArg))
			}
		}
		cmd.Action = func() { fmt.Println(hp) }
		cmd.Command("get", "Show a hardware vendor's details", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(hp, nil) }
		})
		cmd.Command("delete rm", "Remove a hardware product", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteHardwareProduct(hp.ID)
				display(conch.GetHardwareProducts())
			}
		})
	})

	cmd.Command("vendors vs", "Work with hardware vendors", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(conch.GetAllHardwareVendors()) }

		cmd.Command("get ls", "Get a list of all hardware vendors", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetAllHardwareVendors()) }
		})

		cmd.Command("create", "Create a hardware vendor", func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name of the hardware vendor.")
			cmd.Action = func() {
				conch.FindOrCreateHardwareVendor(*name)
			}
		})
	})

	cmd.Command("vendor v", "Work a specific hardware vendor", func(cmd *cli.Cmd) {
		var hv types.HardwareVendor
		idArg := cmd.StringArg("NAME", "", "The name, or UUID of the hardware vendor.")

		// grab the Vendor for the given ID
		cmd.Before = func() {
			var e error
			hv, e = conch.GetHardwareVendorByName(*idArg)
			if e != nil {
				fatal(e)
			}
			if (hv == types.HardwareVendor{}) {
				fmt.Println("Hardware Vendor not found for " + *idArg)
				cli.Exit(1)
			}
		}

		cmd.Action = func() { display(hv, nil) }
		cmd.Command("get", "Show a hardware vendor's details", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(hv, nil) }
		})
		cmd.Command("delete rm", "Remove a hardware vendor", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteHardwareVendor(hv.ID)
			}
		})
	})
}
