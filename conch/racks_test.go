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

func TestRacks(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/rack/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.CreateRack(types.RackCreate{}) },
		},
		{
			URL:    "/rack/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRackByName("foo") },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRackByID(types.UUID{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.UpdateRack(types.UUID{}, types.RackUpdate{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteRack(types.UUID{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/layout/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRackLayout(types.UUID{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/layout/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.UpdateRackLayout(types.UUID{}, []types.RackLayoutUpdate{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRackAssignments(types.UUID{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.UpdateRackAssignments(types.UUID{}, types.RackAssignmentUpdates{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteRackAssignments(types.UUID{}, types.RackAssignmentDeletes{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/phase/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRackPhase(types.UUID{}, types.RackPhase{}, false)
			},
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/phase/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRackPhase(types.UUID{}, types.RackPhase{}, true)
			},
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/links/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRackLinks(types.UUID{}, types.RackLinks{})
			},
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/links/",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				c.DeleteRackLinks(types.UUID{}, types.RackLinks{})
			},
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/layout/01/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetSingleRackLayoutByRU(types.UUID{}, "01") },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/layout/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetSingleRackLayoutByID(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/layout/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateSingleRackLayout(
					types.UUID{},
					types.UUID{},
					types.RackLayoutUpdate{},
				)
			},
		},
		{
			URL:    "/rack/00000000-0000-0000-0000-000000000000/layout/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteSingleRackLayout(types.UUID{}, types.UUID{}) },
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
