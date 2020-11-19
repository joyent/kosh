package cli

import (
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func adminCmd(cmd *cli.Cmd) {
	cmd.Before = func() {
		config.requireAuth()
		config.requireSysAdmin()
	}
	cmd.Command("users", "List all Users", adminUsersCmd)
	cmd.Command("user u", "Administrate a single User", adminUserCmd)
}

func adminUsersCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display Renderer

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()
	}

	cmd.Action = func() { display(conch.GetAllUsers()) }

	cmd.Command("get ls", "display all users", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(conch.GetAllUsers()) }
	})

	cmd.Command("create new add", "Add a new user to the system", func(cmd *cli.Cmd) {
		email := cmd.StringArg("EMAIL", "", "A user's email")
		name := cmd.StringArg("NAME", "", "A user's name")
		password := cmd.StringArg("PASS", "", "A user's initial password")
		admin := cmd.BoolOpt("admin", false, "make user a system admin")
		notify := cmd.BoolOpt("send-email", false, "notify the user via email")

		cmd.Spec = "EMAIL NAME PASS"

		cmd.Action = func() {
			display(conch.CreateUser(types.NewUser{
				Email:    types.EmailAddress(*email),
				IsAdmin:  *admin,
				Name:     types.NonEmptyString(*name),
				Password: types.NonEmptyString(*password),
			}, *notify))
		}
	})

	cmd.Command("import", "Import a new user from the JSON output", func(cmd *cli.Cmd) {
		notify := cmd.BoolOpt("send-email", false, "notify the user via email")
		filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file that defines the user. '-' indicates STDIN")

		cmd.Action = func() {
			input, e := getInputReader(*filePathArg)
			fatalIf(e)

			u := conch.ReadUser(input)
			display(
				conch.CreateUser(types.NewUser{
					Email:   u.Email,
					IsAdmin: u.IsAdmin,
					Name:    u.Name,
				}, *notify),
			)
		}
	})
}

func adminUserCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display Renderer
	var user types.UserDetailed

	email := cmd.StringArg("EMAIL", "", "A user's email")
	// cmd.Spec = "EMAIL"

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()

		var e error
		user, e = conch.GetUserByEmail(*email)
		fatalIf(e)
		// TODO check to see if user is empty
	}

	cmd.Action = func() { display(user, nil) }

	cmd.Command("get", "display all users", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(user, nil) }
	})

	cmd.Command("update", "update the information for a user", func(cmd *cli.Cmd) {
		email := cmd.StringArg("EMAIL", "", "A user's email")
		name := cmd.StringArg("NAME", "", "A user's name")
		admin := cmd.BoolOpt("admin", false, "make user a system admin")
		notify := cmd.BoolOpt("send-email", false, "notify the user via email")

		cmd.Spec = "EMAIL NAME"
		cmd.Action = func() {
			update := types.UpdateUser{}
			if *email != "" && types.EmailAddress(*email) != user.Email {
				update.Email = types.EmailAddress(*email)
			}
			if *name != "" && types.NonEmptyString(*name) != user.Name {
				update.Name = types.NonEmptyString(*name)
			}
			if *admin != user.IsAdmin {
				user.IsAdmin = *admin
			}
			conch.UpdateUser(string(user.Email), update, *notify)
		}
	})

	cmd.Command("delete rm", "remove the specified user", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch.DeleteUser(string(user.Email))
			display(conch.GetAllUsers())
		}
	})

	cmd.Command("tokens", "operate on the user's tokens", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(conch.GetUserTokens(string(user.Email))) }

		cmd.Command("get ls", "list the tokens for the current user", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(conch.GetUserTokens(string(user.Email))) }
		})
	})

	cmd.Command("token", "operate on a user's tokens", func(cmd *cli.Cmd) {
		var token types.UserToken

		name := cmd.StringArg("NAME", "", "The string name of a setting")
		cmd.Spec = "NAME"
		cmd.Before = func() {
			var e error
			token, e = conch.GetUserTokenByName(string(user.Email), *name)
			fatalIf(e)
		}

		cmd.Action = func() { display(token, nil) }

		cmd.Command("get", "information about a single token for the given user", func(cmd *cli.Cmd) {
			cmd.Action = func() { display(token, nil) }
		})

		cmd.Command("delete rm", "remove a token for the given user", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				conch.DeleteUserToken(string(user.Email), token.Name)
				display(conch.GetUserTokens(string(user.Email)))
			}
		})
	})
}
