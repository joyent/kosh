package conch

import "github.com/joyent/kosh/conch/types"

// GET /dc
func (c *Client) GetDatacenters() (dc types.Datacenters) {
	c.DC("").Receive(dc)
	return
}

// POST /dc
func (c *Client) CreateDatacenter(dc types.DatacenterCreate) error {
	_, e := c.DC("").Post(dc).Send()
	return e
}

// GET /dc/:datacenter_id
func (c *Client) GetDatacenterByID(id string) (dc types.Datacenter) {
	c.DC(id).Receive(dc)
	return
}

// POST /dc/:datacenter_id
func (c *Client) UpdateDatacenter(id string, update types.DatacenterUpdate) error {
	_, e := c.DC(id).Post(update).Send()
	return e
}

// DELETE /dc/:datacenter_id
func (c *Client) DeleteDatacenter(id string) error {
	_, e := c.DC(id).Delete().Send()
	return e
}

// GET /dc/:datacenter_id/rooms
func (c *Client) GetDatacenterRooms(id string) (rooms types.DatacenterRoomsDetailed) {
	c.DC(id).Rooms().Receive(rooms)
	return
}
