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

func TestHardwareProductCreate(t *testing.T) {
	spy := requestSpy{}
	var got HardwareProduct

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		body, _ := ioutil.ReadAll(r.Body)
		assertJSONSchema(t, body, "request/HardwareProductCreate")
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&got)
		json.NewEncoder(w).Encode(got)
	}))
	defer server.Close()

	API.URL = server.URL

	h := API.Hardware()

	mock := newTestHardwareProduct()
	want := h.Create(
		mock.Name,
		mock.Alias,
		mock.HardwareVendorID,
		mock.SKU,
		mock.RackUnitSize+1,
		mock.ID,
		newTestHardwareProductProfile(),
	)

	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, "/hardware_product")
	assertRequestMethod(t, spy.requestMethod, "POST")
	assertData(t, got, want)
}

func TestHardwareProductDelete(t *testing.T) {
	spy := requestSpy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		w.WriteHeader(204)

	}))
	defer server.Close()

	API.URL = server.URL

	hp := newTestHardwareProduct()
	h := API.Hardware()

	h.Delete(hp.ID)

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/hardware_product/%s", hp.ID))
}

func TestHardwareVendorCreate(t *testing.T) {
	spy := requestSpy{}

	mock := newTestHardwareVendor()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		json.NewEncoder(w).Encode(mock)
	}))
	defer server.Close()

	API.URL = server.URL

	h := API.Hardware()

	_ = h.CreateVendor(mock.Name)

	assertRequestCount(t, spy.requestCount, 2)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/hardware_vendor/%s", mock.Name))
	assertRequestMethod(t, spy.requestMethod, "GET")
}

func TestHardwareVendorDelete(t *testing.T) {
	spy := requestSpy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spy.onRequest(r)
		w.WriteHeader(204)

	}))
	defer server.Close()

	API.URL = server.URL

	hv := newTestHardwareVendor()
	h := API.Hardware()

	h.DeleteVendor(hv.ID.String())

	assertRequestMethod(t, spy.requestMethod, "DELETE")
	assertRequestCount(t, spy.requestCount, 1)
	assertRequestPath(t, spy.requestPath, fmt.Sprintf("/hardware_vendor/%s", hv.ID))
}
