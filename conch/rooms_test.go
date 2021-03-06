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

func TestRooms(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/room/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetAllRooms() },
		},
		{
			URL:    "/room/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.CreateRoom(types.DatacenterRoomCreate{}) },
		},
		{
			URL:    "/room/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRoomByAlias("foo") },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRoomByID(types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoom(types.UUID{}, types.DatacenterRoomUpdate{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteRoom(types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetAllRoomRacks(types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRoomRackByName(types.UUID{}, "foo") },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRoomRackByID(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoomRack(types.UUID{}, types.UUID{}, types.RackUpdate{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteRoomRack(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetRoomRackLayout(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoomRackLayout(types.UUID{}, types.UUID{}, types.RackLayoutUpdate{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "GET",
			Do: func(c *conch.Client) {
				c.GetRoomRackAssignments(types.UUID{}, types.UUID{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoomRackAssignments(
					types.UUID{},
					types.UUID{},
					types.RackAssignmentUpdates{},
				)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				c.DeleteRoomRackAssignments(
					types.UUID{},
					types.UUID{},
					types.RackAssignmentDeletes{},
				)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/phase/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoomRackPhase(types.UUID{}, types.UUID{}, types.RackPhase{}, false)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/phase/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoomRackPhase(types.UUID{}, types.UUID{}, types.RackPhase{}, true)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/links/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateRoomRackLinks(types.UUID{}, types.UUID{}, types.RackLinks{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/links/",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				c.DeleteRoomRackLinks(types.UUID{}, types.UUID{}, types.RackLinks{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/01/",
			Method: "GET",
			Do: func(c *conch.Client) {
				c.GetSingleRoomRackLayoutByRU(
					types.UUID{},
					types.UUID{},
					"01",
				)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do: func(c *conch.Client) {
				c.GetSingleRoomRackLayoutByID(
					types.UUID{},
					types.UUID{},
					types.UUID{},
				)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.UpdateSingleRoomRackLayout(
					types.UUID{},
					types.UUID{},
					types.UUID{},
					types.RackLayoutUpdate{},
				)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				c.DeleteSingleRoomRackLayout(
					types.UUID{},
					types.UUID{},
					types.UUID{},
				)
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
