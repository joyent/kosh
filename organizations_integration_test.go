package main

import (
	"testing"
)

var org Org

func TestOrganizationAPIIntegration(t *testing.T) {
	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/organizations")
	defer r() // Make sure recorder is stopped once done with it

	t.Run("create", func(t *testing.T) {
		mock := newTestOrganization()
		org = API.Organizations().Create(
			mock.Name,
			mock.Description,
			[]map[string]string{{"email": "conch@example.com"}},
		)
	})

	t.Run("get-all", func(t *testing.T) {
		_ = API.Organizations().GetAll()
	})

	t.Run("get-one", func(t *testing.T) {
		_ = API.Organizations().Get(org.ID)
	})

	t.Run("delete", func(t *testing.T) {
		API.Organizations().Delete(org.ID)
	})
}
