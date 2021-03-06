package cli

import (
	"errors"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func rolesCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display Renderer

	cmd.Before = func() {
		config.requireSysAdmin()
		conch = config.ConchClient()
		display = config.Renderer()
	}

	cmd.Command("get", "Get a list of all rack roles", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(conch.GetAllRackRoles()) }
	})

	cmd.Command("create", "Create a new rack role", func(cmd *cli.Cmd) {
		var (
			nameOpt     = cmd.StringOpt("name", "", "The name of the role")
			rackSizeOpt = cmd.IntOpt("rack-size", 0, "Size of the rack necessary for this role")
		)

		cmd.Spec = "--name --rack-size"
		cmd.Action = func() {
			if *nameOpt == "" {
				fatalIf(errors.New("--name is required"))
			}

			if *rackSizeOpt == 0 {
				fatalIf(errors.New("--rack-size is required and cannot be 0"))
			}
			conch.CreateRackRole(types.RackRoleCreate{
				Name:     types.MojoStandardPlaceholder(*nameOpt),
				RackSize: types.PositiveInteger(*rackSizeOpt),
			})
		}
	})
}

func roleCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)
	var role types.RackRole

	nameArg := cmd.StringArg(
		"NAME",
		"",
		"The name of the rack role",
	)

	cmd.Spec = "NAME"

	cmd.Before = func() {
		config.requireAuth()
		config.requireSysAdmin()

		conch := config.ConchClient()
		display = config.Renderer()

		var e error
		role, e = conch.GetRackRoleByName(*nameArg)
		if e != nil {
			fatalIf(e)
		}
		if (role == types.RackRole{}) {
			fatalIf(errors.New("couldn't find the role"))
		}
	}

	cmd.Command("get", "Get information about a single rack role", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(role, nil) }
	})

	cmd.Command("update", "Update information about a single rack role", func(cmd *cli.Cmd) {
		var (
			nameOpt     = cmd.StringOpt("name", "", "The name of the role")
			rackSizeOpt = cmd.IntOpt("rack-size", 0, "Size of the rack necessary for this role")
		)

		cmd.Action = func() {
			conch.UpdateRackRole(role.ID, types.RackRoleUpdate{
				Name:     types.MojoStandardPlaceholder(*nameOpt),
				RackSize: types.PositiveInteger(*rackSizeOpt),
			})
			display(conch.GetAllRackRoles())
		}
	})

	cmd.Command("delete", "Delete a single rack role", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch.DeleteRackRole(role.ID)
			display(conch.GetAllRackRoles())
		}
	})
}
