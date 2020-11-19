package cli

import (
	"errors"
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func whoamiCmd(cmd *cli.Cmd) {
	cmd.Before = config.requireAuth
	profileCmd(cmd)
}

func userCmd(cmd *cli.Cmd) {
	cmd.Before = config.requireAuth
	cmd.Command("profile", "View your Conch profile", profileCmd)
	cmd.Command("settings", "Get the settings for the current user", settingsCmd)
	cmd.Command("setting", "Commands for dealing with a single setting for the current user", userSetting)
	cmd.Command("tokens", "Commands for dealing with the current user's tokens", tokensCmd)
	cmd.Command("token", "Commands for dealing with a single token for the current user", tokenCmd)
}

func tokensCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display Renderer

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()
	}
	cmd.Action = func() { display(conch.GetCurrentUserTokens()) }

	cmd.Command("get ls", "list the tokens for the current user", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(conch.GetCurrentUserTokens()) }
	})
	cmd.Command("create new add", "Get the settings for the current user", func(cmd *cli.Cmd) {
		name := cmd.StringArg("NAME", "", "The string name of a setting")
		user := cmd.StringOpt("user u", "", "User name to use for authentication")
		pass := cmd.StringOpt("pass p", "", "Password to use for authentication")
		cmd.Action = func() {
			if *user != "" && *pass != "" {
				loginToken, e := conch.Login(*user, *pass)
				if e != nil {
					fatalIf(e)
				}
				config.Debug(fmt.Sprintf("%+v", loginToken))
				conch = conch.Authorization("Bearer " + loginToken.JwtToken)
			}
			display(conch.CreateCurrentUserToken(types.NewUserTokenRequest{Name: *name}))
		}
	})
}

func tokenCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display Renderer
	var token types.UserToken

	name := cmd.StringArg("NAME", "", "The string name of a setting")
	cmd.Spec = "NAME"

	cmd.Before = func() {
		conch = config.ConchClient()
		display = config.Renderer()

		var e error
		if name == nil {
			fatalIf(errors.New("must provide a valid token name"))
		}
		token, e = conch.GetCurrentUserTokenByName(*name)
		if e != nil {
			fatalIf(e)
		}
	}

	cmd.Action = func() { display(token, nil) }

	cmd.Command("get", "display the user token information", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(token, nil) }
	})

	cmd.Command("delete rm", "display the user token information", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch.DeleteCurrentUserToken(token.Name)
			display(conch.GetCurrentUserTokens())
		}
	})
}

func profileCmd(cmd *cli.Cmd) {
	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()

		display(conch.GetCurrentUser())
	}
}

func settingsCmd(cmd *cli.Cmd) {
	cmd.Action = func() {
		conch := config.ConchClient()
		display := config.Renderer()
		display(conch.GetCurrentUserSettings())
	}
}

func userSetting(cmd *cli.Cmd) {
	name := *cmd.StringArg("NAME", "", "The string name of a setting")
	cmd.Spec = "NAME"
	cmd.Command("get", "Get a setting for the current user", userSettingGet(name))
	cmd.Command("set", "Set a setting for the current user", userSettingSet(name))
	cmd.Command("delete rm", "Delete a setting for the current user", userSettingDelete(name))
}

func userSettingGet(setting string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()
			display(conch.GetCurrentUserSettingByName(setting))
		}
	}
}

func userSettingSet(setting string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		value := *cmd.StringArg("VALUE", "", "The new value of the setting")

		cmd.Spec = "VALUE"

		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()
			if e := conch.SetCurrentUserSettingByName(setting, types.UserSetting(value)); e != nil {
				fatalIf(e)
			}
			display(conch.GetCurrentUserSettingByName(setting))
		}
	}
}

func userSettingDelete(setting string) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			if e := conch.DeleteCurrentUserSetting(setting); e != nil {
				fatalIf(e)
			}
		}
	}
}
