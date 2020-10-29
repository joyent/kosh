package cli

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
)

func schemaCmd(cfg Config) func(*cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		cmd.Before = func() {
			conch = cfg.ConchClient()
		}

		cmd.Command("request", "View your Conch profile", func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The string name of a request schema")
			cmd.Spec = "NAME"

			cmd.Action = func() {
				display(conch.GetSchema(fmt.Sprintf("request/%s", *name)))
			}
		})

		cmd.Command("response", "Get the settings for the current user", func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The string name of a response schema")
			cmd.Spec = "NAME"

			cmd.Action = func() {
				display(conch.GetSchema(fmt.Sprintf("response/%s", *name)))
			}
		})
	}
}
