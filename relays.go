// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Relays struct {
	*Conch
}

func (c *Conch) Relays() *Relays {
	return &Relays{c}
}

type RelayList []Relay

func (r RelayList) Len() int {
	return len(r)
}

func (r RelayList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RelayList) Less(i, j int) bool {
	return r[i].Updated.Before(r[j].Updated)
}

func (rl RelayList) String() string {
	sort.Sort(rl)
	if API.JsonOnly {
		return API.AsJSON(rl)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Serial",
		"Name",
		"Version",
		"IP",
		"SSH Port",
		"Updated",
	})

	for _, r := range rl {
		table.Append([]string{
			r.SerialNumber,
			r.Name,
			r.Version,
			r.IpAddr,
			strconv.Itoa(r.SshPort),
			TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()

}

type Relay struct {
	ID           uuid.UUID `json:"id" faker:"uuid"`
	SerialNumber string    `json:"serial_number"`
	Name         string    `json:"name,omitempty"`
	Version      string    `json:"version,omitempty"`
	IpAddr       string    `json:"ipaddr,omitempty" faker:"ipv4"`
	SshPort      int       `json:"ssh_port"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	LastSeen     time.Time `json:"last_seen,omitempty"`
}

func (r Relay) String() string {
	if API.JsonOnly {
		return API.AsJSON(r)
	}

	t, err := NewTemplate().Parse(relayTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, r); err != nil {
		panic(err)
	}

	return buf.String()

}

func (r *Relays) GetAll() RelayList {
	rl := make(RelayList, 0)

	res := r.Do(r.Sling().New().Get("/relay/?no_devices=1"))
	if ok := res.Parse(&rl); !ok {
		panic(res)
	}

	return rl
}

func (r *Relays) Get(identifier string) (relay Relay) {
	uri := fmt.Sprintf(
		"/relay/%s",
		identifier,
	)

	res := r.Do(r.Sling().New().Get(uri))
	if ok := res.Parse(&relay); !ok {
		panic(res)
	}
	return relay
}

func (r *Relays) Register(
	serial string,
	version string,
	ipaddr string,
	name string,
	sshPort int,
) Relay {
	if serial == "" {
		panic("please provide a serial number")
	}

	out := struct {
		Serial  string `json:"serial"`
		Version string `json:"version,omitempty"`
		IpAddr  string `json:"ipaddr,omitempty"`
		Name    string `json:"name,omitempty"`
		SshPort int    `json:"ssh_port,omitempty"`
	}{
		serial,
		version,
		ipaddr,
		name,
		sshPort,
	}

	uri := fmt.Sprintf(
		"/relay/%s/register",
		url.PathEscape(serial),
	)

	res := r.Do(
		r.Sling().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(out),
	)

	var relay Relay
	if ok := res.Parse(&relay); !ok {
		panic(res)
	}
	return relay
}

func (r *Relays) Delete(identifier string) {
	uri := fmt.Sprintf("/relay/%s", url.PathEscape(identifier))
	res := r.Do(r.Sling().New().Delete(uri))

	if res.StatusCode() != 204 {
		// I know this is weird. Like in other places, it should be impossible
		// to reach here unless the status code is 204. The API returns 204
		// (which gets us here) or 409 (which will explode before it gets here).
		// If we got here via some other code, then there's some new behavior
		// that we need to know about.
		panic(res)
	}
}

func init() {
	App.Command("relays", "Perform actions against the whole list of relays", func(cmd *cli.Cmd) {
		cmd.Command("get", "Get a list of relays", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Relays().GetAll()) }
		})

		cmd.Command("find", "Find relays by name", func(cmd *cli.Cmd) {
			var (
				relays = cmd.StringsArg("RELAYS", nil, "List of regular expressions to match against relay IDs")
				andOpt = cmd.BoolOpt("and", false, "Match the list as a logical AND")
			)

			cmd.Spec = "[OPTIONS] RELAYS..."
			cmd.LongDesc = `
Takes a list of regular expressions and matches those against the IDs of all known relays.

The default behavior is to match as a logical OR but this behavior can be changed by providing the --and flag

For instance:

* "conch relays find drd" will find all relays with 'drd' in their ID. For perl folks, this is essentially 'm/drd/'
* "conch relays find '^ams-'" will find all relays with IDs that begin with 'ams-'
* "conch relays find drd '^ams-' will find all relays with IDs that contain 'drd' OR begin with 'ams-'
* "conch relays find --and drd '^ams-' will find all relays with IDs that contain 'drd' AND begin with '^ams-'`

			cmd.Action = func() {
				if *relays == nil {
					panic("please provide a list of regular expressions")
				}

				// If a user for some strange reason gives us a relay name of "", the
				// cli lib will pass it on to us. That name is obviously useless so
				// let's filter it out.
				relayREs := make([]*regexp.Regexp, 0)
				for _, matcher := range *relays {
					if matcher == "" {
						continue
					}
					re, err := regexp.Compile(matcher)
					if err != nil {
						panic(err)
					}

					relayREs = append(relayREs, re)
				}
				if len(relayREs) == 0 {
					panic("please provide a list of regular expressions")
				}

				results := make(RelayList, 0)
				for _, relay := range API.Relays().GetAll() {
					matched := 0
					for _, re := range relayREs {
						if re.MatchString(relay.SerialNumber) {
							if *andOpt {
								matched++
							} else {
								results = append(results, relay)
								continue
							}
						}
					}
					if *andOpt {
						if matched == len(relayREs) {
							results = append(results, relay)
						}
					}
				}

				fmt.Println(results)
			}
		})
	})

	App.Command("relay", "Perform actions against a single relay", func(cmd *cli.Cmd) {
		relayArg := cmd.StringArg(
			"RELAY",
			"",
			"ID of the relay",
		)

		cmd.Spec = "RELAY"

		cmd.Command("get", "Get data about a single relay", func(cmd *cli.Cmd) {
			cmd.Action = func() { fmt.Println(API.Relays().Get(*relayArg)) }
		})

		cmd.Command("register", "Register a relay with the API", func(cmd *cli.Cmd) {

			var (
				versionOpt = cmd.StringOpt("version", "", "The version of the relay")
				sshPortOpt = cmd.IntOpt("ssh_port port", 22, "The SSH port for the relay")
				ipAddrOpt  = cmd.StringOpt("ipaddr ip", "", "The IP address for the relay")
				nameOpt    = cmd.StringOpt("name", "", "The name of the relay")
			)

			cmd.Action = func() {
				fmt.Println(API.Relays().Register(
					*relayArg,
					*versionOpt,
					*ipAddrOpt,
					*nameOpt,
					*sshPortOpt,
				))
			}
		})
		cmd.Command("delete rm", "Delete a relay", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				API.Relays().Delete(*relayArg)
				fmt.Println(API.Relays().GetAll())
			}
		})
	})

}
