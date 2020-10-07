package cli

import (
	cli "github.com/jawher/mow.cli"
)

const (
	ProductionURL = "https://conch.joyent.us"
	StagingURL    = "https://staging.conch.joyent.us"
)

func NewApp(config Config) *cli.Cli {

	app := cli.App("kosh", "Command line interface for Conch")

	app.Version("V version", config.Version)

	config.ConchToken = *app.String(cli.StringOpt{
		Name:   "t token",
		Value:  "",
		Desc:   "API token",
		EnvVar: "KOSH_TOKEN",
	})

	config.ConchEnvironment = *app.String(cli.StringOpt{
		Name:   "environment env",
		Value:  "production",
		Desc:   "Specify the environment to be used: production, staging, development (provide URL in the --url parameter)",
		EnvVar: "KOSH_ENV",
	})

	config.ConchURL = *app.String(cli.StringOpt{
		Name:   "u url",
		Value:  "",
		Desc:   "If the environment is 'development', this specifies the API URL. Ignored if --environment is 'production' or 'staging'",
		EnvVar: "KOSH_URL",
	})

	config.OutputJSON = *app.Bool(cli.BoolOpt{
		Name:   "j json",
		Value:  false,
		Desc:   "Output JSON only",
		EnvVar: "KOSH_JSON_ONLY",
	})

	config.StrictParse = *app.Bool(cli.BoolOpt{
		Name:   "strict",
		Value:  false,
		Desc:   "Intended for developers. Parse API responses strictly, not allowing new fields",
		EnvVar: "KOSH_DEVEL_STRICT",
	})

	config.DevMode = *app.Bool(cli.BoolOpt{
		Name:   "developer",
		Value:  false,
		Desc:   "Activate developer mode. This disables most user-friendly protections, is noisy, and switches to developer-friendly output where appropriate",
		EnvVar: "KOSH_DEVEL_MODE",
	})

	app.Command("whoami", "Display details of the current user", whoamiCmd(config))
	app.Command("user", "Commands for dealing with the current user (you)", userCmd(config))
	app.Command("devices ds", "Commands for dealing with multiple devices", devicesCmd(config))
	app.Command("device d", "Perform actions against a single device", deviceCmd(config))

	return app
}
