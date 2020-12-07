package conch

import "github.com/joyent/kosh/conch/types"

// GetAllHardwareVendors (GET /hardware_vendor) returns a list of all hardware
// vendors
func (c *Client) GetAllHardwareVendors() (vendors types.HardwareVendors, e error) {
	_, e = c.HardwareVendor().Receive(&vendors)
	return
}

// GetHardwareVendorByName (GET /hardware_vendor/:hardware_vendor_id_or_name)
// returns a specific hardware vendor by the given name
func (c *Client) GetHardwareVendorByName(name string) (vendor types.HardwareVendor, e error) {
	_, e = c.HardwareVendor(name).Receive(&vendor)
	return
}

// GetHardwareVendorByID (GET /hardware_vendor/:hardware_vendor_id_or_name)
// returns a specific hardware vendor by the given UUID
func (c *Client) GetHardwareVendorByID(id types.UUID) (vendor types.HardwareVendor, e error) {
	_, e = c.HardwareVendor(id.String()).Receive(&vendor)
	return
}

// DeleteHardwareVendor (DELETE /hardware_vendor/:hardware_vendor_id_or_name)
// removes the hardware vendor with the given UUID
func (c *Client) DeleteHardwareVendor(id types.UUID) error {
	_, e := c.HardwareVendor(id.String()).Delete().Send()
	return e
}

// FindOrCreateHardwareVendor optionally creates a new hardawre vendor with a
// given name if it does not already exist
func (c *Client) FindOrCreateHardwareVendor(name string) (vendor types.HardwareVendor, e error) {
	vendor, e = c.GetHardwareVendorByName(name)
	if (vendor == types.HardwareVendor{}) {
		c.CreateHardwareVendor(name)
		vendor, e = c.GetHardwareVendorByName(name)
	}
	return
}

// CreateHardwareVendor (POST /hardware_vendor/:hardware_vendor_id_or_name)
// createa a new hardware vendor with the given name
func (c *Client) CreateHardwareVendor(name string) error {
	_, e := c.HardwareVendor(name).Post().Send()
	return e
}
