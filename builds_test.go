package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/gofrs/uuid"
)

func TestBuildsGetAll(t *testing.T) {
	spy := requestSpy{}
	buildList := newTestBuildList()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(buildList)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	got := b.GetAll()

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/build")
	assertData(t, got, buildList)
}

func TestBuildsGet(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(build)
	}))
	defer server.Close()

	API.URL = server.URL
	b := API.Builds()

	got := b.Get(build.ID)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s", build.ID))
	assertData(t, got, build)
}

func TestBuildsGetByName(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(build)
	}))
	defer server.Close()

	API.URL = server.URL
	b := API.Builds()

	got := b.GetByName(build.Name)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s", build.Name))
	assertData(t, got, build)
}

func TestBuildsCreate(t *testing.T) {
	spy := requestSpy{}
	var got Build

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/BuildCreate")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
		json.NewEncoder(w).Encode(got)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	want := b.Create("Z", "Z Build", []map[string]string{{"email": "admin@example.com"}})

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/build")
	assertRequestMethod(t, spy.requestMethod, "POST")
	if got.Name != "Z" {
		t.Errorf("Invalid requestBody. Got %v expected something with name Z", got)
	}

	// now let's check what we made in the server is what we got in the client
	assertData(t, got, want)
}

func TestBuildsGetUsers(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	list := newTestUserAndRoles()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(list)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	got := b.GetUsers(build.ID)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/user", build.ID))
	assertData(t, got, list)
}

func TestBuildsAddUser(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	user := newTestUser()
	var got UserAndRole

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/BuildAddUser")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	b.AddUser(build.ID, user.Email, "admin", false)

	assertRequestMethod(t, spy.requestMethod, "POST")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/user", build.ID))
}

func TestBuildsRemoveUser(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	user := newTestUser()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()
	b.RemoveUser(build.ID, user.Email, false)

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/user/%s", build.ID, user.Email))
}

func TestBuildsGetOrgs(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	list := newTestOrgAndRoles()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(list)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	got := b.GetOrgs(build.ID)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/organization", build.ID))
	assertData(t, got, list)
}

func TestBuildsAddOrg(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	org := newTestOrg(t)
	var got OrgAndRole

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/BuildAddOrganization")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	b.AddOrg(build.ID, org.ID.String(), "admin", false)

	assertRequestMethod(t, spy.requestMethod, "POST")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/organization", build.ID))
}

func TestBuildsRemoveOrg(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	org := newTestOrg(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()
	b.RemoveOrg(build.ID, org.ID.String(), false)

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/organization/%s", build.ID, org.ID))
}

func TestBuildsGetDevices(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	list := newTestDeviceList()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(list)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	got := b.GetDevices(build.ID)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/device", build.ID))
	assertData(t, got, list)
}

func TestBuildsAddDevice(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	device := newTestDevice()
	var got Device

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	b.AddDevice(build.ID, device.ID.String())

	assertRequestMethod(t, spy.requestMethod, "POST")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/device/%s", build.ID, device.ID))
}

func TestBuildsRemoveDevice(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	device := newTestDevice()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()
	b.RemoveDevice(build.ID, device.ID.String())

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/device/%s", build.ID, device.ID))
}

func TestBuildsGetRacks(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	list := newTestRackList()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(list)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	got := b.GetRacks(build.ID)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/rack", build.ID))
	assertData(t, got, list)
}

func TestBuildsAddRack(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	rack := newTestRack()
	var got Rack

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()

	b.AddRack(build.ID, rack.ID.String())

	assertRequestMethod(t, spy.requestMethod, "POST")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/rack/%s", build.ID, rack.ID))
}

func TestBuildsRemoveRack(t *testing.T) {
	spy := requestSpy{}
	build := newTestBuild()
	rack := newTestRack()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
	}))
	defer server.Close()

	API.URL = server.URL

	b := API.Builds()
	b.RemoveRack(build.ID, rack.ID.String())

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/build/%s/rack/%s", build.ID, rack.ID))
}

// ----

func setupFaker() {
	faker.AddProvider("uuid", func(v reflect.Value) (interface{}, error) {
		return uuid.NewV4()
	})
}

func newTestBuildList() (list BuildList) {
	setupFaker()
	_ = faker.FakeData(&list)
	return
}

func newTestBuild() (build Build) {
	setupFaker()
	_ = faker.FakeData(&build)
	return
}

func newTestUser() (user UserAndRole) {
	setupFaker()
	_ = faker.FakeData(&user)
	return
}

func newTestUserAndRoles() (list UserAndRoles) {
	setupFaker()
	_ = faker.FakeData(&list)
	return
}

func newTestOrgAndRoles() (list OrgAndRoles) {
	setupFaker()
	_ = faker.FakeData(&list)
	return
}

func newTestDevice() (device Device) {
	setupFaker()
	_ = faker.FakeData(&device)
	return
}

func newTestDeviceList() (list DeviceList) {
	setupFaker()
	_ = faker.FakeData(&list)
	return
}

func newTestRack() (rack Rack) {
	setupFaker()
	_ = faker.FakeData(&rack)
	return
}

func newTestRackList() (list RackList) {
	setupFaker()
	_ = faker.FakeData(&list)
	return
}
