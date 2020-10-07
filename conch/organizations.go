package conch

import "github.com/joyent/kosh/conch/types"

// GET /organization
func (c *Client) GetOrganizations() (orgs types.Organizations) {
	c.Organization().Receive(orgs)
	return
}

// POST /organization
func (c *Client) CreateOrganization(org types.OrganizationCreate) error {
	_, e := c.Organization().Post(org).Send()
	return e
}

// GET /organization/:organization_id_or_name
func (c *Client) GetOrganizationByID(id string) (org types.Organization) {
	c.Organization(id).Receive(org)
	return
}

// POST /organization/:organization_id_or_name
func (c *Client) UpdateOrganization(id string, update types.OrganizationUpdate) error {
	_, e := c.Organization(id).Post(update).Send()
	return e
}

// DELETE /organization/:organization_id_or_name
func (c *Client) DeleteOrganization(id string) error {
	_, e := c.Organization(id).Delete().Send()
	return e
}

// POST /organization/:organization_id_or_name/user?send_mail=<1|0>
func (c *Client) AddOrganizationUser(id string, user types.OrganizationAddUser) error {
	_, e := c.Organization(id).Post(user).Send()
	return e
}

// DELETE /organization/:organization_id_or_name/user/#target_user_id_or_email?send_mail=<1|0>
func (c *Client) DeleteOrganizationUser(id, user string) error {
	_, e := c.Organization(id).User(user).Delete().Send()
	return e
}
