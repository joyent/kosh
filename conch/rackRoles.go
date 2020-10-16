package conch

import "github.com/joyent/kosh/conch/types"

// GET /rack_role
func (c *Client) GetAllRackRoles() (roles types.RackRoles) {
	c.RackRole().Receive(roles)
	return
}

// POST /rack_role
func (c *Client) CreateRackRole(role types.RackRoleCreate) error {
	_, e := c.RackRole().Post(role).Send()
	return e
}

// GET /rack_role/:rack_role_id_or_name
func (c *Client) GetRackRoleByName(name string) (role types.RackRole) {
	c.RackRole(name).Receive(role)
	return
}

// GET /rack_role/:rack_role_id_or_name
func (c *Client) GetRackRoleByID(id types.UUID) (role types.RackRole) {
	c.RackRole(id.String()).Receive(role)
	return
}

// POST /rack_role/:rack_role_id_or_name
func (c *Client) UpdateRackRole(id types.UUID, update types.RackRoleUpdate) error {
	_, e := c.RackRole(id.String()).Post(update).Send()
	return e
}

// DELETE /rack_role/:rack_role_id_or_name
func (c *Client) DeleteRackRole(id types.UUID) error {
	_, e := c.RackRole(id.String()).Delete().Send()
	return e
}
