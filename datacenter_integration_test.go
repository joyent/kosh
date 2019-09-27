package main

import (
	"testing"
)

var dc Datacenter

func TestDatacenterAPIIntegration(t *testing.T) {
	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/datacenter")
	defer r() // Make sure recorder is stopped once done with it

	t.Run("create", func(t *testing.T) {

		dc = API.Datacenters().Create(
			"Atlantis",
			"Asguardians",
			"Hala",
			"",
		)
		t.Logf("created %v", dc)
	})

	t.Run("get-all", func(t *testing.T) {
		h := API.Datacenters().GetAll()
		t.Logf("got %v", h)
	})

	t.Run("get-one", func(t *testing.T) {
		h := API.Datacenters().Get(dc.ID)
		t.Logf("got %v", h)
	})
}
