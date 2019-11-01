package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRackRoleCreateFromStruct(t *testing.T) {
	spy := requestSpy{}
	var got RackRole

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/RackRoleCreate")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
		json.NewEncoder(w).Encode(got)
	}))
	defer server.Close()

	API.URL = server.URL

	r := API.RackRoles()

	want := r.CreateFromStruct(newTestRackRole())

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/rack_role")
	assertRequestMethod(t, spy.requestMethod, "POST")
	assertData(t, got, want)
}
