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

func TestRackRole(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/rack_role/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetAllRackRoles() },
		},
		{
			URL:    "/rack_role/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.CreateRackRole(types.RackRoleCreate{}) },
		},
		{
			URL:    "/rack_role/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRackRoleByName("foo") },
		},
		{
			URL:    "/rack_role/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRackRoleByID(types.UUID{}) },
		},
		{
			URL:    "/rack_role/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.UpdateRackRole(types.UUID{}, types.RackRoleUpdate{}) },
		},
		{
			URL:    "/rack_role/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteRackRole(types.UUID{}) },
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
			test.Do(conch.New(newConfig(ts.URL)))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}
