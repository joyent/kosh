package cli

import (
	"errors"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func relaysCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()
	}

	cmd.Action = func() { display(conch.GetAllRelays()) }

	cmd.Command("get ls", "Get a list of relays", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(conch.GetAllRelays()) }
	})
}

func relayCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)

	var relay types.Relay
	relayArg := cmd.StringArg(
		"RELAY",
		"",
		"ID of the relay",
	)

	cmd.Spec = "RELAY"

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()

		var e error
		relay, e = conch.GetRelayBySerial(*relayArg)
		if e != nil {
			fatalIf(e)
		}
		if (relay == types.Relay{}) {
			fatalIf(errors.New("relay not found"))
		}
	}
	// default action is to display the relay
	cmd.Action = func() { display(relay, nil) }

	cmd.Command("get", "Get data about a single relay", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(relay, nil) }
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
