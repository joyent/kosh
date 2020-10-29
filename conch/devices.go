package conch

import (
	"fmt"

	"github.com/joyent/kosh/conch/types"
)

// FindDevicesBySetting (GET /device?:key=:value) returns a list of devices
// that have the matching setting
func (c *Client) FindDevicesBySetting(key, value string) (device types.Devices) {
	c.Device("").WithParams(map[string]string{key: value}).Receive(&device)
	return
}

// FindDevicesByTag (GET /device?:key=:value) returns a list of devices that
// have the matching tag
func (c *Client) FindDevicesByTag(key, value string) (device types.Devices) {
	key = fmt.Sprintf("tag_%s", key)
	c.Device("").WithParams(map[string]string{key: value}).Receive(&device)
	return
}

// FindDevicesByField (GET /device?:key=:value) returns a list of devices that
// have the matching field
func (c *Client) FindDevicesByField(key, value string) (device types.Device) {
	c.Device("").WithParams(map[string]string{key: value}).Receive(&device)
	return
}

// GetDeviceBySerial (GET /device/:device_id_or_serial_number) retrieves a
// specific device by the given serial number string
func (c *Client) GetDeviceBySerial(serial string) (device types.DetailedDevice) {
	c.Device(serial).Receive(&device)
	return
}

// GetDeviceByID (GET /device/:device_id_or_serial_number) retrieves a specific
// device by the given UUID
func (c *Client) GetDeviceByID(id types.UUID) (device types.DetailedDevice) {
	c.Device(id.String()).Receive(&device)
	return
}

// GetDevicePXE (GET /device/:device_id_or_serial_number/pxe) returns the PXE
// information for the given device
func (c *Client) GetDevicePXE(id string) (pxe types.DevicePXE) {
	c.Device(id).PXE().Receive(&pxe)
	return
}

// GetDevicePhase (GET /device/:device_id_or_serial_number/phase) returns teh
// Phase for the given device
func (c *Client) GetDevicePhase(id string) (phase types.DevicePhase) {
	c.Device(id).Phase().Receive(&phase)
	return
}

// GetDeviceSKU (GET /device/:device_id_or_serial_number/sku) returns the SKU
// for the given device
func (c *Client) GetDeviceSKU(id string) (sku types.DeviceSku) {
	c.Device(id).SKU().Receive(&sku)
	return
}

// SetDeviceAssetTag (POST /device/:device_id_or_serial_number/asset_tag) sets
// a new asset tag for the given device
func (c *Client) SetDeviceAssetTag(id string, tag types.DeviceAssetTag) error {
	_, e := c.Device(id).AssetTag().Post(tag).Send()
	return e
}

// SetDeviceValidated (POST /device/:device_id_or_serial_number/validated) sets
// the validated field on the given device. THIS IS DEPRECATED.
func (c *Client) SetDeviceValidated(id string) error {
	_, e := c.Device(id).Validated().Post("").Send()
	return e
}

// SetDevicePhase (POST /device/:device_id_or_serial_number/phase) this sets
// the phase on the given device
func (c *Client) SetDevicePhase(id, phase string) error {
	_, e := c.Device(id).Phase().Post(types.DevicePhase(phase)).Send()
	return e
}

// SetDeviceLinks (POST /device/:device_id_or_serial_number/links) sets the
// links on the given device
func (c *Client) SetDeviceLinks(id string, links types.DeviceLinks) error {
	_, e := c.Device(id).Links().Post(links).Send()
	return e
}

// DeleteDeviceLinks (DELETE /device/:device_id_or_serial_number/links) removes
// the links for the given device
func (c *Client) DeleteDeviceLinks(id string) error {
	_, e := c.Device(id).Links().Delete().Send()
	return e
}

// SetDeviceSKU (POST /device/:device_id_or_serial_number/sku)  updates teh SKU
// for the given device
func (c *Client) SetDeviceSKU(id string, hardware types.DeviceHardware) error {
	_, e := c.Device(id).SKU().Post(hardware).Send()
	return e
}

// SetDeviceBuild (POST /device/:device_id_or_serial_number/build) updates teh
// build for the given device
func (c *Client) SetDeviceBuild(id string, build types.DeviceBuild) error {
	_, e := c.Device(id).Build("").Post(build).Send()
	return e
}

// GetDeviceLocation (GET /device/:device_id_or_serial_number/location) returns
// the known location data for the given device, this data is only accurate
// whhile the device is in preflight
func (c *Client) GetDeviceLocation(id string) (location types.DeviceLocation) {
	c.Device(id).Location().Receive(&location)
	return
}

