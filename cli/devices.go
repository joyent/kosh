package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	cli "github.com/jawher/mow.cli"
)

type Renderable interface {
	JSON() ([]byte, error)
	String() string
}

func display(cfg Config, item Renderable) {
	if cfg.OutputJSON {
		fmt.Println(item.JSON())
	} else {
		fmt.Println(item.String())
	}
}

func devicesCmd(cfg Config) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("search s", "Search for devices", deviceSearchCmd(cfg))
	}
}

func deviceSearchCmd(cfg Config) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("setting", "Search for devices by exact setting value", searchBySettingCmd(cfg))
		cmd.Command("tag", "Search for devices by exact tag value", searchByTagCmd(cfg))
		cmd.Command("hostname", "Search for devices by exact hostname", searchByHostnameCmd(cfg))
	}
}

func searchBySettingCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		key := *cmd.StringArg("KEY", "", "Setting name")
		value := *cmd.StringArg("VALUE", "", "Setting Value")
		cmd.Spec = "KEY VALUE"

		cmd.Action = func() {
			display(cfg, conch.FindDevicesBySetting(key, value))
		}
	}
}

func searchByTagCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		key := *cmd.StringArg("KEY", "", "Tag name")
		value := *cmd.StringArg("VALUE", "", "Tag Value")
		cmd.Spec = "KEY VALUE"

		cmd.Action = func() {
			display(cfg, conch.FindDevicesByTag(key, value))
		}
	}
}

func searchByHostnameCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		hostname := *cmd.StringArg("HOSTNAME", "", "hostname")
		cmd.Spec = "HOSTNAME"

		cmd.Action = func() {
			display(cfg, conch.FindDevicesByField("hostname", hostname))
		}
	}
}

// Single Device Commands
func deviceCmd(cfg Config) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		id := *cmd.StringArg(
			"DEVICE",
			"",
			"UUID or serial number of the device. Short UUIDs are *not* accepted",
		)

		cmd.Spec = "DEVICE"

		cmd.Command("get", "Get information about a single device", deviceGetCmd(cfg, id))
		cmd.Command("validations", "Get the most recent validation results for a single device", deviceValidationsCmd(cfg, id))
		cmd.Command("settings", "See all settings for a device", deviceSettingsCmd(cfg, id))
		cmd.Command("setting", "See a single setting for a device", deviceSettingCmd(cfg, id))
		cmd.Command("tags", "See all tags for a device", deviceTagsCmd(cfg, id))
		cmd.Command("tag", "See a single tag for a device", deviceTagCmd(cfg, id))
		cmd.Command("interface", "Information about a single interface", deviceInterfaceCmd(cfg, id))
		cmd.Command("preflight", "Data that is only accurate inside preflight", devicePreflightCmd(cfg, id))
		cmd.Command("phase", "Actions on the lifecycle phase of the device", devicePhaseCmd(cfg, id))

		cmd.Command("report", "Get the most recently recorded report for this device", deviceReportCmd(cfg, id))
	}
}

func deviceGetCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() { display(cfg, conch.GetDeviceById(id)) }
	}
}

func deviceValidationsCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() { display(cfg, conch.GetDeviceValidationStates(id)) }
	}
}

func deviceSettingsCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() { display(cfg, conch.GetDeviceSettings(id)) }
	}
}

func deviceSettingCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		key := *cmd.StringArg(
			"NAME",
			"",
			"Name of the setting",
		)

		cmd.Spec = "NAME"

		cmd.Action = func() {
			display(cfg, conch.GetDeviceSettingByName(id, key))
		}

		cmd.Command("get", "Get a particular device setting", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(cfg, conch.GetDeviceSettingByName(id, key)) }
		})

		cmd.Command("set", "Set a particular device setting", func(cmd *cli.Cmd) {
			value := *cmd.StringArg("VALUE", "", "Value of the setting")
			cmd.Spec = "VALUE"

			cmd.Action = func() {
				conch.SetDeviceSetting(id, key, value)
				display(cfg, conch.GetDeviceSettings(id))
			}
		})

		cmd.Command("delete rm", "Delete a particular device setting", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteDeviceSetting(id, key)
				display(cfg, conch.GetDeviceSettings(id))
			}
		})
	}
}

func deviceTagsCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() { display(cfg, conch.GetDeviceTags(id)) }
	}
}

func deviceTagCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		name := *cmd.StringArg("NAME", "", "Name of the tag")

		cmd.Spec = "NAME"

		cmd.Action = func() { display(cfg, conch.GetDeviceTagByName(id, name)) }

		cmd.Command("get", "Get a particular device tag", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(cfg, conch.GetDeviceTagByName(id, name)) }
		})

		cmd.Command("set", "Set a particular device tag", func(cmd *cli.Cmd) {
			value := *cmd.StringArg("VALUE", "", "Value of the tag")
			cmd.Spec = "VALUE"

			cmd.Action = func() {
				conch.SetDeviceTag(id, name, value)
				display(cfg, conch.GetDeviceTags(id))
			}
		})

		cmd.Command("delete rm", "Delete a particular device tag", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteDeviceTag(id, name)
				display(cfg, conch.GetDeviceTags(id))
			}
		})
	}
}

func deviceInterfaceCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		name := *cmd.StringArg("NAME", "", "Name of the interface")
		cmd.Spec = "NAME"
		cmd.Action = func() { display(cfg, conch.GetDeviceInterfaceByName(id, name)) }
	}
}

func devicePreflightCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Before = func() {
			if conch.GetDevicePhase(id) != "integration" {
				os.Stderr.WriteString("Warning: This device is no longer in the 'integration' phase. This data is likely to be inaccurate\n")
			}
		}

		cmd.Command("location", "The location of a device in preflight", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(cfg, conch.GetDeviceLocation(id)) }
		})

		cmd.Command("ipmi", "IPMI address for a device in preflight", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				iface := conch.GetDeviceInterfaceByName(id, "ipmi1")
				fmt.Println(iface.Ipaddr)
			}
		})
	}
}

/***/
var PhasesList = []string{
	"integration",
	"installation",
	"production",
	"diagnostics",
	"decommissioned",
}

func prettyPhasesList() string {
	return strings.Join(PhasesList, ", ")
}

func okPhase(phase string) bool {
	for _, b := range PhasesList {
		if phase == b {
			return true
		}
	}
	return false
}

func devicePhaseCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Command("get", "Get the phase of the device", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(cfg, conch.GetDevicePhase(id)) }
		})

		cmd.Command("set", "Set the phase of the device [one of: "+prettyPhasesList()+"]", func(cmd *cli.Cmd) {
			phase := *cmd.StringArg("PHASE", "", "Name of the phase [one of: "+prettyPhasesList()+"]")
			cmd.Spec = "PHASE"
			cmd.Action = func() {
				if !okPhase(phase) {
					log.Fatal("Phase must be one of: " + prettyPhasesList())
				}
				conch.SetDevicePhase(id, phase)
				display(cfg, conch.GetDevicePhase(id))
			}
		})
	}
}

func deviceReportCmd(cfg Config, id string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() { display(cfg, conch.GetDeviceById(id).LatestReport) }
	}
}
