package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
)

func devicesCmd(cmd *cli.Cmd) {
	cmd.Before = config.Before(requireAuth)
	cmd.Command("search s", "Search for devices", deviceSearchCmd)
}

func deviceSearchCmd(cmd *cli.Cmd) {
	cmd.Before = config.Before(requireAuth)
	cmd.Command("setting", "Search for devices by exact setting value", searchBySettingCmd)
	cmd.Command("tag", "Search for devices by exact tag value", searchByTagCmd)
	cmd.Command("hostname", "Search for devices by exact hostname", searchByHostnameCmd)
}

func searchBySettingCmd(cmd *cli.Cmd) {
	key := *cmd.StringArg("KEY", "", "Setting name")
	value := *cmd.StringArg("VALUE", "", "Setting Value")
	cmd.Spec = "KEY VALUE"

	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()

		display(conch.FindDevicesBySetting(key, value))
	}
}

func searchByTagCmd(cmd *cli.Cmd) {
	key := *cmd.StringArg("KEY", "", "Tag name")
	value := *cmd.StringArg("VALUE", "", "Tag Value")
	cmd.Spec = "KEY VALUE"

	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()

		display(conch.FindDevicesByTag(key, value))
	}
}

func searchByHostnameCmd(cmd *cli.Cmd) {
	hostname := *cmd.StringArg("HOSTNAME", "", "hostname")
	cmd.Spec = "HOSTNAME"

	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()

		display(conch.FindDevicesByField("hostname", hostname))
	}
}

// Single Device Commands
func deviceCmd(cmd *cli.Cmd) {
	id := cmd.StringArg(
		"DEVICE",
		"",
		"UUID or serial number of the device. Short UUIDs are *not* accepted",
	)
	cmd.Spec = "DEVICE"

	cmd.Command("get", "Get information about a single device", deviceGetCmd(id))
	cmd.Command("validations", "Get the most recent validation results for a single device", deviceValidationsCmd(id))
	cmd.Command("settings", "See all settings for a device", deviceSettingsCmd(id))
	cmd.Command("setting", "See a single setting for a device", deviceSettingCmd(id))
	cmd.Command("tags", "See all tags for a device", deviceTagsCmd(id))
	cmd.Command("tag", "See a single tag for a device", deviceTagCmd(id))
	cmd.Command("interface", "Information about a single interface", deviceInterfaceCmd(id))
	cmd.Command("preflight", "Data that is only accurate inside preflight", devicePreflightCmd(id))
	cmd.Command("phase", "Actions on the lifecycle phase of the device", devicePhaseCmd(id))
	cmd.Command("report", "Get the most recently recorded report for this device", deviceDeviceReportCmd(id))
}

func deviceGetCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()

			display(conch.GetDeviceBySerial(*id))
		}
	}
}

func deviceValidationsCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()

			display(conch.GetDeviceValidationStates(*id))
		}
	}
}

func deviceSettingsCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()

			display(conch.GetDeviceSettings(*id))
		}
	}
}

func deviceSettingCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		var conch *conch.Client
		var display func(interface{}, error)

		key := *cmd.StringArg(
			"NAME",
			"",
			"Name of the setting",
		)
		cmd.Spec = "NAME"

		cmd.Before = func() {
			conch = config.ConchClient()
			display = config.Renderer()
		}
		cmd.Action = func() {
			display(conch.GetDeviceSettingByName(*id, key))
		}

		cmd.Command("get", "Get a particular device setting", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetDeviceSettingByName(*id, key)) }
		})

		cmd.Command("set", "Set a particular device setting", func(cmd *cli.Cmd) {
			value := *cmd.StringArg("VALUE", "", "Value of the setting")
			cmd.Spec = "VALUE"

			cmd.Action = func() {
				conch.SetDeviceSetting(*id, key, value)
				display(conch.GetDeviceSettings(*id))
			}
		})

		cmd.Command("delete rm", "Delete a particular device setting", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteDeviceSetting(*id, key)
				display(conch.GetDeviceSettings(*id))
			}
		})
	}
}

func deviceTagsCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()

			display(conch.GetDeviceTags(*id))
		}
	}
}

func deviceTagCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		var conch *conch.Client
		var display func(interface{}, error)

		name := *cmd.StringArg("NAME", "", "Name of the tag")
		cmd.Spec = "NAME"

		cmd.Before = func() {
			conch = config.ConchClient()
			display = config.Renderer()
		}

		cmd.Action = func() { display(conch.GetDeviceTagByName(*id, name)) }

		cmd.Command("get", "Get a particular device tag", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetDeviceTagByName(*id, name)) }
		})

		cmd.Command("set", "Set a particular device tag", func(cmd *cli.Cmd) {
			value := *cmd.StringArg("VALUE", "", "Value of the tag")
			cmd.Spec = "VALUE"

			cmd.Action = func() {
				conch.SetDeviceTag(*id, name, value)
				display(conch.GetDeviceTags(*id))
			}
		})

		cmd.Command("delete rm", "Delete a particular device tag", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteDeviceTag(*id, name)
				display(conch.GetDeviceTags(*id))
			}
		})
	}
}

func deviceInterfaceCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		name := *cmd.StringArg("NAME", "", "Name of the interface")
		cmd.Spec = "NAME"
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()

			display(conch.GetDeviceInterfaceByName(*id, name))
		}
	}
}

func devicePreflightCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		var conch *conch.Client
		var display func(interface{}, error)

		cmd.Before = func() {
			conch = config.ConchClient()
			display = config.Renderer()
			phase, e := conch.GetDevicePhase(*id)
			if e != nil {
				fatal(e)
			}
			if phase != "integration" {
				os.Stderr.WriteString("Warning: This device is no longer in the 'integration' phase. This data is likely to be inaccurate\n")
			}
		}

		cmd.Command("location", "The location of a device in preflight", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetDeviceLocation(*id)) }
		})

		cmd.Command("ipmi", "IPMI address for a device in preflight", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				iface, e := conch.GetDeviceInterfaceByName(*id, "ipmi1")
				if e != nil {
					fatal(e)
				}
				fmt.Println(iface.Ipaddr)
			}
		})
	}
}

/***/
var phasesList = []string{
	"integration",
	"installation",
	"production",
	"diagnostics",
	"decommissioned",
}

func prettyPhasesList() string {
	return strings.Join(phasesList, ", ")
}

func okPhase(phase string) bool {
	for _, b := range phasesList {
		if phase == b {
			return true
		}
	}
	return false
}

func devicePhaseCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		var conch *conch.Client
		var display func(interface{}, error)

		cmd.Before = func() {
			conch = config.ConchClient()
			display = config.Renderer()
		}

		cmd.Command("get", "Get the phase of the device", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetDevicePhase(*id)) }
		})

		cmd.Command("set", "Set the phase of the device [one of: "+prettyPhasesList()+"]", func(cmd *cli.Cmd) {
			phase := *cmd.StringArg("PHASE", "", "Name of the phase [one of: "+prettyPhasesList()+"]")
			cmd.Spec = "PHASE"
			cmd.Action = func() {
				if !okPhase(phase) {
					log.Fatal("Phase must be one of: " + prettyPhasesList())
				}
				conch.SetDevicePhase(*id, phase)
				display(conch.GetDevicePhase(*id))
			}
		})
	}
}

func deviceDeviceReportCmd(id *string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()
			d, e := conch.GetDeviceBySerial(*id)
			display(d.LatestReport, e)
		}
	}
}
