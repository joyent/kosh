package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomsCreateFromStruct(t *testing.T) {
	spy := requestSpy{}
	var got Room

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/DatacenterRoomCreate")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
		json.NewEncoder(w).Encode(got)
	}))
	defer server.Close()

	API.URL = server.URL

	r := API.Rooms()

	want := r.CreateFromStruct(newTestRoom())

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/room")
	assertRequestMethod(t, spy.requestMethod, "POST")
	assert.Equal(t, got, want)
}
