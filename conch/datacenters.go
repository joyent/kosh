package conch

import "github.com/joyent/kosh/conch/types"

// GET /dc
func (c *Client) GetAllDatacenters() (dc types.Datacenters) {
	c.DC("").Receive(&dc)
	return
}

// POST /dc
func (c *Client) CreateDatacenter(dc types.DatacenterCreate) error {
	_, e := c.DC("").Post(dc).Send()
	return e
}

// GET /dc/:datacenter_id
func (c *Client) GetDatacenterByName(name string) (dc types.Datacenter) {
	c.DC(name).Receive(&dc)
	return
}

func (c *Client) GetDatacenterByID(id types.UUID) (dc types.Datacenter) {
	c.DC(id.String()).Receive(&dc)
	return
}

// POST /dc/:datacenter_id
func (c *Client) UpdateDatacenter(id types.UUID, update types.DatacenterUpdate) error {
	_, e := c.DC(id.String()).Post(update).Send()
	return e
}

// DELETE /dc/:datacenter_id
func (c *Client) DeleteDatacenter(id types.UUID) error {
	_, e := c.DC(id.String()).Delete().Send()
	return e
}

// GET /dc/:datacenter_id/rooms
func (c *Client) GetAllDatacenterRooms(id types.UUID) (rooms types.DatacenterRoomsDetailed) {
	c.DC(id.String()).Rooms().Receive(&rooms)
	return
}
