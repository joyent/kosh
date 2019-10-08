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

	"github.com/gofrs/uuid"
)

func newOrgList() Orgs {
	return Orgs{
		newStubOrg("A", "An Example Organization"),
		newStubOrg("B", "Another Example Organization"),
	}
}

func TestOrganizationsGet(t *testing.T) {
	orgList := newOrgList()
	spy := requestSpy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		orgID := r.URL.Path[len("/organization/"):]
		for _, org := range orgList {
			if org.ID.String() == orgID {
				json.NewEncoder(w).Encode(org)
				return
			}
		}
	}))
	defer server.Close()

	API.URL = server.URL
	o := API.Organizations()

	for i, org := range orgList {
		t.Run(org.Name, func(t *testing.T) {
			got := o.Get(org.ID)

			assertRequestMethod(t, spy.requestMethod, "GET")
			assertRequestCount(t, spy.requestCount, i+1)
			assertRequestPath(t, spy.requestPath, fmt.Sprintf("/organization/%s", org.ID))
			assertData(t, got, org)
		})
	}

}

func TestOrganizationsGetByName(t *testing.T) {
	orgList := newOrgList()
	spy := requestSpy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		orgName := r.URL.Path[len("/organization/"):]
		for _, org := range orgList {
			if org.Name == orgName {
				json.NewEncoder(w).Encode(org)
				return
			}
		}
	}))
	defer server.Close()

	API.URL = server.URL
	o := API.Organizations()

	for i, org := range orgList {
		t.Run(org.Name, func(t *testing.T) {
			got := o.GetByName(org.Name)

			assertRequestMethod(t, spy.requestMethod, "GET")
			assertRequestCount(t, spy.requestCount, i+1)
			assertRequestPath(t, spy.requestPath, fmt.Sprintf("/organization/%s", org.Name))
			assertData(t, got, org)
		})
	}

}

func TestOrganizationsGetAll(t *testing.T) {
	orgList := newOrgList()
	spy := requestSpy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(orgList)
	}))
	defer server.Close()

	API.URL = server.URL

	o := API.Organizations()

	got := o.GetAll()

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/organization")
	assertData(t, got, orgList)
}

func TestOrganizationsCreate(t *testing.T) {
	spy := requestSpy{}
	var got Org

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/OrganizationCreate")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
		json.NewEncoder(w).Encode(got)
	}))
	defer server.Close()

	API.URL = server.URL

	o := API.Organizations()

	want := o.Create("Z", "Z Org", []map[string]string{{"email": "admin@example.com"}})

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/organization")
	assertRequestMethod(t, spy.requestMethod, "POST")
	if got.Name != "Z" {
		t.Errorf("Invalid requestBody. Got %v expected something with name Z", got)
	}

	// now let's check what we made in the server is what we got in the client
	assertData(t, got, want)
}

func TestOrganizationsDelete(t *testing.T) {
	orgList := newOrgList()
	org := orgList[0]
	spy := requestSpy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
	}))
	defer server.Close()

	API.URL = server.URL

	o := API.Organizations()

	o.Delete(org.ID)

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/organization/%s", org.ID))
}

func TestOrganizationsGetUsers(t *testing.T) {
	spy := requestSpy{}
	org := newStubOrg("A", "An Organization")
	userList := []OrganizationUser{
		newOrgUser("Timmy"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(userList)
	}))
	defer server.Close()

	API.URL = server.URL

	o := API.Organizations()

	got := o.GetUsers(org.ID)

	assertRequestMethod(t, spy.requestMethod, "GET")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/organization/%s/user", org.ID))
	assertData(t, got, userList)
}

func TestOrganizationsAddUser(t *testing.T) {
	spy := requestSpy{}
	org := newStubOrg("A", "An Organization")
	user := newOrgUser("Timmy")
	var got OrganizationUser

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/OrganizationAddUser")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
	}))
	defer server.Close()

	API.URL = server.URL

	o := API.Organizations()

	o.AddUser(org.ID, user.Email, "admin", false)

	assertRequestMethod(t, spy.requestMethod, "POST")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/organization/%s/user", org.ID))
}

func TestOrganizationsRemoveUser(t *testing.T) {
	spy := requestSpy{}
	user := newOrgUser("Timmy")
	org := newStubOrg("A", "An Organization")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
	}))
	defer server.Close()

	API.URL = server.URL

	o := API.Organizations()

	o.RemoveUser(org.ID, user.Email, false)

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/organization/%s/user/%s", org.ID, user.Email))
}

/*
	Helpers
*/

type requestSpy struct {
	requestCount  int
	requestPath   string
	requestMethod string
}

func (rs *requestSpy) onRequest(r *http.Request) {
	rs.requestCount++
	rs.requestMethod = r.Method
	rs.requestPath = r.URL.Path
}

func assertRequestMethod(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Request method wrong, got %s wanted %s", got, want)
	}

}

func assertRequestPath(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Request path wrong, got %s wanted %s", got, want)
	}
}

func newStubOrg(name, description string) Org {
	uuid, _ := uuid.NewV4()

	return Org{
		ID:          uuid,
		Name:        name,
		Description: description,
		//Created:     time.Now(),
		Admins:     DetailedUsers{},
		Workspaces: WorkspaceAndRoles{},
	}
}

func newOrgUser(name string) OrganizationUser {
	uuid, _ := uuid.NewV4()
	return OrganizationUser{
		ID:    uuid,
		Name:  name,
		Email: fmt.Sprintf("%s@example.com", name),
		Role:  "",
	}
}

func assertRequestCount(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Wrong number of requests, got %d wanted %d", got, want)
	}
}

func assertData(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got wrong results, got %v wanted %v", got, want)
	}
}
