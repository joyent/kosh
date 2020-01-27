// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	"encoding/json"
	"fmt"
	"net/http/httputil"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
)

func errorHandler() {
	if r := recover(); r != nil {
		if API.DevelMode {

			if reflect.TypeOf(r).String() == "*main.ConchResponse" {
				res := r.(*ConchResponse)

				reqDump, err := httputil.DumpRequest(res.Request, true)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Fprintf(os.Stderr,
					"HTTP Request: %s\n\n",
					reqDump,
				)

				fmt.Fprintf(os.Stderr,
					"HTTP Status Code: %d\nHTTP Status: %s\nString Error: %s\n\n",
					res.StatusCode(),
					res.Status(),
					res.Error.Error(),
				)

				fmt.Fprintf(
					os.Stderr,
					"RAW HTTP RESPONSE:\n%s\n\n",
					res.Body,
				)

				var s interface{}
				if err := json.Unmarshal([]byte(res.Body), &s); err == nil {
					fmt.Fprintln(os.Stderr, "MARSHALLED RESPONSE:")
					spew.Fdump(os.Stderr, s)
				}

			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n\n", r)
				fmt.Fprintf(os.Stderr, "RAW ERROR: ")
				spew.Fdump(os.Stderr, r)
			}

			fmt.Fprintln(os.Stderr)
			debug.PrintStack()
			os.Exit(1)
		}

		var msg string
		if reflect.TypeOf(r).String() == "*main.ConchResponse" {
			res := r.(*ConchResponse)
			if res.Error != nil {
				msg = res.Error.Error()
			} else {
				msg = fmt.Sprintf("An HTTP error occured: %s", res.Status())
			}
		} else {
			msg = fmt.Sprintf("An error occurred: %s", r)
		}

		if API.JsonOnly {
			fmt.Fprintf(os.Stderr, "{\"error\":\"%s\"}\n", msg)
		} else {
			fmt.Fprintln(os.Stderr, msg)
		}

		os.Exit(1)
	}
}

func Spew(d interface{}) {
	spew.Fdump(os.Stderr, d)
}

func init() {
	spew.Config = spew.ConfigState{
		Indent:                  "    ",
		SortKeys:                true,
		DisablePointerAddresses: true,
		DisableMethods:          true,
		DisableCapacities:       true,
		DisablePointerMethods:   true,
		SpewKeys:                true,
	}
}
