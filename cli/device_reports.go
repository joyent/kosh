package cli

import (
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
)

func deviceReportCmd(cmd *cli.Cmd) {
	cmd.Command("post", "Post a new device report", func(cmd *cli.Cmd) {
		var conch *conch.Client

		filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file that defines the layout. '-' indicates STDIN")

		input, err := getInputReader(*filePathArg)
		if err != nil {
			fatal(err)
		}

		cmd.Before = func() { conch = config.ConchClient() }
		cmd.Action = func() { conch.SendDeviceReport(input) }
	})
}
