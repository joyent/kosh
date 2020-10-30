package cli

import (
	"errors"
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch/types"
)

func relaysCmd(cfg Config) func(cmd *cli.Cmd) {
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := cfg.ConchClient()
			display(conch.GetAllRelays())
		}
		cmd.Command("get ls", "Get a list of relays", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch := cfg.ConchClient()
				display(conch.GetAllRelays())
			}
		})
	}
}

func relayCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		var relay types.Relay
		relayArg := cmd.StringArg(
			"RELAY",
			"",
			"ID of the relay",
		)

		cmd.Spec = "RELAY"

		cmd.Before = func() {
			conch = cfg.ConchClient()
			relay = conch.GetRelayBySerial(*relayArg)
		}
		cmd.Command("get", "Get data about a single relay", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if (relay == types.Relay{}) {
					fatal(errors.New("relay not found"))
				}
				fmt.Println(relay)
			}
		})

		cmd.Command("register", "Register a relay with the API", func(cmd *cli.Cmd) {
			var (
				versionOpt = cmd.StringOpt("version", "", "The version of the relay")
				sshPortOpt = cmd.IntOpt("ssh_port port", 22, "The SSH port for the relay")
				ipAddrOpt  = cmd.StringOpt("ipaddr ip", "", "The IP address for the relay")
				nameOpt    = cmd.StringOpt("name", "", "The name of the relay")
			)

			cmd.Action = func() {
				conch.RegisterRelay(*relayArg, types.RegisterRelay{
					Version: *versionOpt,
					Ipaddr:  *ipAddrOpt,
					Name:    types.NonEmptyString(*nameOpt),
					SSHPort: types.NonNegativeInteger(*sshPortOpt),
				})
			}
		})

		cmd.Command("delete rm", "Delete a relay", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteRelay(relay.ID.String())
				display(conch.GetAllRelays())
			}
		})
	}
}
