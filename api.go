// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"
)

type Conch struct {
	URL   string
	Token string

	UserAgent map[string]string

	Debug         bool
	Trace         bool
	JsonOnly      bool
	StrictParsing bool
	DevelMode     bool

	HTTP *http.Client
}

var defaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		DualStack: true,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

func (c *Conch) PrintJSON(d interface{}) {
	fmt.Println(c.AsJSON(d))
}

func (c *Conch) AsJSON(d interface{}) string {
	if j, err := json.MarshalIndent(d, "", "    "); err != nil {
		panic(err)
	} else {
		return string(j)
	}
}

func (c *Conch) Sling() *sling.Sling {
	userAgent := fmt.Sprintf("Conch/%s", Version)
	if len(c.UserAgent) > 0 {
		for k, v := range c.UserAgent {
			userAgent = fmt.Sprintf("%s %s/%s", userAgent, k, v)
		}
	}

	if c.HTTP == nil {
		c.HTTP = &http.Client{
			Transport: defaultTransport,
		}
	}

	s := sling.New().
		Client(c.HTTP).
		Set("User-Agent", userAgent)

	if c.URL != "" {
		s = s.Base(c.URL)
	}

	if c.Token != "" {
		s = s.Set("Authorization", "Bearer "+c.Token)
	}

	return s.New()
}

func (c *Conch) DoBadly(s *sling.Sling) *ConchResponse {
	response := ConchResponse{Strict: c.StrictParsing}

	req, err := s.Request()
	if err != nil {
		response.Error = err
		panic(&response)
	}

	response.Request = req

	// BUG(sungo) Logging

	res, err := c.HTTP.Do(req)
	response.Response = res

	if (res == nil) || (err != nil) {
		response.Error = err
		panic(&response)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		response.Error = err
		panic(&response)
	}

	response.Body = string(bodyBytes)
	return &response
}

func (c *Conch) Do(s *sling.Sling) *ConchResponse {
	response := ConchResponse{Strict: c.StrictParsing}

	req, err := s.Request()
	if err != nil {
		response.Error = err
		panic(&response)
	}

	response.Request = req

	// BUG(sungo) Logging

	res, err := c.HTTP.Do(req)
	response.Response = res

	if (res == nil) || (err != nil) {
		response.Error = err
		panic(&response)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		response.Error = err
		panic(&response)
	}

	response.Body = string(bodyBytes)

	if response.IsError() {
		e := struct {
			Error string `json:"error"`
		}{}

		if ok := response.Parse(&e); ok {
			response.Error = fmt.Errorf(
				"HTTP Error %d: %s",
				response.StatusCode(),
				e.Error,
			)
		} else {
			response.Error = fmt.Errorf(
				"HTTP Error %d: %s",
				response.StatusCode(),
				response.Status(),
			)
		}
		panic(&response)
	}

	switch response.StatusCode() {
	case 201:
		fallthrough
	case 204:
		if location, err := response.Response.Location(); err == nil {
			if location != nil {
				return c.Do(c.Sling().Get(location.String()))
			}
		}
	}

	return &response
}

func (c *Conch) Version() string {
	res := c.Do(c.Sling().Get("/version"))

	v := struct {
		Version string `json:"version"`
	}{}

	if ok := res.Parse(&v); !ok {
		panic(res)
	}

	return v.Version
}

/*****************/

// ErrNoResponse indicates an operation was attempted on the structure when no
// HTTP response is present
var ErrNoResponse = errors.New("no HTTP response found")

// ConchResponse holds the notions of what happened to an HTTP request and
// convenience functions around payload parsing and the like
type ConchResponse struct {
	Request  *http.Request
	Response *http.Response

	Strict bool
	Body   string
	Error  error
}

/***/

// StatusCode provides the HTTP response status code. If we don't have a
// response, -1 is returned.
func (r *ConchResponse) StatusCode() int {
	if r.Response == nil {
		return -1
	}
	return r.Response.StatusCode
}

// Status provides the string version of the HTTP status code. If we don't have
// a response, "" is returned.
func (r *ConchResponse) Status() string {
	if r.Response == nil {
		return ""
	}
	return r.Response.Status
}

// IsError provides a really simplistic notion of when an HTTP response has
// gone awry. Specifically, if the status code is between 400 and 600, it is
// considered in error. The response is also considered to be ok/successful if
// it hasn't happened yet.
func (r *ConchResponse) IsError() bool {
	if r.Response == nil {
		return false
	}

	if (r.StatusCode() >= 400) && (r.StatusCode() < 600) {
		return true
	}

	return false
}

// IsErrorOurFault is a convenience function to spot when an HTTP error code is
// in the 400s
func (r *ConchResponse) IsErrorOurFault() bool {
	if !r.IsError() {
		return false
	}

	if r.StatusCode() >= 400 && r.StatusCode() < 500 {
		return true
	}

	return false
}

// IsErrorTheirFault is a convenience function to spot when an HTTP error code
// is in the 500s
func (r *ConchResponse) IsErrorTheirFault() bool {
	if !r.IsError() {
		return false
	}

	if r.StatusCode() >= 500 && r.StatusCode() < 600 {
		return true
	}

	return false
}

// Parse, well, parses the JSON payload in the HTTP response and tries to shove
// it into the provided structure. If Strict is true, the parser disallows
// any unknown fields. To quote the go docs, an error is returned when "the
// input contains object keys which do not match any non-ignored, exported
// fields in the destination."
//
// So, if the API sends us data we aren't expecting, you'll get an error.
// However, you will *not* get an error if the API fails to send data you're
// expecting.
func (r *ConchResponse) Parse(data interface{}) bool {
	r.Error = nil
	if r.Response == nil {
		r.Error = ErrNoResponse
		return false
	}

	dec := json.NewDecoder(strings.NewReader(r.Body))
	if r.Strict {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(data)
	r.Error = err
	return (err == nil)
}
