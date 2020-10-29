package conch

import "github.com/joyent/kosh/conch/types"

// GetAllBuilds retrieves a list of all Builds
// GET /builds
func (c *Client) GetAllBuilds(options ...types.BuildQueryOptions) (builds types.Builds) {
	c.Build("").Receive(&builds)
	return
}

// CreateBuild creates a new Build
// POST /build
func (c *Client) CreateBuild(create types.BuildCreate) error {
	_, e := c.Build("").Post(create).Send()
	return e
}

// GetBuildByName get's a single build by "name" ie a string that may be a Name or UUID
// GET /build/:build_id_or_name
func (c *Client) GetBuildByName(name string) (build types.Build) {
	c.Build(name).Receive(&build)
	return
}

// UpdateBuild  updates a named build
// POST /build/:build_id_or_name
func (c *Client) UpdateBuild(name string, update types.BuildUpdate) error {
	_, e := c.Build(name).Post(update).Send()
	return e
}

// GetBuildUsers retrieves a list of users associated with the given build
// GET /build/:build_id_or_name/user
func (c *Client) GetBuildUsers(name string) (build types.BuildUsers) {
	c.Build(name).User("").Receive(&build)
	return
}

// AddBuildUser associates a new user with the build, optionally tell the API to email the user too
// POST /build/:build_id_or_name/user
func (c *Client) AddBuildUser(name string, update types.BuildAddUser, sendEmail bool) error {
	_, e := c.Build(name).User("").Post(update).Send()
	return e
}

// DeleteBuildUser removes a user from being associated with the build
// DELETE /build/:build_id_or_name/user/#target_user_id_or_email
func (c *Client) DeleteBuildUser(name, user string, sendEmail bool) error {
	_, e := c.Build(name).User(user).Delete().Send()
	return e
}

// GetAllBuildOrganizations - GET /build/:build_id_or_name/user
func (c *Client) GetAllBuildOrganizations(name string) (build types.BuildOrganizations) {
	c.Build(name).Organization("").Receive(&build)
	return
}

// AddBuildOrganization - POST /build/:build_id_or_name/user
func (c *Client) AddBuildOrganization(name string, update types.BuildAddOrganization, sendEmail bool) error {
	_, e := c.Build(name).Organization("").Post(update).Send()
	return e
}

// DeleteBuildOrganization - DELETE /build/:build_id_or_name/user/#target_user_id_or_email
func (c *Client) DeleteBuildOrganization(build, org string, sendEmail bool) error {
	_, e := c.Build(build).Organization(org).Delete().Send()
	return e
}

// GetAllBuildDevices - GET /build/:build_id_or_name/device
func (c *Client) GetAllBuildDevices(name string) (build types.Build) {
	c.Build(name).Device("").Receive(&build)
	return
}

// GetBuildDevicesPXE - GET /build/:build_id_or_name/device/pxe
func (c *Client) GetBuildDevicesPXE(name string) (build types.Build) {
	c.Build(name).Device("").PXE().Receive(&build)
	return
}

// AddNewBuildDevice - POST /build/:build_id_or_name/device
func (c *Client) AddNewBuildDevice(name string, device types.BuildCreateDevices) error {
	_, e := c.Build(name).Device("").Post(device).Send()
	return e
}

// AddBuildDeviceByName - POST /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) AddBuildDeviceByName(name, device string) error {
	_, e := c.Build(name).Device(device).Post("").Send()
	return e
}

// AddBuildDeviceByID - POST /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) AddBuildDeviceByID(buildID, deviceID types.UUID) error {
	_, e := c.Build(buildID.String()).Device(deviceID.String()).Post("").Send()
	return e
}

// DeleteBuildDeviceByID - DELETE /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) DeleteBuildDeviceByID(buildID, deviceID types.UUID) error {
	_, e := c.Build(buildID.String()).Device(deviceID.String()).Delete().Send()
	return e
}

// GetBuildRacks - GET /build/:build_id_or_name/rack
func (c *Client) GetBuildRacks(name string) (racks types.Racks) {
	c.Build(name).Rack("").Receive(&racks)
	return
}

// AddBuildRackByID - POST /build/:build_id_or_name/rack/:rack_id_or_name
func (c *Client) AddBuildRackByID(name, rack string) error {
	_, e := c.Build(name).Rack(rack).Post("").Send()
	return e
}

// DeleteBuildRackByID - DELETE /build/:build_id_or_name/rack/:rack_id_or_name
func (c *Client) DeleteBuildRackByID(name, rack string) error {
	_, e := c.Build(name).Rack(rack).Delete().Send()
	return e
}
