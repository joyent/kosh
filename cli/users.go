package cli

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch/types"
)

func whoamiCmd(cfg Config) func(*cli.Cmd) { return profileCmd(cfg) }

func userCmd(cfg Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("profile", "View your Conch profile", profileCmd(cfg))
		cmd.Command("settings", "Get the settings for the current user", settings(cfg))
		cmd.Command("setting", "Commands for dealing with a single setting for the current user", userSetting(cfg))
	}
}

func profileCmd(cfg Config) func(cmd *cli.Cmd) {
	display := cfg.Renderer()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := cfg.ConchClient()
			log := cfg.GetLogger()
			log.Debug("display(conch.GetCurrentUser())")
			display(conch.GetCurrentUser())
		}
	}
}

func settings(cfg Config) func(cmd *cli.Cmd) {
	display := cfg.Renderer()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := cfg.ConchClient()
			display(conch.GetCurrentUserSettings())
		}
	}
}

func userSetting(cfg Config) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		name := *cmd.StringArg("NAME", "", "The string name of a setting")
		cmd.Spec = "NAME"
		cmd.Command("get", "Get a setting for the current user", userSettingGet(cfg, name))
		cmd.Command("set", "Set a setting for the current user", userSettingSet(cfg, name))
		cmd.Command("delete", "Delete a setting for the current user", userSettingDelete(cfg, name))
	}
}

func userSettingGet(cfg Config, setting string) func(cmd *cli.Cmd) {
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := cfg.ConchClient()
			display(conch.GetCurrentUserSettingByName(setting))
		}
	}
}

func userSettingSet(cfg Config, setting string) func(cmd *cli.Cmd) {
	display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		value := *cmd.StringArg("VALUE", "", "The new value of the setting")

		cmd.Spec = "VALUE"

		cmd.Action = func() {
			conch := cfg.ConchClient()
			if e := conch.SetCurrentUserSettingByName(setting, types.UserSetting(value)); e != nil {
				fatal(e)
			}
			display(conch.GetCurrentUserSettingByName(setting))
		}
	}
}

func userSettingDelete(cfg Config, setting string) func(cmd *cli.Cmd) {
	// display := cfg.Renderer()

	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := cfg.ConchClient()
			if e := conch.DeleteCurrentUserSetting(setting); e != nil {
				fatal(e)
			}
			if !cfg.GetOutputJSON() {
				fmt.Println("OK")
			}
		}
	}
}
