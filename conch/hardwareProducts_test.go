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

func TestHardwareProductRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/hardware_product/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetHardwareProducts() },
		},
		{
			URL:    "/hardware_product/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.CreateHardwareProduct(types.HardwareProductCreate{})
			},
		},
		{
			URL:    "/hardware_product/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetHardwareProductByID("foo") },
		},
		{
			URL:    "/hardware_product/foo/",
			Method: "POST",
			Do: func(c *conch.Client) {
				_ = c.UpdateHardwareProduct("foo", types.HardwareProductUpdate{})
			},
		},
		{
			URL:    "/hardware_product/00000000-0000-0000-0000-000000000000/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteHardwareProduct(types.UUID{}) },
		},
		{
			URL:    "/hardware_product/foo/specification?path=%5B%2Fbar%5D%2F",
			Method: "PUT",
			Do: func(c *conch.Client) {
				_ = c.UpdateHardwareProductSpecification("foo", "/bar", types.HardwareProductSpecification{})
			},
		},
		{
			URL:    "/hardware_product/foo/specification?path=%5B%2Fbar%5D%2F",
			Method: "DELETE",
			Do: func(c *conch.Client) {
				_ = c.DeleteHardwareProductSpecification("foo", "/bar")
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
