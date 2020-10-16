package conch

import "github.com/joyent/kosh/conch/types"

// GET /hardware_vendor
func (c *Client) GetHardwareVendors() (vendors types.HardwareVendors) {
	c.HardwareVendor().Receive(&vendors)
	return
}

// GET /hardware_vendor/:hardware_vendor_id_or_name
func (c *Client) GetHardwareVendorByID(id string) (vendor types.HardwareVendor) {
	c.HardwareVendor(id).Receive(&vendor)
	return
}

// DELETE /hardware_vendor/:hardware_vendor_id_or_name
func (c *Client) DeleteHardwareVendor(id types.UUID) error {
	_, e := c.HardwareVendor(id.String()).Delete().Send()
	return e
}

func (c *Client) FindOrCreateHardwareVendor(id string) (vendor types.HardwareVendor) {
	vendor = c.GetHardwareVendorByID(id)
	if (vendor == types.HardwareVendor{}) {
		c.CreateHardwareVendor(id)
		vendor = c.GetHardwareVendorByID(id)
	}
	return
}

// POST /hardware_vendor/:hardware_vendor_id_or_name
func (c *Client) CreateHardwareVendor(id string) error {
	_, e := c.HardwareVendor(id).Post("").Send()
	return e
}
