package cli

import (
	"errors"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch/types"
)

func validationCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		cmd.Before = func() {
			conch = cfg.ConchClient()
		}
		cmd.Command("plans", "Work with validation plans", func(cmd *cli.Cmd) {
			cmd.Command("get ls", "Get a list of all plans", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					display(conch.GetAllValidationPlans())
				}
			})
		})
		cmd.Command("plan", "Work with a specific validation plan", func(cmd *cli.Cmd) {
			var plan types.ValidationPlan

			idArg := cmd.StringArg("UUID", "", "UUID of the Validation Plan, Short IDs accepted")
			cmd.Spec = "UUID"

			cmd.Before = func() {
				plan = conch.GetValidationPlanByName(*idArg)
				if (plan == types.ValidationPlan{}) {
					fatal(errors.New("could not find the validation plan"))
				}
			}

			cmd.Command("get", "Get information about a single build by its name", func(cmd *cli.Cmd) {
				cmd.Action = func() { display(plan) }
			})
		})
	}
}
