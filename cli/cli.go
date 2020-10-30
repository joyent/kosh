package cli

import (
	"fmt"
	"io"
	"os"

	cli "github.com/jawher/mow.cli"
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

var config Config

// NewApp creates a new kosh app, takes a cli.Config and returns an instance of cli.Cli
func NewApp(cfg Config) *cli.Cli {
	config = cfg

	app := cli.App("kosh", "Command line interface for Conch")
	app.Spec = "[-vVdj]"

	app.Version("V version", config.Version)

	app.StringPtr(&config.ConchToken, cli.StringOpt{
		Name:   "t token",
		Value:  "",
		Desc:   "API token",
		EnvVar: "KOSH_TOKEN",
	})

	app.StringPtr(&config.ConchURL, cli.StringOpt{
		Name:   "u url",
		Value:  productionURL,
		Desc:   "This specifies the API URL.",
		EnvVar: "KOSH_URL",
	})

	app.BoolPtr(&config.OutputJSON, cli.BoolOpt{
		Name:   "j json",
		Value:  false,
		Desc:   "Output JSON only",
		EnvVar: "KOSH_JSON_ONLY",
	})

	app.BoolPtr(&config.Logger.LevelDebug, cli.BoolOpt{
		Name:   "d debug",
		Value:  false,
		Desc:   "Enable Debugging output (for debugging purposes *very* noisy). ",
		EnvVar: "KOSH_DEBUG_MODE",
	})

	app.BoolPtr(&config.Logger.LevelInfo, cli.BoolOpt{
		Name:   "v verbose",
		Value:  false,
		Desc:   "Enable Verbose Output",
		EnvVar: "KOSH_VERBOSE_MODE",
	})

	app.Command("build b", "Work with a specific build", buildCmd)
	app.Command("builds bs", "Work with builds", buildsCmd)
	app.Command("datacenter dc", "Deal with a single datacenter", datacenterCmd)
	app.Command("datacenters dcs", "Work with the datacenters you have access to", datacentersCmd)
	app.Command("device d", "Perform actions against a single device", deviceCmd)
	app.Command("device-report dr", "Deal with device reports", deviceReportCmd)
	app.Command("devices ds", "Commands for dealing with multiple devices", devicesCmd)
	app.Command("hardware h", "Work with hardware profiles and vendors", hardwareCmd)
	app.Command("organization org", "Work with a specific organization", organizationCmd)
	app.Command("organizations orgs", "Work with organizations", organizationsCmd)
	app.Command("rack r", "Work with a single rack", rackCmd)
	app.Command("racks rs", "Work with datacenter racks", racksCmd)
	app.Command("relay", "Perform actions against a single relay", relayCmd)
	app.Command("relays", "Perform actions against the whole list of relays", relaysCmd)
	app.Command("roles", "Work with datacenter rack roles", rolesCmd)
	app.Command("role", "Work with a single rack role", roleCmd)
	app.Command("room", "Deal with a single datacenter room", roomCmd)
	app.Command("rooms", "Work with datacenter rooms", roomsCmd)
	app.Command("schema", "Get the server JSON Schema for a given request or response", schemaCmd)
	app.Command("user u", "Commands for dealing with the current user (you)", userCmd)
	app.Command("validation v", "Work with validations", validationCmd)
	app.Command("whoami", "Display details of the current user", whoamiCmd)

	app.Command("version", "Get more detailed version info than --version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch := config.ConchClient()
			display := config.Renderer()
			fmt.Printf(
				"Kosh %s\n"+"  Git Revision: %s\n",
				config.Version,
				config.GitRev,
			)
			display(conch.Version())
		}
	})

	app.Before = func() {
		if config.ConchToken == "" {
			fmt.Println("Need to provide --token or set KOSH_TOKEN")
			cli.Exit(1)
		}

		config.Debug("Starting App")
		config.Info(config)
	}

	return app
}
