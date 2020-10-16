package cli

import (
	cli "github.com/jawher/mow.cli"
)

func deviceReportCmd(cfg Config) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("post", "Post a new device report", deviceReportPostCmd(cfg))
	}
}

func deviceReportPostCmd(cfg Config) func(*cli.Cmd) {
	conch := cfg.ConchClient()

	return func(cmd *cli.Cmd) {
		filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file that defines the layout. '-' indicates STDIN")

		input, err := getInputReader(*filePathArg)
		if err != nil {
			fatal(err)
		}

		cmd.Before = func() { conch = cfg.ConchClient() }
		cmd.Action = func() { conch.SendDeviceReport(input) }
	}
}
