package conch_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/v3/conch"
	"github.com/stretchr/testify/assert"
)

func TestDeviceReportRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/device_report/",
			Method: "POST",
			Do: func(c *conch.Client) {
				j, _ := json.Marshal("{}")
				_ = c.SendDeviceReport(bytes.NewBuffer(j))
			},
		},
		{
			URL:    "/device_report/",
			Method: "POST",
			Do: func(c *conch.Client) {
				j, _ := json.Marshal("{}")
				_ = c.ValidateDeviceReport(bytes.NewBuffer(j))
			},
		},
		{
			URL:    "/device_report/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetDeviceReport("foo") },
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s %s", test.Method, test.URL)
		t.Run(name, func(t *testing.T) {
			seen := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.URL, r.URL.String())
				assert.Equal(t, test.Method, r.Method)
				seen = true

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()
			test.Do(conch.New(conch.API(ts.URL)))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}
