package main

import (
	"testing"
)

func TestUserAPIIntergration(t *testing.T) {
	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/user")
	defer r() // Make sure recorder is stopped once done with it

	t.Run("me", func(t *testing.T) {
		_ = API.Users().Me()
	})

	t.Run("me-settings", func(t *testing.T) {
		_ = API.Users().MySettings()
	})
}
