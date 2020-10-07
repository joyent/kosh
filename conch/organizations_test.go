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

func TestOrganizationRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/organization/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetOrganizations() },
		},
		{
			URL:    "/organization/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.CreateOrganization(types.OrganizationCreate{}) },
		},
		{
			URL:    "/organization/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetOrganizationByID("foo") },
		},
		{
			URL:    "/organization/foo/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateOrganization("foo", types.OrganizationUpdate{})
			},
		},
		{
			URL:    "/organization/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteOrganization("foo") },
		},
		{
			URL:    "/organization/foo/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddOrganizationUser("foo", types.OrganizationAddUser{})
			},
		},
		{
			URL:    "/organization/foo/user/bar/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteOrganizationUser("foo", "bar") },
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
