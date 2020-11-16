package conch

import "github.com/joyent/kosh/conch/types"

// GetAllDatacenters ( GET /dc ) retrieves a list of all datacenters
func (c *Client) GetAllDatacenters() (dc types.Datacenters, e error) {
	c.DC("").Receive(&dc)
	return
}

// CreateDatacenter (POST /dc) posts a new Datacenter to teh API
func (c *Client) CreateDatacenter(dc types.DatacenterCreate) error {
	_, e := c.DC("").Post(dc).Send()
	return e
}

// GetDatacenterByName ( GET /dc/:datacenter_id ) fetches a new datacenter
// using the given string
func (c *Client) GetDatacenterByName(name string) (dc types.Datacenter, e error) {
	c.DC(name).Receive(&dc)
	return
}

// GetDatacenterByID ( GET /dc/:datacenter_id ) fetches a new datacenter using
// the given UUID
func (c *Client) GetDatacenterByID(id types.UUID) (dc types.Datacenter, e error) {
	c.DC(id.String()).Receive(&dc)
	return
}

// UpdateDatacenter (POST /dc/:datacenter_id) updates datacenter with the given
// UUID
func (c *Client) UpdateDatacenter(id types.UUID, update types.DatacenterUpdate) error {
	_, e := c.DC(id.String()).Post(update).Send()
	return e
}

// DeleteDatacenter (DELETE /dc/:datacenter_id) removes the given datacenter
// from the API
func (c *Client) DeleteDatacenter(id types.UUID) error {
	_, e := c.DC(id.String()).Delete().Send()
	return e
}

// GetAllDatacenterRooms ( GET /dc/:datacenter_id/rooms ) retrieves a list fo
// rooms in the given datacenter
func (c *Client) GetAllDatacenterRooms(id types.UUID) (rooms types.DatacenterRoomsDetailed, e error) {
	c.DC(id.String()).Rooms().Receive(&rooms)
	return
}
