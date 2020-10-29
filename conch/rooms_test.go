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
			Do:     func(c *conch.Client) { _ = c.GetAllRooms() },
		},
		{
			URL:    "/room/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.CreateRoom(types.DatacenterRoomCreate{}) },
		},
		{
			URL:    "/room/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRoomByAlias("foo") },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRoomByID(types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateRoom(types.UUID{}, types.DatacenterRoomUpdate{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteRoom(types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetAllRoomRacks(types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRoomRackByName(types.UUID{}, "foo") },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRoomRackByID(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateRoomRack(types.UUID{}, types.UUID{}, types.RackUpdate{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteRoomRack(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRoomRackLayout(types.UUID{}, types.UUID{}) },
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateRoomRackLayout(types.UUID{}, types.UUID{}, types.RackLayoutUpdate{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "GET",
			Do: func(c *conch.Client) {
				_ = c.GetRoomRackAssignments(types.UUID{}, types.UUID{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/assignment/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateRoomRackAssignments(
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
				_ = c.DeleteRoomRackAssignments(
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
				_ = c.UpdateRoomRackPhase(types.UUID{}, types.UUID{}, types.RackPhase{}, false)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/phase/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateRoomRackPhase(types.UUID{}, types.UUID{}, types.RackPhase{}, true)
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/links/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateRoomRackLinks(types.UUID{}, types.UUID{}, types.RackLinks{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/links/",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				_ = c.DeleteRoomRackLinks(types.UUID{}, types.UUID{}, types.RackLinks{})
			},
		},
		{
			URL:    "/room/00000000-0000-0000-0000-000000000000/rack/00000000-0000-0000-0000-000000000000/layout/01/",
			Method: "GET",
			Do: func(c *conch.Client) {
				_ = c.GetSingleRoomRackLayoutByRU(
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
				_ = c.GetSingleRoomRackLayoutByID(
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
				_ = c.UpdateSingleRoomRackLayout(
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
				_ = c.DeleteSingleRoomRackLayout(
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
