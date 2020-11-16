package cli

import (
	"errors"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func validationCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)
	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()
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
			var e error
			plan, e = conch.GetValidationPlanByName(*idArg)
			if e != nil {
				fatal(e)
			}

			if (plan == types.ValidationPlan{}) {
				fatal(errors.New("could not find the validation plan"))
			}
		}

		cmd.Command("get", "Get information about a single build by its name", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(plan, nil) }
		})
	})
}
