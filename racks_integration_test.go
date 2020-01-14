package main

import (
	"testing"
)

const TestRack = "7632cbbc"

func TestRacksAPIIntegration(t *testing.T) {
	defer errorHandler()

	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/racks")
	defer r() // Make sure recorder is stopped once done with it

	f := newFixture()
	f.setupRackRole()
	f.setupRoom()
	f.setupHardwareProducts()
	f.setupBuild()
	defer f.reset()

	var testRack Rack
	t.Run("create a new Rack", func(t *testing.T) {

		mock := newTestRack()
		testRack = API.Racks().Create(
			mock.Name,
			f.room.ID,
			f.role.ID,
			"integration",
			f.build.ID,
		)
	})

	/*	t.Run("fetch all Racks", func(t *testing.T) {
			defer errorHandler()
			list := API.Racks().GetAll()
			t.Logf("got %v", list)
		})
	*/
	t.Run("fetch a single rack", func(t *testing.T) {
		defer errorHandler()
		list := API.Racks().Get(testRack.ID)
		t.Logf("got %v", list)
	})

	t.Run("create a new Rack Layout", func(t *testing.T) {
		defer errorHandler()

		mrl := RackLayoutUpdates{
			{
				RU:        1,
				ProductID: f.serverProduct.ID,
			},
			{
				RU:        1 + f.serverProduct.RackUnitSize,
				ProductID: f.switchProduct.ID,
			},
		}
		_ = API.Racks().CreateLayout(
			testRack.ID,
			mrl,
		)
	})
	/*
		t.Run("create a new Rack Assignment", func(t *testing.T) {
			mock := newTestRackAssignmentUpdates()
			_ = API.Racks().UpdateAssignments(
				testRack.ID,
				mock,
			)
		})
	*/
	t.Run("remove a Rack", func(t *testing.T) {
		defer errorHandler()

		for _, row := range API.Racks().Layouts(testRack.ID) {
			API.Racks().DeleteLayoutSlot(row.ID)
		}

		API.Racks().Delete(testRack.ID)
	})

	// other tests
	t.Run("create a new Rack from struct", func(t *testing.T) {
		defer errorHandler()
		mock := newTestRack()
		mock.RoomID = f.room.ID
		mock.RoleID = f.role.ID
		mock.BuildID = f.build.ID
		mock.Phase = "integration"
		testRack := API.Racks().CreateFromStruct(mock)
		API.Racks().Delete(testRack.ID)
	})

}
