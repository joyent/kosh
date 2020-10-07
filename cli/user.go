package cli

import (
	"fmt"
	"log"

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
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			user := conch.GetCurrentUser()
			if cfg.OutputJSON {
				fmt.Println(user.JSON())
			} else {
				fmt.Println(user.String())
			}
		}
	}
}

func settings(cfg Config) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := conch.GetCurrentUserSettings()
			if cfg.OutputJSON {
				fmt.Println(settings.JSON())
			} else {
				fmt.Println(settings.String())
			}
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
	conch := cfg.ConchClient()

	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			setting := conch.GetCurrentUserSettingByName(setting)
			if cfg.OutputJSON {
				fmt.Println(setting.JSON())
			} else {
				fmt.Println(setting.String())
			}
		}
	}
}

func userSettingSet(cfg Config, setting string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()

	return func(cmd *cli.Cmd) {
		value := *cmd.StringArg("VALUE", "", "The new value of the setting")

		cmd.Spec = "VALUE"

		cmd.Action = func() {
			if e := conch.SetCurrentUserSettingByName(setting, types.UserSetting(value)); e != nil {
				log.Fatal(e)
			}
			setting := conch.GetCurrentUserSettingByName(setting)
			if cfg.OutputJSON {
				fmt.Println(setting.JSON())
			} else {
				fmt.Println(setting.String())
			}
		}
	}
}

func userSettingDelete(cfg Config, setting string) func(cmd *cli.Cmd) {
	conch := cfg.ConchClient()

	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			if e := conch.DeleteCurrentUserSetting(setting); e != nil {
				log.Fatal(e)
			}
			if !cfg.OutputJSON {
				fmt.Println("OK")
			}
		}
	}
}
