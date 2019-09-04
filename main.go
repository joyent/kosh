// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

const (
	ProductionURL = "https://conch.joyent.us"
	StagingURL    = "https://staging.conch.joyent.us"
	DateFormat    = "2006-01-02 15:04:05 -0700 MST"
)

var (
	Version string
	GitRev  string

	API = &Conch{}
	App = cli.App("kosh", "Command line interface for Conch")
)

func buildUserAgent() map[string]string {
	var isRoot bool
	if current, err := user.Current(); err == nil {
		if current.Uid == "0" {
			isRoot = true
		}
	}

	agentBits := make(map[string]string)
	agent := fmt.Sprintf(
		"%s (%s; %s; r=%v)",
		GitRev,
		runtime.GOOS,
		runtime.GOARCH,
		isRoot,
	)

	agentBits["Kosh"] = agent
	return agentBits
}

func TimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(DateFormat)
}

func CutUUID(id string) string {
	re := regexp.MustCompile("^(.+?)-")
	bits := re.FindStringSubmatch(id)
	if len(bits) > 0 {
		return bits[1]
	}
	return id
}

func FindUUID(id string, list []uuid.UUID) (bool, uuid.UUID) {
	re := regexp.MustCompile(fmt.Sprintf("^%s", id))
	for _, item := range list {
		if re.MatchString(item.String()) {
			return true, item
		}
	}
	return false, uuid.UUID{}
}

func Table() (table *tablewriter.Table) {
	table = tablewriter.NewWriter(os.Stdout)
	TableToMarkdown(table)
	return table
}

func TableToMarkdown(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
}

func IsSysAdmin() bool {
	return API.Users().Me().IsAdmin
}

func RequireSysAdmin() {
	if !IsSysAdmin() {
		panic("This action requires Conch systems administrator privileges")
	}
}

func NewTemplate() *template.Template {
	return template.New("wat").Funcs(template.FuncMap{
		"CutUUID": func(id string) string { return CutUUID(id) },
		"TimeStr": func(t time.Time) string { return TimeStr(t) },
	})
}

/***/

func init() {
	tokenOpt := App.String(cli.StringOpt{
		Name:   "token",
		Value:  "",
		Desc:   "API token",
		EnvVar: "KOSH_TOKEN",
	})

	environmentOpt := App.String(cli.StringOpt{
		Name:   "environment env",
		Value:  "production",
		Desc:   "Specify the environment to be used: production, staging, development (provide URL in the --url parameter)",
		EnvVar: "KOSH_ENV",
	})

	urlOpt := App.String(cli.StringOpt{
		Name:   "url",
		Value:  "",
		Desc:   "If the environment is 'development', this specifies the API URL. Ignored if --environment is 'production' or 'staging'",
		EnvVar: "KOSH_URL",
	})

	jsonOnlyOpt := App.Bool(cli.BoolOpt{
		Name:   "json",
		Value:  false,
		Desc:   "Output JSON only",
		EnvVar: "KOSH_JSON_ONLY",
	})

	strictParseOpt := App.Bool(cli.BoolOpt{
		Name:   "strict",
		Value:  false,
		Desc:   "Intended for developers. Parse API responses strictly, not allowing new fields",
		EnvVar: "KOSH_DEVEL_STRICT",
	})

	develOpt := App.Bool(cli.BoolOpt{
		Name:   "developer",
		Value:  false,
		Desc:   "Activate developer mode. This disables most user-friendly protections, is noisy, and switches to developer-friendly output where appropriate",
		EnvVar: "KOSH_DEVEL_MODE",
	})

	App.Before = func() {
		if len(*environmentOpt) > 0 {
			if (*environmentOpt == "development") && (len(*urlOpt) == 0) {
				panic("--url must be provided if --environment=development is set")
			}
		}

		switch *environmentOpt {
		case "staging":
			API.URL = StagingURL
		case "development":
			API.URL = *urlOpt
		default:
			API.URL = ProductionURL
		}

		if *tokenOpt == "" {
			panic("please provide a token")
		}

		API.JsonOnly = *jsonOnlyOpt
		API.Token = *tokenOpt
		API.UserAgent = buildUserAgent()
		API.StrictParsing = *strictParseOpt
		API.DevelMode = *develOpt
	}
}

func main() {
	defer errorHandler()
	// BUG(sungo): github version check foo
	_ = App.Run(os.Args)
}
