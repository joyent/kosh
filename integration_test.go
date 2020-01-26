package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func setupAPIClient() {
	Version = "3.x-testing"
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
		panic(err)
	}
}

func setupRecorder(fixture string) func() {
	defer errorHandler()
	// TODO: we need to re-think the test fixtures entirely
	r, err := recorder.New(fixture)
	if err != nil {
		panic(err)
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

func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}

}

func captureOutput(f func()) string {
	done := capture()
	f()
	out, err := done()
	if err != nil {
		panic(err)
	}
	return out
}
