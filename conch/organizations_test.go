package conch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/v3/conch"
	"github.com/joyent/kosh/v3/conch/types"
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
			Do:     func(c *conch.Client) { _ = c.GetAllOrganizations() },
		},
		{
			URL:    "/organization/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.CreateOrganization(types.OrganizationCreate{}) },
		},
		{
			URL:    "/organization/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetOrganizationByName("foo") },
		},
		{
			URL:    "/organization/foo/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateOrganization("foo", types.OrganizationUpdate{})
			},
		},
		{
			URL:    "/organization/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteOrganization(types.UUID{}) },
		},
		{
			URL:    "/organization/00000000-0000-0000-0000-000000000000/user/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddOrganizationUser(types.UUID{}, types.OrganizationAddUser{}, false)
			},
		},
		{
			URL:    "/organization/00000000-0000-0000-0000-000000000000/user/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddOrganizationUser(types.UUID{}, types.OrganizationAddUser{}, true)
			},
		},
		{
			URL:    "/organization/00000000-0000-0000-0000-000000000000/user/bar/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteOrganizationUser(types.UUID{}, "bar", false) },
		},
		{
			URL:    "/organization/00000000-0000-0000-0000-000000000000/user/bar/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteOrganizationUser(types.UUID{}, "bar", true) },
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
