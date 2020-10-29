package cli

import (
	"fmt"
	"strings"

	cli "github.com/jawher/mow.cli"
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

func buildsCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()
	log := cfg.GetLogger()
	return func(cmd *cli.Cmd) {
		cmd.Before = func() {
			conch = cfg.ConchClient()
		}
		cmd.Action = func() {
			log.Debug("display(conch.GetAllBuilds())")
			display(conch.GetAllBuilds())
		}
		cmd.Command("get ls", "Get a list of all builds", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				log.Debug("display(conch.GetAllBuilds())")
				display(conch.GetAllBuilds())
			}
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
				display(conch.GetAllBuilds())
			}
		})
	}
}

func buildCmd(cfg Config) func(*cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		var b types.Build
		buildNameArg := cmd.StringArg("NAME", "", "Name or ID of the build")

		cmd.Spec = "NAME"
		cmd.Before = func() {
			conch = cfg.ConchClient()
			b = conch.GetBuildByName(*buildNameArg)
		}

		cmd.Action = func() { display(b) }

		cmd.Command("get", "Get information about a single build by its name", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				display(b)
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
						fatal(fmt.Errorf(
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
						fatal(fmt.Errorf(
							"'role' value must be one of: %s",
							prettyBuildRoleList(),
						))
					}
					org := conch.GetOrganizationByName(*orgNameArg)

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
			cmd.Command("get ls", "Get a list of devices in an build", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					display(conch.GetAllBuildDevices(*buildNameArg))
				}
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
					build := conch.GetBuildByName(*buildNameArg)
					device := conch.GetDeviceBySerial(*deviceIDArg)
					conch.DeleteBuildDeviceByID(build.ID, device.ID)
					display(conch.GetAllBuildDevices(*buildNameArg))
				}
			})
		})

		cmd.Command("racks", "Manage racks in a specific build", func(cmd *cli.Cmd) {
			cmd.Command("get ls", "Get a list of racks in an build", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					display(conch.GetBuildRacks(*buildNameArg))
				}
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
}
