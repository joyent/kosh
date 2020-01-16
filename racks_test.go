package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRackAssignments(t *testing.T) {

	t.Run("update rack assignments", func(t *testing.T) {
		spy := requestSpy{}
		var got RackAssignments

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spy.onRequest(r)
			if r.Method == "POST" {
				body, _ := ioutil.ReadAll(r.Body)
				assertJSONSchema(t, body, "request/RackAssignmentUpdates")
				json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
			} else {
				json.NewEncoder(w).Encode(got)
			}
		}))

		defer server.Close()

		API.URL = server.URL

		newAssignments := newTestRackAssignmentUpdates()
		rack := newTestRack()
		_ = API.Racks().UpdateAssignments(rack.ID, newAssignments)

		assertRequestCount(t, spy.requestCount, 2)
		assertRequestPath(t, spy.requestPath, fmt.Sprintf("/rack/%s/assignment", rack.ID))
	})

}

func TestRackLayout(t *testing.T) {

	t.Run("create rack layout", func(t *testing.T) {
		//defer errorHandler()
		spy := requestSpy{}
		var got RackLayout

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spy.onRequest(r)
			if r.Method == "POST" {
				body, _ := ioutil.ReadAll(r.Body)
				assertJSONSchema(t, body, "request/RackLayoutCreate")
				json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
				json.NewEncoder(w).Encode(got)
			} else {
				json.NewEncoder(w).Encode([]string{})
			}
		}))

		defer server.Close()

		API.URL = server.URL
		svp := newTestHardwareProduct()
		mrl := RackLayoutUpdates{
			{
				RU:        1,
				ProductID: svp.ID,
			},
		}
		rack := newTestRack()
		_ = API.Racks().CreateLayout(rack.ID, mrl)

		assertRequestCount(t, spy.requestCount, 3)
		assertRequestPath(t, spy.requestPath, fmt.Sprintf("/rack/%s/layout", rack.ID))
	})

}
