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

func TestDevices(t *testing.T) {

	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/device?foo=bar",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.FindDevicesBySetting("foo", "bar") },
		},
		{
			URL:    "/device?tag_foo=bar",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.FindDevicesByTag("foo", "bar") },
		},
		{
			URL:    "/device?hostname=bar",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.FindDevicesByField("hostname", "bar") },
		},
		{
			URL:    "/device/DEADBEEF/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceById("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/pxe/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDevicePXE("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/phase/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDevicePhase("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/sku/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceSKU("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/asset_tag/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceAssetTag("DEADBEEF", "123456") },
		},
		{
			URL:    "/device/DEADBEEF/validated/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceValidated("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/phase/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDevicePhase("DEADBEEF", "DEAD") },
		},
		{
			URL:    "/device/DEADBEEF/links/",
			Method: "POST",
			Do: func(c *conch.Client) {
				links := types.NewDeviceLinks(
					"http://example.com",
					"http://different.example.com",
				)
				_ = c.SetDeviceLinks("DEADBEEF", links)
			},
		},
		{
			URL:    "/device/DEADBEEF/links/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteDeviceLinks("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/sku/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceSKU("DEADBEEF", "123456") },
		},
		{
			URL:    "/device/DEADBEEF/build/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceBuild("DEADBEEF", "123456") },
		},
		{
			URL:    "/device/DEADBEEF/location/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceLocation("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/location/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceLocation("DEADBEEF", "elsewhere") },
		},
		{
			URL:    "/device/DEADBEEF/location/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteDeviceLocation("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/settings/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceSettings("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/settings/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceSettingByName("DEADBEEF", "foo") },
		},
		{
			URL:    "/device/DEADBEEF/settings/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceTags("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/settings/tag_foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceTagByName("DEADBEEF", "foo") },
		},
		{
			URL:    "/device/DEADBEEF/settings/tag_foo/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceTag("DEADBEEF", "foo", "bar") },
		},
		{
			URL:    "/device/DEADBEEF/settings/tag_foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteDeviceTag("DEADBEEF", "foo") },
		},
		{
			URL:    "/device/DEADBEEF/settings/foo/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.SetDeviceSetting("DEADBEEF", "foo", "bar") },
		},
		{
			URL:    "/device/DEADBEEF/settings/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteDeviceSetting("DEADBEEF", "foo") },
		},
		{
			URL:    "/device/DEADBEEF/validation/0D15EA5E/",
			Method: "POST",
			Do: func(c *conch.Client) {
				report := types.DeviceReport{}
				_ = c.RunValidationForDevice("DEADBEEF", "0D15EA5E", report)
			},
		},
		{
			URL:    "/device/DEADBEEF/interface/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceInterfaces("DEADBEEF") },
		},
		{
			URL:    "/device/DEADBEEF/interface/eth0/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceInterfaceByName("DEADBEEF", "eth0") },
		},
		{
			URL:    "/device/DEADBEEF/interface/eth0/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetDeviceInterfaceField("DEADBEEF", "eth0", "foo") },
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
			test.Do(conch.New(ts.URL, "token", &logger{}))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}
