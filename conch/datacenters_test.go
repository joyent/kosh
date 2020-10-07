package conch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
	"github.com/stretchr/testify/assert"
)

func TestDatacenterRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/dc/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDatacenters() },
		},
		{
			URL:    "/dc/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.CreateDatacenter(types.DatacenterCreate{}) },
		},
		{
			URL:    "/dc/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDatacenterByID("foo") },
		},
		{
			URL:    "/dc/foo/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.UpdateDatacenter("foo", types.DatacenterUpdate{}) },
		},
		{
			URL:    "/dc/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteDatacenter("foo") },
		},
		{
			URL:    "/dc/foo/rooms/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDatacenterRooms("foo") },
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
			test.Do(conch.New(ts.URL, "token", &logger{}))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}
