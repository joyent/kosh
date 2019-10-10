package main

import (
	"testing"

	"github.com/bxcodec/faker"
)

var org Org

func TestIntegrationOrganizationsCreate(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/organizations/create")
	defer r() // Make sure recorder is stopped once done with it
	fake := newTestOrg(t)
	org = API.Organizations().Create(
		fake.Name,
		fake.Description,
		[]map[string]string{{"email": "conch@example.com"}},
	)
}

func TestIntegrationOrganizatiosGetAll(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/organizations/get-all")
	defer r() // Make sure recorder is stopped once done with it

	_ = API.Organizations().GetAll()
}

func TestIntegrationOrganizatiosGet(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/organizations/get")
	defer r() // Make sure recorder is stopped once done with it

	_ = API.Organizations().Get(org.ID)
}

func TestIntegrationOrganizatiosDelete(t *testing.T) {
	setupAPIClient()
	r := setupRecorder(t, "fixtures/conch-v3/organizations/delete")
	defer r() // Make sure recorder is stopped once done with it

	API.Organizations().Delete(org.ID)
}

// ---

func newTestOrg(t *testing.T) (org Org) {
	t.Helper()

	err := faker.FakeData(&org)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return org
}
