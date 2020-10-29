package cli

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch/types"
)

func organizationsCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()
	return func(cmd *cli.Cmd) {
		cmd.Before = func() {
			conch = cfg.ConchClient()
		}
		cmd.Command("get ls", "Get a list of all organizations", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				display(conch.GetAllOrganizations())
			}
		})

		cmd.Command("create", "Create a new organization", func(cmd *cli.Cmd) {
			nameArg := cmd.StringArg("NAME", "", "Name of the new organization")

			descOpt := cmd.StringOpt("description", "", "A description of the organization")
			adminEmailArg := cmd.StringOpt(
				"admin",
				"",
				"Email address for the (initial) admin user for the organization. This does *not* create the user.",
			)

			cmd.Spec = "NAME [OPTIONS]"
			cmd.Action = func() {
				conch.CreateOrganization(types.OrganizationCreate{
					Name:        types.MojoStandardPlaceholder(*nameArg),
					Description: types.NonEmptyString(*descOpt),
					Admins: []types.Admin{
						types.Admin{Email: types.EmailAddress(*adminEmailArg)},
					},
				})
			}
		})
	}
}

func organizationCmd(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		var o types.Organization
		organizationNameArg := cmd.StringArg("NAME", "", "Name or ID of the Organization")

		cmd.Spec = "NAME"
		cmd.Before = func() {
			conch = cfg.ConchClient()
			o = conch.GetOrganizationByName(*organizationNameArg)
		}

		cmd.Command("get", "Get information about a single organization by its name", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(o)
			}
		})

		cmd.Command("delete rm", "Remove a specific organization", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteOrganization(o.ID)
			}
		})

		cmd.Command("users", "Manage users in a specific organization", func(cmd *cli.Cmd) {
			cmd.Command("get ls", "Get a list of users in an organization", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					display(conch.GetOrganizationByID(o.ID))
				}
			})

			cmd.Command("add", "Add a user to an organization", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email of the user to add to the organization. Does *not* create the user",
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
					conch.AddOrganizationUser(
						o.ID,
						types.OrganizationAddUser{
							Email: types.EmailAddress(*userEmailArg),
							Role:  types.Role(*roleOpt),
						},
						*sendEmailOpt,
					)
					display(conch.GetOrganizationByID(o.ID))
				}
			})

			cmd.Command("remove rm", "remove a user from an organization", func(cmd *cli.Cmd) {
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
					conch.DeleteOrganizationUser(
						o.ID,
						*userEmailArg,
						*sendEmailOpt,
					)
					display(conch.GetOrganizationByID(o.ID))
				}
			})
		})
	}
}
