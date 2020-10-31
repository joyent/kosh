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

func TestBuilds(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/build/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetAllBuilds() },
		},
		{
			URL:    "/build/",
			Method: "POST",
			Do: func(c *conch.Client) {
				build := types.BuildCreate{}
				_ = c.CreateBuild(build)
			},
		},
		{
			URL:    "/build/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetBuildByName("foo") },
		},
		{
			URL:    "/build/foo/",
			Method: "POST",
			Do: func(c *conch.Client) {
				update := types.BuildUpdate{}
				_ = c.UpdateBuild("foo", update)
			},
		},
		{
			URL:    "/build/foo/user/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetBuildUsers("foo") },
		},
		{
			URL:    "/build/foo/user/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.AddBuildUser("foo", types.BuildAddUser{}, false) },
		},
		{
			URL:    "/build/foo/user/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.AddBuildUser("foo", types.BuildAddUser{}, true) },
		},

		{
			URL:    "/build/foo/user/alice/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteBuildUser("foo", "alice", false) },
		},
		{
			URL:    "/build/foo/user/alice/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteBuildUser("foo", "alice", true) },
		},

		{
			URL:    "/build/foo/organization/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetAllBuildOrganizations("foo") },
		},
		{
			URL:    "/build/foo/organization/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddBuildOrganization("foo", types.BuildAddOrganization{}, false)
			},
		},
		{
			URL:    "/build/foo/organization/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddBuildOrganization("foo", types.BuildAddOrganization{}, true)
			},
		},
		{
			URL:    "/build/foo/organization/lemmings/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteBuildOrganization("foo", "lemmings", false) },
		},
		{
			URL:    "/build/foo/organization/lemmings/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteBuildOrganization("foo", "lemmings", true) },
		},

		{
			URL:    "/build/foo/device/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetAllBuildDevices("foo") },
		},
		{
			URL:    "/build/foo/device/pxe/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetBuildDevicesPXE("foo") },
		},
		{
			URL:    "/build/foo/device/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddNewBuildDevice("foo", types.BuildCreateDevices{})
			},
		},
		{
			URL:    "/build/00000000-0000-0000-0000-000000000000/device/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddBuildDeviceByID(types.UUID{}, types.UUID{})
			},
		},
		{
			URL:    "/build/00000000-0000-0000-0000-000000000000/device/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddBuildDeviceByID(types.UUID{}, types.UUID{})
			},
		},
		{
			URL:    "/build/00000000-0000-0000-0000-000000000000/device/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				_ = c.DeleteBuildDeviceByID(types.UUID{}, types.UUID{})
			},
		},
		{
			URL:    "/build/foo/rack/",
			Method: "GET",
			Do: func(c *conch.Client) {
				_ = c.GetBuildRacks("foo")
			},
		},
		{
			URL:    "/build/foo/rack/DEADBEEF/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.AddBuildRackByID("foo", "DEADBEEF")
			},
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
