package cli

import (
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch/types"
)

func whoamiCmd(cmd *cli.Cmd) { profileCmd(cmd) }

func userCmd(cmd *cli.Cmd) {
	cmd.Command("profile", "View your Conch profile", profileCmd)
	cmd.Command("settings", "Get the settings for the current user", settings)
	cmd.Command("setting", "Commands for dealing with a single setting for the current user", userSetting)
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
	cmd.Command("delete", "Delete a setting for the current user", userSettingDelete(name))
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
				fatal(e)
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
				fatal(e)
			}
		}
	}
}
