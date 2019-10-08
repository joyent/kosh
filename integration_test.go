package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func setupAPIClient() {
	Version += "-testing"
	API.URL = os.Getenv("KOSH_URL")
	API.Token = os.Getenv("KOSH_TOKEN")
	API.StrictParsing = true
	API.DevelMode = true

	if _, err := os.Stat("fixtures/conch-v3"); err == nil {
		return
	} else if os.IsNotExist(err) {
		if API.Token == "" {
			panic("Must supply a token in KOSH_TOKEN")
		}
	} else {
		panic(fmt.Sprintf("%v", err))
	}
}

func setupRecorder(t *testing.T, fixture string) func() {
	t.Helper()
	r, err := recorder.New(fixture)
	if err != nil {
		t.Fatal(err)
	}

	// strip out our authentication headers
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
	})

	// ignore hostnames when fetching from the casset
	r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		iURL, _ := url.Parse(i.URL)
		return r.Method == i.Method && r.URL.Path == iURL.Path
	})

	oldClient := API.HTTP
	API.HTTP = &http.Client{
		Transport: r, // Inject as transport!
	}

	return func() {
		API.HTTP = oldClient
		r.Stop()
	}
}
