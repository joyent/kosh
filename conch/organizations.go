package conch

import "github.com/joyent/kosh/conch/types"

// GetAllOrganizations (GET /organization) returns the list of organizations
func (c *Client) GetAllOrganizations() (orgs types.Organizations, e error) {
	_, e = c.Organization().Receive(&orgs)
	return
}

// CreateOrganization (POST /organization) creates a new organization
func (c *Client) CreateOrganization(org types.OrganizationCreate) error {
	_, e := c.Organization().Post(org).Send()
	return e
}

// GetOrganizationByName (GET /organization/:organization_id_or_name) retrieves
// an organziation by the given name
func (c *Client) GetOrganizationByName(name string) (org types.Organization, e error) {
	_, e = c.Organization(name).Receive(&org)
	return
}

// GetOrganizationByID (GET /organization/:organization_id_or_name) retrieves
// an organization by the given UUID
func (c *Client) GetOrganizationByID(id types.UUID) (org types.Organization, e error) {
	_, e = c.Organization(id.String()).Receive(&org)
	return
}

// UpdateOrganization (POST /organization/:organization_id_or_name) updates the
// organzation with the given name
func (c *Client) UpdateOrganization(name string, update types.OrganizationUpdate) error {
	_, e := c.Organization(name).Post(update).Send()
	return e
}

// DeleteOrganization (DELETE /organization/:organization_id_or_name)
// removes the organzation with the given UUID
func (c *Client) DeleteOrganization(id types.UUID) error {
	_, e := c.Organization(id.String()).Delete().Send()
	return e
}

// AddOrganizationUser (POST /organization/:organization_id_or_name/user?send_mail=<1|0>)
// adds a user to the orgaanization with the given UUID. Optionally sends an email to that user.
// BUG(perigrin): sendEmail flag is currently not implemented
func (c *Client) AddOrganizationUser(id types.UUID, user types.OrganizationAddUser, sendEmail bool) error {
	_, e := c.Organization(id.String()).User().Post(user).Send()
	return e
}

// DeleteOrganizationUser (DELETE /organization/:organization_id_or_name/user/#target_user_id_or_email?send_mail=<1|0>)
// removes a user from the organization with the given ID, optionally sends them an email
// BUG(perigrin): sendEmail flag is currently not implemented
func (c *Client) DeleteOrganizationUser(id types.UUID, user string, sendEmail bool) error {
	_, e := c.Organization(id.String()).User(user).Delete().Send()
	return e
}
