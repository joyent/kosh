package main

import (
	"testing"
)

func TestRackRoleAPIIntegration(t *testing.T) {
	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/racks-roles")
	defer r() // Make sure recorder is stopped once done with it

	var testRackRole RackRole
	t.Run("create a rack role", func(t *testing.T) {
		mockRole := newTestRackRole()
		testRackRole = API.RackRoles().CreateFromStruct(mockRole)
	})

	t.Run("get all rack roles", func(t *testing.T) {

		list := API.RackRoles().GetAll()
		t.Logf("got %v", list)
	})

	t.Run("get a rack role", func(t *testing.T) {

		list := API.RackRoles().Get(testRackRole.ID)
		t.Logf("got %v", list)
	})

	t.Run("get a role by name", func(t *testing.T) {

		list := API.RackRoles().GetByName(testRackRole.Name)
		t.Logf("got %v", list)
	})
}
