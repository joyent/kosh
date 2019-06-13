// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jawher/mow.cli"
)

func init() {
	App.Command(
		"api",
		"",
		func(cmd *cli.Cmd) {
			cmd.Command(
				"get",
				"Perform an HTTP GET against the provided URL",
				func(cmd *cli.Cmd) {
					var urlArg = cmd.StringArg("URL", "", "The API path to GET. Must *not* include the hostname or port")
					cmd.Spec = "URL"
					cmd.Action = func() {
						fmt.Println(API.DoBadly(API.Sling().Get(*urlArg)).Body)
					}
				},
			)

			cmd.Command(
				"delete",
				"Perform an HTTP DELETE against the provided URL",
				func(cmd *cli.Cmd) {
					var urlArg = cmd.StringArg("URL", "", "The API path to DELETE. Must *not* include the hostname or port")
					cmd.Spec = "URL"
					cmd.Action = func() {
						fmt.Println(API.DoBadly(API.Sling().New().Delete(*urlArg)).Body)
					}
				},
			)

			cmd.Command(
				"post",
				"Perform an HTTP POST against the provided URL",
				func(cmd *cli.Cmd) {
					var (
						urlArg      = cmd.StringArg("URL", "", "The API path to POST. Must *not* include the hostname or port")
						filePathArg = cmd.StringArg("FILE", "-", "Path to a JSON file to use as the request body. '-' indicates STDIN")
					)
					cmd.Spec = "URL [FILE]"

					cmd.Action = func() {
						var b []byte
						var err error
						if *filePathArg == "-" {
							b, err = ioutil.ReadAll(os.Stdin)
						} else {
							b, err = ioutil.ReadFile(*filePathArg)
						}
						if err != nil {
							panic(err)
						}

						fmt.Println(API.DoBadly(
							API.Sling().New().Post(*urlArg).
								Set("Content-Type", "application/json").
								Body(bytes.NewReader(b)),
						).Body)
					}
				},
			)
			cmd.Command(
				"version",
				"Get the version of the API we are talking to",
				func(cmd *cli.Cmd) {
					cmd.Action = func() {
						if API.JsonOnly {
							fmt.Printf("{\"version\":\"%s\"}\n", API.Version())
						} else {
							fmt.Println(API.Version())
						}
					}
				},
			)

		},
	)
}
