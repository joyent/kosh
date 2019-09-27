package main

import (
	"testing"
)

func TestWorkspaceAPIIntergration(t *testing.T) {

	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/workspace")
	defer r() // Make sure recorder is stopped once done with it

	t.Run("get-all", func(t *testing.T) {
		list := API.Workspaces().GetAll()
		t.Logf("got %v", list)
	})

	t.Run("get-by-name", func(t *testing.T) {
		_ = API.Workspaces().GetByName("GLOBAL")
	})

	t.Run("get-users", func(t *testing.T) {
		_ = API.Workspaces().GetByName("GLOBAL")
	})
}
