package cli

import (
	"fmt"
	"io"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/logger"
)

const (
	productionURL = "https://conch.joyent.us"
	stagingURL    = "https://staging.conch.joyent.us"
	edgeURL       = "https://edge.conch.joyent.us"
)

func fatal(e error) {
	fmt.Println(e)
	cli.Exit(1)
}

func getInputReader(filePathArg string) (io.Reader, error) {
	if filePathArg == "-" {
		return os.Stdin, nil
	}
	return os.Open(filePathArg)
}

func requireSysAdmin(c Config) func() {
	return func() {
		if !c.ConchClient().IsSysAdmin() {
			fmt.Println("This action requires Conch systems administrator privileges")
			cli.Exit(1)
		}
	}
}

// NewApp creates a new kosh app, takes a cli.Config and returns an instance of cli.Cli
func NewApp(config Config) *cli.Cli {
	app := cli.App("kosh", "Command line interface for Conch")
	app.Spec = "[-vVdj]"

	app.Version("V version", config.GetVersion())

	conchToken := app.String(cli.StringOpt{
		Name:   "t token",
		Value:  "",
		Desc:   "API token",
		EnvVar: "KOSH_TOKEN",
	})

	conchURL := app.String(cli.StringOpt{
		Name:   "u url",
		Value:  productionURL,
		Desc:   "This specifies the API URL.",
		EnvVar: "KOSH_URL",
	})

	outputJSON := app.Bool(cli.BoolOpt{
		Name:   "j json",
		Value:  false,
		Desc:   "Output JSON only",
		EnvVar: "KOSH_JSON_ONLY",
	})

	levelDebug := app.Bool(cli.BoolOpt{
		Name:   "d debug",
		Value:  false,
		Desc:   "Enable Debugging output (for debugging purposes *very* noisy). ",
		EnvVar: "KOSH_DEBUG_MODE",
	})

	levelInfo := app.Bool(cli.BoolOpt{
		Name:   "v verbose",
		Value:  false,
		Desc:   "Enable Verbose Output",
		EnvVar: "KOSH_VERBOSE_MODE",
	})

	app.Command("build b", "Work with a specific build", buildCmd(config))
	app.Command("builds bs", "Work with builds", buildsCmd(config))
	app.Command("datacenter dc", "Deal with a single datacenter", datacenterCmd(config))
	app.Command("datacenters dcs", "Work with the datacenters you have access to", datacentersCmd(config))
	app.Command("device d", "Perform actions against a single device", deviceCmd(config))
	app.Command("device-report dr", "Deal with device reports", deviceReportCmd(config))
	app.Command("devices ds", "Commands for dealing with multiple devices", devicesCmd(config))
	app.Command("hardware h", "Work with hardware profiles and vendors", hardwareCmd(config))
	app.Command("organization org", "Work with a specific organization", organizationCmd(config))
	app.Command("organizations orgs", "Work with organizations", organizationsCmd(config))
	app.Command("rack r", "Work with a single rack", rackCmd(config))
	app.Command("racks rs", "Work with datacenter racks", racksCmd(config))
	app.Command("relay", "Perform actions against a single relay", relayCmd(config))
	app.Command("relays", "Perform actions against the whole list of relays", relaysCmd(config))
	app.Command("roles", "Work with datacenter rack roles", rolesCmd(config))
	app.Command("role", "Work with a single rack role", roleCmd(config))
	app.Command("room", "Deal with a single datacenter room", roomCmd(config))
	app.Command("rooms", "Work with datacenter rooms", roomsCmd(config))
	app.Command("schema", "Get the server JSON Schema for a given request or response", schemaCmd(config))
	app.Command("user u", "Commands for dealing with the current user (you)", userCmd(config))
	app.Command("validation v", "Work with validations", validationCmd(config))
	app.Command("whoami", "Display details of the current user", whoamiCmd(config))

	app.Command("version", "Get more detailed version info than --version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Printf(
				"Kosh %s\n"+"  Git Revision: %s\n",
				config.GetVersion(),
				config.GetGitRev(),
			)
		}
	})

	app.Before = func() {
		if *conchToken == "" {
			fmt.Println("Need to provide --token or set KOSH_TOKEN")
			cli.Exit(1)
		}

		config.SetURL(*conchURL)
		config.SetToken(*conchToken)
		config.SetOutputJSON(*outputJSON)
		config.SetLogger(logger.Logger{
			LevelDebug: *levelDebug,
			LevelInfo:  *levelInfo,
		})
		config.GetLogger().Debug("Starting App")
		config.GetLogger().Info(config)
	}

	return app
}
