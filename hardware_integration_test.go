package main

import (
	"testing"
)

// the validations are currently hardcoded in every instance of conch
// fixing that is on the slate for v3.1

func TestHarwareProductAPIIntegration(t *testing.T) {
	defer errorHandler()

	setupAPIClient()

	r := setupRecorder("fixtures/conch-v3/hardware")
	defer r() // Make sure recorder is stopped once done with it

	f := newFixture()
	f.setupHardwareVendor()
	f.setupValidationPlan()
	defer f.reset()

	var testHardwareProduct HardwareProduct

	t.Run("create", func(t *testing.T) {

		mock := newTestHardwareProduct()
		testHardwareProduct = API.Hardware().Create(
			mock.Name,
			mock.Alias,
			f.hardwareVendor.ID,
			mock.SKU,
			mock.RackUnitSize,
			f.validationPlan.ID,
			newTestHardwareProductProfile(),
		)
	})

	t.Run("get-product-by-name", func(t *testing.T) {
		_ = API.Hardware().GetProductByName(testHardwareProduct.Name)
	})

	t.Run("delete", func(t *testing.T) {
		API.Hardware().Delete(testHardwareProduct.ID)
	})
}
