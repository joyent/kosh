package main

import (
	"testing"
	"time"

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

type TestOrg struct {
	Name        string
	Description string
	Created     time.Time
	Admins      []map[string]string
}

func newTestOrg(t *testing.T) TestOrg {
	t.Helper()

	org := TestOrg{}
	err := faker.FakeData(&org)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("using fake org: %+v", org)
	return org
}
