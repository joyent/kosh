package conch

import "github.com/joyent/kosh/conch/types"

// GET /builds
func (c *Client) GetBuilds() (builds types.Builds) {
	c.Build("").Receive(builds)
	return
}

// POST /build
func (c *Client) CreateBuild(create types.BuildCreate) error {
	_, e := c.Build("").Post(create).Send()
	return e
}

// GET /build/:build_id_or_name
func (c *Client) GetBuildByName(name string) (build types.Build) {
	c.Build(name).Receive(build)
	return
}

// POST /build/:build_id_or_name
func (c *Client) UpdateBuild(name string, update types.BuildUpdate) error {
	_, e := c.Build(name).Post(update).Send()
	return e
}

// GET /build/:build_id_or_name/user
func (c *Client) GetBuildUsers(name string) (build types.BuildUsers) {
	c.Build(name).User("").Receive(build)
	return
}

// POST /build/:build_id_or_name/user
func (c *Client) AddBuildUser(name string, update types.BuildAddUser) error {
	_, e := c.Build(name).User("").Post(update).Send()
	return e
}

// DELETE /build/:build_id_or_name/user/#target_user_id_or_email
func (c *Client) DeleteBuildUser(name, user string) error {
	_, e := c.Build(name).User(user).Delete().Send()
	return e
}

// GET /build/:build_id_or_name/user
func (c *Client) GetBuildOrganizations(name string) (build types.BuildOrganizations) {
	c.Build(name).Organization("").Receive(build)
	return
}

// POST /build/:build_id_or_name/user
func (c *Client) AddBuildOrganization(name string, update types.BuildAddOrganization) error {
	_, e := c.Build(name).Organization("").Post(update).Send()
	return e
}

// DELETE /build/:build_id_or_name/user/#target_user_id_or_email
func (c *Client) DeleteBuildOrganization(build, org string) error {
	_, e := c.Build(build).Organization(org).Delete().Send()
	return e
}

// GET /build/:build_id_or_name/device
func (c *Client) GetBuildDevices(name string) (build types.Build) {
	c.Build(name).Device("").Receive(build)
	return
}

// GET /build/:build_id_or_name/device/pxe
func (c *Client) GetBuildDevicesPXE(name string) (build types.Build) {
	c.Build(name).Device("").PXE().Receive(build)
	return
}

// POST /build/:build_id_or_name/device
func (c *Client) AddNewBuildDevice(name string, device types.BuildCreateDevices) error {
	_, e := c.Build(name).Device("").Post(device).Send()
	return e
}

// POST /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) AddBuildDeviceByID(name, device string) error {
	_, e := c.Build(name).Device(device).Post("").Send()
	return e
}

// DELETE /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) DeleteBuildDeviceByID(name, device string) error {
	_, e := c.Build(name).Device(device).Delete().Send()
	return e
}

// GET /build/:build_id_or_name/rack
func (c *Client) GetBuildRacks(name string) (racks types.Racks) {
	c.Build(name).Rack("").Receive(racks)
	return
}

// GET /build/:build_id_or_name/rack/:rack_id_or_name
func (c *Client) AddBuildRackByID(name, rack string) error {
	_, e := c.Build(name).Rack(rack).Post("").Send()
	return e
}
