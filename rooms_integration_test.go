package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomsAPIIntegration(t *testing.T) {
	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/rooms")
	defer r() // Make sure recorder is stopped once done with it

	dc := API.Datacenters().CreateFromStruct(newTestDatacenter())
	defer API.Datacenters().Delete(dc.ID)

	var testRoom Room
	t.Run("Create a New Room", func(t *testing.T) {
		want := newTestRoom()
		testRoom = API.Rooms().Create(
			dc.ID,
			want.AZ,
			want.Alias,
			want.VendorName,
		)
		assert.NotNil(t, testRoom.ID)
		// TODO write a functional test here -- perigrin
	})

	t.Run("Get all rooms", func(t *testing.T) {
		defer errorHandler()
		list := API.Rooms().GetAll()
		t.Logf("got %v", list)
	})

	t.Run("Get one room", func(t *testing.T) {
		defer errorHandler()
		list := API.Rooms().Get(testRoom.ID)
		t.Logf("got %v", list)
	})

	t.Run("List all racks in a room", func(t *testing.T) {
		defer errorHandler()
		list := API.Rooms().Racks(testRoom.ID)
		t.Logf("got %v", list)
	})

	t.Run("Remove a room", func(t *testing.T) {
		defer errorHandler()
		API.Rooms().Delete(testRoom.ID)
	})
}
