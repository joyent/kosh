package conch

import (
	"fmt"

	"github.com/joyent/kosh/conch/types"
)

// GetAllBuilds (GET /builds) retrieves a list of all Builds. The optional
// parameters may be "started" or "completed" which when set to 1 or 0 filters
// those builds that have been started or completed in or out of the results
func (c *Client) GetAllBuilds(options ...map[string]string) (builds types.Builds, e error) {
	params := make(map[string]string)
	for _, o := range options {
		for k, v := range o {
			params[k] = v
		}
	}
	if len(params) > 0 {
		_, e = c.Build("").WithParams(params).Receive(&builds)
	} else {
		_, e = c.Build("").Receive(&builds)
	}
	return
}

// CreateBuild (POST /build) creates a new Build
func (c *Client) CreateBuild(create types.BuildCreate) error {
	c.Info(fmt.Sprintf("creating build: %v", create))
	_, e := c.Build("").Post(create).Send()
	return e
}

// GetBuildByName (GET /build/:build_id_or_name) get's a single build by name
func (c *Client) GetBuildByName(name string) (build types.Build, e error) {
	c.Info(fmt.Sprintf("getting build by name: %s", name))
	_, e = c.Build(name).Receive(&build)
	return
}

// GetBuildByID (GET /build/:build_id_or_name) get's a single build by UUID
func (c *Client) GetBuildByID(id types.UUID) (build types.Build, e error) {
	c.Info(fmt.Sprintf("getting build by id: %s", id))
	_, e = c.Build(id.String()).Receive(&build)
	return
}

// UpdateBuild  updates a named build
// POST /build/:build_id_or_name
func (c *Client) UpdateBuild(name string, update types.BuildUpdate) error {
	return c.UpdateBuildByName(name, update)
}

// UpdateBuildByID (POST /build/:build_id_or_name)  updates a build
func (c *Client) UpdateBuildByID(buildID types.UUID, update types.BuildUpdate) error {
	c.Info(fmt.Sprintf("updating build %v: %v", buildID, update))
	_, e := c.Build(buildID.String()).Post(update).Send()
	return e
}

// UpdateBuild  updates a named build
// POST /build/:build_id_or_name
func (c *Client) UpdateBuildByName(name string, update types.BuildUpdate) error {
	c.Info(fmt.Sprintf("updating build %v: %v", name, update))
	_, e := c.Build(name).Post(update).Send()
	return e
}

// GetBuildUsers retrieves a list of users associated with the given build
// GET /build/:build_id_or_name/user
func (c *Client) GetBuildUsers(name string) (build types.BuildUsers, e error) {
	c.Info(fmt.Sprintf("getting users for build: %s", name))
	_, e = c.Build(name).User("").Receive(&build)
	return
}

// AddBuildUser associates a new user with the build, optionally tell the API to email the user too
// POST /build/:build_id_or_name/user
func (c *Client) AddBuildUser(name string, update types.BuildAddUser, sendEmail bool) error {
	c.Info(fmt.Sprintf("adding users to build %v: %v", name, update))
	_, e := c.Build(name).User("").Post(update).Send()
	return e
}

// DeleteBuildUser removes a user from being associated with the build
// DELETE /build/:build_id_or_name/user/#target_user_id_or_email
func (c *Client) DeleteBuildUser(name, user string, sendEmail bool) error {
	c.Info(fmt.Sprintf("removing user from build %v: %v", name, user))
	_, e := c.Build(name).User(user).Delete().Send()
	return e
}

// GetAllBuildOrganizations - GET /build/:build_id_or_name/user
func (c *Client) GetAllBuildOrganizations(name string) (build types.BuildOrganizations, e error) {
	c.Info(fmt.Sprintf("getting organizations for build: %s", name))
	_, e = c.Build(name).Organization("").Receive(&build)
	return
}

// AddBuildOrganization - POST /build/:build_id_or_name/user
func (c *Client) AddBuildOrganization(name string, update types.BuildAddOrganization, sendEmail bool) error {
	c.Info(fmt.Sprintf("adding organization to build %v: %v", name, update))
	_, e := c.Build(name).Organization("").Post(update).Send()
	return e
}

// DeleteBuildOrganization - DELETE /build/:build_id_or_name/user/#target_user_id_or_email
func (c *Client) DeleteBuildOrganization(build, org string, sendEmail bool) error {
	c.Info(fmt.Sprintf("removing organization from build %v: %v", build, org))
	_, e := c.Build(build).Organization(org).Delete().Send()
	return e
}

// GetAllBuildDevices - GET /build/:build_id_or_name/device
func (c *Client) GetAllBuildDevices(name string) (list types.Devices, e error) {
	c.Info(fmt.Sprintf("getting devices for build: %s", name))
	_, e = c.Build(name).Device("").Receive(&list)
	return
}

// GetBuildDevicesPXE - GET /build/:build_id_or_name/device/pxe
func (c *Client) GetBuildDevicesPXE(name string) (list types.DevicePxes, e error) {
	c.Info(fmt.Sprintf("getting device PXE info for build: %s", name))
	_, e = c.Build(name).Device("").PXE().Receive(&list)
	return
}

// AddNewBuildDevice - POST /build/:build_id_or_name/device
func (c *Client) AddNewBuildDevice(name string, device types.BuildCreateDevices) error {
	c.Info(fmt.Sprintf("adding device to build %v: %v", name, device))
	_, e := c.Build(name).Device("").Post(device).Send()
	return e
}

// AddBuildDeviceByName - POST /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) AddBuildDeviceByName(name, device string) error {
	c.Info(fmt.Sprintf("adding device to build %v: %v", name, device))
	_, e := c.Build(name).Device(device).Post("").Send()
	return e
}

// AddBuildDeviceByID - POST /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) AddBuildDeviceByID(buildID, deviceID types.UUID) error {
	c.Info(fmt.Sprintf("adding device to build %v: %v", buildID, deviceID))
	_, e := c.Build(buildID.String()).Device(deviceID.String()).Post("").Send()
	return e
}

// DeleteBuildDeviceByID - DELETE /build/:build_id_or_name/device/:device_id_or_serial_number
func (c *Client) DeleteBuildDeviceByID(buildID, deviceID types.UUID) error {
	c.Info(fmt.Sprintf("removing device from build %v: %v", buildID, deviceID))
	_, e := c.Build(buildID.String()).Device(deviceID.String()).Delete().Send()
	return e
}

// GetBuildRacks - GET /build/:build_id_or_name/rack
func (c *Client) GetBuildRacks(name string) (racks types.Racks, e error) {
	c.Info(fmt.Sprintf("getting racks for build: %s", name))
	_, e = c.Build(name).Rack("").Receive(&racks)
	return
}

// AddBuildRackByID - POST /build/:build_id_or_name/rack/:rack_id_or_name
func (c *Client) AddBuildRackByID(name, rack string) error {
	c.Info(fmt.Sprintf("adding rack to build %v: %v", name, rack))
	_, e := c.Build(name).Rack(rack).Post("").Send()
	return e
}

// DeleteBuildRackByID - DELETE /build/:build_id_or_name/rack/:rack_id_or_name
func (c *Client) DeleteBuildRackByID(name, rack string) error {
	c.Info(fmt.Sprintf("removing rack from build %v: %v", name, rack))
	_, e := c.Build(name).Rack(rack).Delete().Send()
	return e
}
