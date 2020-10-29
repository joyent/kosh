package main

import (
	"os"

	"github.com/joyent/kosh/cli"
)

var (
	// Version holds the current version of the app, injected as part of the build process
	Version string

	// GitRev holds the current Git Revision, injected as part of the build process
	GitRev string
)

func main() {
	app := cli.NewApp(cli.NewConfig(Version, GitRev))
	app.Run(os.Args)
}
