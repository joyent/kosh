package main

import (
	"testing"
)

var build Build

func TestIntegrationBuildsCreate(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/builds/create")
	defer r() // Make sure recorder is stopped once done with it
	fake := newTestBuild()
	build = API.Builds().Create(
		fake.Name,
		fake.Description,
		[]map[string]string{{"email": "conch@example.com"}},
	)
}

func TestIntegrationBuildsGetAll(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/builds/get-all")
	defer r() // Make sure recorder is stopped once done with it

	_ = API.Builds().GetAll()
}

func TestIntegrationBuildsGet(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/builds/get")
	defer r() // Make sure recorder is stopped once done with it

	_ = API.Builds().Get(build.ID)
}