// SetDeviceLocation (POST /device/:device_id_or_serial_number/location)
// udpates the current device location information for the given device
func (c *Client) SetDeviceLocation(id string, location types.DeviceLocationUpdate) error {
	_, e := c.Device(id).Location().Post(location).Send()
	return e
}

// DeleteDeviceLocation (DELETE /device/:device_id_or_serial_number/location)
// removes the given device location information for the given device
func (c *Client) DeleteDeviceLocation(id string) error {
	_, e := c.Device(id).Location().Delete().Send()
	return e
}

// GetDeviceSettings (GET /device/:device_id_or_serial_number/settings)
// retrieves the current settings for teh given device
func (c *Client) GetDeviceSettings(id string) (settings types.DeviceSettings) {
	c.Device(id).Settings("").Receive(&settings)
	return
}

// SetDeviceSettings (POST /device/:device_id_or_serial_number/settings)
// updates the current settings for the given device
func (c *Client) SetDeviceSettings(id string, settings types.DeviceSettings) error {
	_, e := c.Device(id).Settings("").Post(settings).Send()
	return e
}

// GetDeviceSettingByName (GET /device/:device_id_or_serial_number/settings/:key)
// retrieves a single device setting by name
func (c *Client) GetDeviceSettingByName(id, name string) (setting types.DeviceSetting) {
	c.Device(id).Settings(name).Receive(&setting)
	return
}

// GetDeviceTags (GET /device/:device_id_or_serial_number/settings/:key)
// retrieves the tags for a given device
func (c *Client) GetDeviceTags(id string) (tags types.DeviceSettings) {
	c.Device(id).Settings("").Receive(&tags)
	return
}

// GetDeviceTagByName (GET /device/:device_id_or_serial_number/settings/:key)
// retrieves a current tag by name for a given device
func (c *Client) GetDeviceTagByName(id, name string) (tag types.DeviceSetting) {
	name = fmt.Sprintf("tag_%s", name)
	c.Device(id).Settings(name).Receive(&tag)
	return
}

// SetDeviceTag (POST /device/:device_id_or_serial_number/settings/:key)
// updates a current tag by name for the given device
func (c *Client) SetDeviceTag(id, name, value string) error {
	name = fmt.Sprintf("tag_%s", name)
	_, e := c.Device(id).Settings(name).Post(value).Send()
	return e
}

// DeleteDeviceTag (POST /device/:device_id_or_serial_number/settings/:key)
// removes a current tag by name for the given device
func (c *Client) DeleteDeviceTag(id, name string) error {
	name = fmt.Sprintf("tag_%s", name)
	_, e := c.Device(id).Settings(name).Delete().Send()
	return e
}

// SetDeviceSetting (POST /device/:device_id_or_serial_number/settings/:key)
// updates a device setting by name for the given device
func (c *Client) SetDeviceSetting(id, name, value string) error {
	_, e := c.Device(id).Settings(name).Post(value).Send()
	return e
}

// DeleteDeviceSetting (DELETE /device/:device_id_or_serial_number/settings/:key)
// removes a device setting by name for the given device
func (c *Client) DeleteDeviceSetting(id, name string) error {
	_, e := c.Device(id).Settings(name).Delete().Send()
	return e
}

// RunValidationForDevice (POST /device/:device_id_or_serial_number/validation/:validation_id)
// runs a named validation for the given device and returns (but does not save)
// the results
func (c *Client) RunValidationForDevice(device, validation string, report types.DeviceReport) (results types.ValidationResults) {
	c.Device(device).Validation(validation).Post(report).Receive(&results)
	return
}

// GetDeviceValidationStates ( GET /device/:device_id_or_serial_number/validation_state?status=<pass|fail|error>&status=...)
// retrieves the given validation results for teh given device
func (c *Client) GetDeviceValidationStates(id string, states ...string) (validations types.ValidationStateWithResults) {
	c.Device(id).ValidationStates(states...).Receive(&validations)
	return
}

// GetDeviceInterfaces (GET /device/:device_id_or_serial_number/interface)
// returns the reported interface information for the given device
func (c *Client) GetDeviceInterfaces(id string) (nics types.DeviceNics) {
	c.Device(id).Interface("").Receive(&nics)
	return
}

// GetDeviceInterfaceByName (GET /device/:device_id_or_serial_number/interface)
// returns the information for a specific interface for the given device
func (c *Client) GetDeviceInterfaceByName(id, name string) (nic types.DeviceNic) {
	c.Device(id).Interface(name).Receive(&nic)
	return
}

// GetDeviceInterfaceField (GET /device/:device_id_or_serial_number/interface)
// retrieves a specific property of the named interface for a given device
func (c *Client) GetDeviceInterfaceField(id, name, field string) (nicField types.DeviceNicField) {
	c.Device(id).Interface(name).Field(field).Receive(&nicField)
	return
}
