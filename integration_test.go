package main

import (
	"net/http"
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

	if API.Token == "" {
		panic("Must supply a token in KOSH_TOKEN")
	}

}

func setupRecorder(t *testing.T, fixture string) func() {
	t.Helper()
	r, err := recorder.New(fixture)
	if err != nil {
		t.Fatal(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
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
