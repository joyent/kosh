package main

import (
	"testing"
)

const TestSetting = "test_setting"

func TestDeviceAPIIntegration(t *testing.T) {
	defer errorHandler()

	setupAPIClient()

	r := setupRecorder("fixtures/conch-v3/device")
	defer r() // Make sure recorder is stopped once done with it

	f := newFixture()
	f.setupRackLayout()
	defer f.reset()

	var testDevice DetailedDevice

	t.Run("create a device by adding it to a build with a sku", func(t *testing.T) {
		defer errorHandler()

		mock := newTestDevice()
		mock.Serial = "AAAAAA"
		API.Builds().CreateDevice(f.build.ID, mock.Serial, f.switchProduct.SKU)
		testDevice = API.Devices().Get(mock.Serial)
	})

	t.Run("set and fetching a device setting", func(t *testing.T) {
		defer errorHandler()

		API.Devices().SetSetting(testDevice.ID.String(), TestSetting, "foo")
		s := API.Devices().Setting(testDevice.ID.String(), TestSetting)
		t.Logf("got %v", s)
	})

	t.Run("fetching all device settings", func(t *testing.T) {
		defer errorHandler()

		s := API.Devices().Settings(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	t.Run("fetch device tags", func(t *testing.T) {

		s := API.Devices().Tags(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	t.Run("set and fetch device location", func(t *testing.T) {
		defer errorHandler()

		_ = API.Racks().UpdateAssignments(
			f.rack.ID,
			RackAssignmentUpdates{{testDevice.ID, "", 1}},
		)
		s := API.Devices().GetLocation(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	t.Run("fetch getting a device", func(t *testing.T) {
		defer errorHandler()

		s := API.Devices().Get(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	t.Run("fetch device validation state", func(t *testing.T) {

		s := API.Devices().ValidationState(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	t.Run("fetch device phase", func(t *testing.T) {

		s := API.Devices().GetPhase(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	// GetInterface
	t.Run("fetch device interface", func(t *testing.T) {
		return // TODO: sort out device report submission

		defer errorHandler()

		s := API.Devices().GetIPMI(testDevice.ID.String())
		t.Logf("got %v", s)
	})

	t.Run("remove device", func(t *testing.T) {
		defer errorHandler()

		//		API.Builds().RemoveDevice(f.build.ID, testDevice.ID.String())
		API.Devices().DeleteLocation(testDevice.ID.String())
	})
}
