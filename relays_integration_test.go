package main

import (
	"testing"
)

var testRelay Relay

const serialNumber = "sAWCXAbDHumkCUsrvQpvjFJwv"

func TestRelaysAPIIntegration(t *testing.T) {
	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/relays")
	defer r() // Make sure recorder is stopped once done with it

	t.Run("register relay", func(t *testing.T) {
		mock := newTestRelay()
		testRelay = API.Relays().Register(
			serialNumber,
			mock.Version,
			mock.IpAddr,
			mock.Name,
			mock.SshPort,
		)
	})

	t.Run("get all relays", func(t *testing.T) {
		list := API.Relays().GetAll()
		t.Logf("got %v", list)
	})

	t.Run("get one relay", func(t *testing.T) {
		list := API.Relays().Get(testRelay.SerialNumber)
		t.Logf("got %v", list)
	})

	t.Run("remove a relay", func(t *testing.T) {
		API.Relays().Delete(testRelay.SerialNumber)
	})
}
