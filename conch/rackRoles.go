package conch

import "github.com/joyent/kosh/v3/conch/types"

// GetAllRackRoles (GET /rack_role) returns a list of all rack roles
func (c *Client) GetAllRackRoles() (roles types.RackRoles) {
	c.RackRole().Receive(roles)
	return
}

// CreateRackRole (POST /rack_role) creates a new rack role
func (c *Client) CreateRackRole(role types.RackRoleCreate) error {
	_, e := c.RackRole().Post(role).Send()
	return e
}

// GetRackRoleByName (GET /rack_role/:rack_role_id_or_name)
// retrieves the rack role for the given name
func (c *Client) GetRackRoleByName(name string) (role types.RackRole) {
	c.RackRole(name).Receive(role)
	return
}

// GetRackRoleByID (GET /rack_role/:rack_role_id_or_name) retrieves the rack
// role for the given UUID
func (c *Client) GetRackRoleByID(id types.UUID) (role types.RackRole) {
	c.RackRole(id.String()).Receive(role)
	return
}

// UpdateRackRole (POST /rack_role/:rack_role_id_or_name) updates the rack role
// with the given UUID
func (c *Client) UpdateRackRole(id types.UUID, update types.RackRoleUpdate) error {
	_, e := c.RackRole(id.String()).Post(update).Send()
	return e
}

// DeleteRackRole (DELETE /rack_role/:rack_role_id_or_name) removes the rack
// role with the given UUID
func (c *Client) DeleteRackRole(id types.UUID) error {
	_, e := c.RackRole(id.String()).Delete().Send()
	return e
}
