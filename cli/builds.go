package cli

import (
	"fmt"
	"strings"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

var buildRoleList = []string{"admin", "rw", "ro"}

func prettyBuildRoleList() string {
	return strings.Join(buildRoleList, ", ")
}

func okBuildRole(role string) bool {
	for _, b := range buildRoleList {
		if role == b {
			return true
		}
	}
	return false
}

func buildsCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)

	cmd.Before = func() {
		config.requireAuth()
		conch = config.ConchClient()
		display = config.Renderer()
	}

	var startedSetByUser bool
	var completedSetByUser bool

	started := cmd.Bool(cli.BoolOpt{
		Name:      "s started",
		Value:     false,
		Desc:      "display started builds",
		SetByUser: &startedSetByUser,
	})

	completed := cmd.Bool(cli.BoolOpt{
		Name:      "c completed",
		Value:     false,
		Desc:      "display completed builds",
		SetByUser: &completedSetByUser,
	})

	getAllBuilds := func() {
		params := make(map[string]string)
		if startedSetByUser {
			if *started {
				params["started"] = "1"
			} else {
				params["started"] = "0"
			}
		}
		if completedSetByUser {
			if *completed {
				params["completed"] = "1"
			} else {
				params["completed"] = "0"
			}
		}
		display(conch.GetAllBuilds(params))
	}

	// Default action is to get all builds
	cmd.Action = getAllBuilds

	cmd.Command("get ls", "Get a list of all builds", func(cmd *cli.Cmd) {
		cmd.Action = getAllBuilds
	})

	cmd.Command("create", "Create a new build", func(cmd *cli.Cmd) {
		nameArg := cmd.StringArg("NAME", "", "Name of the new build")

		descOpt := cmd.StringOpt("description", "", "A description of the build")
		adminEmailArg := cmd.StringOpt(
			"admin",
			"",
			"Email address for the (initial) admin user for the build. This does *not* create the user.",
		)

		cmd.Spec = "NAME [OPTIONS]"
		cmd.Action = func() {
			conch.CreateBuild(
				types.BuildCreate{
					Name:        types.MojoStandardPlaceholder(*nameArg),
					Description: types.NonEmptyString(*descOpt),
					Admins:      []types.Admin{types.Admin{Email: types.EmailAddress(*adminEmailArg)}},
				},
			)
			getAllBuilds()
		}
	})
}

func buildCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)
	var build types.Build

	buildNameArg := cmd.StringArg("NAME", "", "Name or ID of the build")
	cmd.Spec = "NAME"

	cmd.Before = func() {
		config.requireAuth()

		conch = config.ConchClient()
		display = config.Renderer()

		var e error
		build, e = conch.GetBuildByName(*buildNameArg)
		fatalIf(e)
	}

	cmd.Action = func() { display(build, nil) }

	cmd.Command("get", "Get information about a single build by its name", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(build, nil) }
	})

	cmd.Command("start", "Mark the build as started", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			e := conch.UpdateBuildByID(build.ID, types.BuildUpdate{
				Started: time.Now(),
			})
			fatalIf(e)
			display(conch.GetBuildByID(build.ID))
		}
	})

	cmd.Command("complete", "Mark the build as completed", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			update := types.BuildUpdate{Completed: time.Now()}

			e := conch.UpdateBuildByID(build.ID, update)
			fatalIf(e)

			display(conch.GetBuildByID(build.ID))
		}
	})

	cmd.Command("users", "Manage users in a specific build", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			display(conch.GetBuildUsers(*buildNameArg))
		}
		cmd.Command("get ls", "Get a list of users in an build", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				display(conch.GetBuildUsers(*buildNameArg))
			}
		})

		cmd.Command("add", "Add a user to an build", func(cmd *cli.Cmd) {
			userEmailArg := cmd.StringArg(
				"EMAIL",
				"",
				"The email of the user to add to the build. Does *not* create the user",
			)

			roleOpt := cmd.StringOpt(
				"role",
				"ro",
				"The role for the user. One of: "+prettyBuildRoleList(),
			)

			sendEmailOpt := cmd.BoolOpt(
				"send-email",
				true,
				"Send email to the target user, notifying them of the change",
			)

			cmd.Spec = "EMAIL [OPTIONS]"
			cmd.Action = func() {
				if !okBuildRole(*roleOpt) {
					fatalIf(fmt.Errorf(
						"'role' value must be one of: %s",
						prettyBuildRoleList(),
					))
				}
				conch.AddBuildUser(
					*buildNameArg,
					types.BuildAddUser{
						Email: types.EmailAddress(*userEmailArg),
						Role:  types.Role(*roleOpt),
					},
					*sendEmailOpt,
				)
				display(conch.GetBuildUsers(*buildNameArg))
			}
		})

		cmd.Command("remove rm", "remove a user from an build", func(cmd *cli.Cmd) {
			userEmailArg := cmd.StringArg(
				"EMAIL",
				"",
				"The email or ID of the user to modify",
			)

			sendEmailOpt := cmd.BoolOpt(
				"send-email",
				true,
				"Send email to the target user, notifying them of the change",
			)
			cmd.Spec = "EMAIL [OPTIONS]"
			cmd.Action = func() {
				conch.DeleteBuildUser(*buildNameArg, *userEmailArg, *sendEmailOpt)
				display(conch.GetBuildUsers(*buildNameArg))
			}
		})
	})

	cmd.Command("organizations orgs", "Manage organizations in a specific build", func(cmd *cli.Cmd) {
		cmd.Command("get ls", "Get a list of organizations in an build", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				display(conch.GetAllBuildOrganizations(*buildNameArg))
			}
		})

		cmd.Command("add", "Add a organization to an build", func(cmd *cli.Cmd) {
			orgNameArg := cmd.StringArg(
				"NAME",
				"",
				"The name of the organization to add to the build. Does *not* create the organization",
			)

			roleOpt := cmd.StringOpt(
				"role",
				"ro",
				"The role for the organization. One of: "+prettyBuildRoleList(),
			)

			sendEmailOpt := cmd.BoolOpt(
				"send-email",
				true,
				"Send email to the organization admins, notifying them of the change",
			)

			cmd.Spec = "NAME [OPTIONS]"
			cmd.Action = func() {
				if !okBuildRole(*roleOpt) {
					fatalIf(fmt.Errorf(
						"'role' value must be one of: %s",
						prettyBuildRoleList(),
					))
				}
				org, e := conch.GetOrganizationByName(*orgNameArg)
				fatalIf(e)

				conch.AddBuildOrganization(*buildNameArg, types.BuildAddOrganization{
					org.ID,
					types.Role(*roleOpt),
				},
					*sendEmailOpt,
				)
				display(conch.GetAllBuildOrganizations(*buildNameArg))
			}
		})

		cmd.Command("remove rm", "remove an organization from a build", func(cmd *cli.Cmd) {
			orgNameArg := cmd.StringArg(
				"NAME",
				"",
				"The name or ID of the organization to modify",
			)

			sendEmailOpt := cmd.BoolOpt(
				"send-email",
				true,
				"Send email to the target organization admins, notifying them of the change",
			)
			cmd.Spec = "EMAIL [OPTIONS]"
			cmd.Action = func() {
				conch.DeleteBuildOrganization(*buildNameArg,
					*orgNameArg,
					*sendEmailOpt,
				)
				display(conch.GetAllBuildOrganizations(*buildNameArg))
			}
		})
	})

	cmd.Command("devices ds", "Manage devices in a specific build", func(cmd *cli.Cmd) {
		// list by default
		cmd.Action = func() { display(conch.GetAllBuildDevices(*buildNameArg)) }

		cmd.Command("get ls", "Get a list of devices in an build", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetAllBuildDevices(*buildNameArg)) }
		})

		cmd.Command("add", "Add a device to an build", func(cmd *cli.Cmd) {
			deviceIDArg := cmd.StringArg(
				"ID",
				"",
				"The ID or serial number of the device to add to the build. Does *not* create the device",
			)

			cmd.Spec = "ID [OPTIONS]"
			cmd.Action = func() {
				conch.AddBuildDeviceByName(*buildNameArg, *deviceIDArg)
				display(conch.GetAllBuildDevices(*buildNameArg))
			}
		})

		cmd.Command("remove rm", "remove a device from a build", func(cmd *cli.Cmd) {
			deviceIDArg := cmd.StringArg(
				"ID",
				"",
				"The ID or serial number of the device to add to the build. Does *not* create the device",
			)

			cmd.Spec = "ID [OPTIONS]"
			cmd.Action = func() {
				b, e := conch.GetBuildByName(*buildNameArg)
				fatalIf(e)

				d, e := conch.GetDeviceBySerial(*deviceIDArg)
				fatalIf(e)

				conch.DeleteBuildDeviceByID(b.ID, d.ID)
				display(conch.GetAllBuildDevices(*buildNameArg))
			}
		})
	})

	cmd.Command("racks", "Manage racks in a specific build", func(cmd *cli.Cmd) {
		// default to list
		cmd.Action = func() { display(conch.GetBuildRacks(*buildNameArg)) }

		cmd.Command("get ls", "Get a list of racks in an build", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetBuildRacks(*buildNameArg)) }
		})

		cmd.Command("add", "Add a rack to an build", func(cmd *cli.Cmd) {
			rackIDArg := cmd.StringArg(
				"ID",
				"",
				"The ID of the rack to add to the build. Does *not* create the rack",
			)

			cmd.Spec = "ID [OPTIONS]"
			cmd.Action = func() {
				conch.AddBuildRackByID(*buildNameArg, *rackIDArg)
				display(conch.GetBuildRacks(*buildNameArg))
			}
		})

		cmd.Command("remove rm", "remove a rack from a build", func(cmd *cli.Cmd) {
			rackIDArg := cmd.StringArg(
				"ID",
				"",
				"The ID of the rack to add to the build. Does *not* create the device",
			)

			cmd.Spec = "ID [OPTIONS]"
			cmd.Action = func() {
				conch.DeleteBuildRackByID(*buildNameArg, *rackIDArg)
				display(conch.GetBuildRacks(*buildNameArg))
			}
		})
	})
}
