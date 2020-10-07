package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHarwareProductIntegration(t *testing.T) {
	defer errorHandler()

	setupAPIClient()

	r := setupRecorder("fixtures/conch-v3/hardware")
	defer r() // Make sure recorder is stopped once done with it

	f := newFixture()
	f.setupHardwareVendor()
	f.setupValidationPlan()
	defer f.reset()

	t.Run("API", func(t *testing.T) {
		var testHardwareProduct HardwareProduct
		var testHardwareProductSummary HardwareProductSummary

		t.Run("create", func(t *testing.T) {
			defer errorHandler()
			mock := newTestHardwareProduct()
			testHardwareProduct = API.Hardware().Create(
				mock.Name,
				mock.Alias,
				f.hardwareVendor.ID,
				mock.SKU,
				mock.RackUnitSize,
				f.validationPlan.ID,
				mock.Purpose,
				mock.BiosFirmware,
				mock.CpuType,
			)

			testHardwareProductSummary = HardwareProductSummary{
				ID:             testHardwareProduct.ID,
				Name:           testHardwareProduct.Name,
				Alias:          testHardwareProduct.Alias,
				GenerationName: testHardwareProduct.GenerationName,
				SKU:            testHardwareProduct.SKU,
				Created:        testHardwareProduct.Created,
				Updated:        testHardwareProduct.Updated,
			}

			assert.NotNil(t, testHardwareProduct.ID)
		})

		t.Run("get all products", func(t *testing.T) {
			got := API.Hardware().GetAllProducts()
			assert.Equal(t, HardwareProductSummaries{testHardwareProductSummary}, got)
		})

		t.Run("get product by name", func(t *testing.T) {
			defer errorHandler()
			got := API.Hardware().GetProductByName(testHardwareProduct.Name)
			assert.Equal(t, testHardwareProduct, got)
		})

		t.Run("delete product", func(t *testing.T) {
			defer errorHandler()
			API.Hardware().Delete(testHardwareProduct.ID)
			got := API.Hardware().GetAllProducts()
			assert.Equal(t, HardwareProductSummaries{}, got)
		})
	})
	/* TODO: Figure out timezones
	t.Run("cli", func(t *testing.T) {
		mock := newTestHardwareProduct()
		mock.SKU = "test-sku-001"
		mock.Name = "Testy McTesterson"

		cases := []struct {
			name string
			cli  []string
			want string
		}{
			{
				"create",
				[]string{
					"kosh", "hardware", "products", "create",
					"--sku", mock.SKU,
					"--name", mock.Name,
					"--alias", mock.Alias,
					"--vendor", f.hardwareVendor.Name,
					"--rack-unit-size", strconv.Itoa(mock.RackUnitSize),
					"--validation-plan", f.validationPlan.ID.String(),
					"--purpose", mock.Purpose,
					"--bios-firmware", mock.BiosFirmware,
					"--cpu-type", mock.CpuType,
				},
				"\nID: 9ad55ceb-2eb7-4125-a492-3c595277b3e3\nName: Testy McTesterson\nSKU: test-sku-001\n\nCreated: 2020-01-26 19:04:27 +0000 UTC\nUpdated: 2020-01-26 19:04:27 +0000 UTC\n\n",
			},
			{
				"products ls",
				[]string{"kosh", "hardware", "products", "ls"},
				"|    ID    |     SKU      |       NAME        |           ALIAS           | GENERATIONNAME |               CREATED                |               UPDATED                |\n|----------|--------------|-------------------|---------------------------|----------------|--------------------------------------|--------------------------------------|\n| 9ad55ceb | test-sku-001 | Testy McTesterson | ujOmocHFAUuWZILajRAzVkeuO |                | 2020-01-26 19:04:27.567389 +0000 UTC | 2020-01-26 19:04:27.567389 +0000 UTC |\n\n",
			},
			{
				"product SKU get",
				[]string{"kosh", "hardware", "product", mock.SKU, "get"},
				"\nID: 9ad55ceb-2eb7-4125-a492-3c595277b3e3\nName: Testy McTesterson\nSKU: test-sku-001\n\nCreated: 2020-01-26 19:04:27 +0000 UTC\nUpdated: 2020-01-26 19:04:27 +0000 UTC\n\n",
			},
			{
				"product SKU delete",
				[]string{"kosh", "hardware", "product", mock.SKU, "rm"},
				HardwareProductSummaries{}.String() + "\n",
			},
		}
		for _, cas := range cases {
			defer errorHandler()
			t.Run(cas.name, func(t *testing.T) {
				defer errorHandler()
				t.Logf("Testing %+v", cas.cli)
				app := cli.App("kosh", "Command line interface for Conch")
				initHardwareCli(app)
				got := captureOutput(func() { app.Run(cas.cli) })
				assert.Equal(t, cas.want, got)
			})
		}
	})
	*/

}
