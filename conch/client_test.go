package conch_test

import (
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/sling"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/joyent/kosh/conch"
)

type logger struct{}

func (*logger) Debug(...interface{}) {}

func NewTestClient(fixture string) *conch.Client {
	api := os.Getenv("KOSH_URL")
	token := os.Getenv("KOSH_TOKEN")

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

	s := sling.New().Base(api).Client(&http.Client{Transport: r}).Set("Authorization", "Bearer "+token)

	return &conch.Client{Sling: s, Logger: &logger{}}
}
