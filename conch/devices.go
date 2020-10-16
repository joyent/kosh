/*

API Wrapper for https://joyent.github.io/conch-api/modules/Conch::Route::Device

*/
package conch

import (
	"fmt"

	"github.com/joyent/kosh/conch/types"
)

// GET /device?:key=:value
func (c *Client) FindDevicesBySetting(key, value string) (device types.Devices) {
	c.Device("").WithParams(map[string]string{key: value}).Receive(&device)
	return
}

func (c *Client) FindDevicesByTag(key, value string) (device types.Devices) {
	key = fmt.Sprintf("tag_%s", key)
	c.Device("").WithParams(map[string]string{key: value}).Receive(&device)
	return
}

func (c *Client) FindDevicesByField(key, value string) (device types.Device) {
	c.Device("").WithParams(map[string]string{key: value}).Receive(&device)
	return
}

// GET /device/:device_id_or_serial_number
func (c *Client) GetDeviceBySerial(serial string) (device types.DetailedDevice) {
	c.Device(serial).Receive(&device)
	return
}

// GET /device/:device_id_or_serial_number
func (c *Client) GetDeviceByID(id types.UUID) (device types.DetailedDevice) {
	c.Device(id.String()).Receive(&device)
	return
}

// GET /device/:device_id_or_serial_number/pxe
func (c *Client) GetDevicePXE(id string) (pxe types.DevicePXE) {
	c.Device(id).PXE().Receive(&pxe)
	return
}

// GET /device/:device_id_or_serial_number/phase
func (c *Client) GetDevicePhase(id string) (phase types.DevicePhase) {
	c.Device(id).Phase().Receive(&phase)
	return
}

// GET /device/:device_id_or_serial_number/sku
func (c *Client) GetDeviceSKU(id string) (sku types.DeviceSku) {
	c.Device(id).SKU().Receive(&sku)
	return
}

// POST /device/:device_id_or_serial_number/asset_tag
func (c *Client) SetDeviceAssetTag(id string, tag types.DeviceAssetTag) error {
	_, e := c.Device(id).AssetTag().Post(tag).Send()
	return e
}

// POST /device/:device_id_or_serial_number/asset_tag
func (c *Client) SetDeviceValidated(id string) error {
	_, e := c.Device(id).Validated().Post("").Send()
	return e
}

// POST /device/:device_id_or_serial_number/phase
func (c *Client) SetDevicePhase(id, phase string) error {
	_, e := c.Device(id).Phase().Post(types.DevicePhase(phase)).Send()
	return e
}

// POST /device/:device_id_or_serial_number/links
func (c *Client) SetDeviceLinks(id string, links types.DeviceLinks) error {
	_, e := c.Device(id).Links().Post(links).Send()
	return e
}

// DELETE /device/:device_id_or_serial_number/links
func (c *Client) DeleteDeviceLinks(id string) error {
	_, e := c.Device(id).Links().Delete().Send()
	return e
}

// POST /device/:device_id_or_serial_number/hardware_product
// POST /device/:device_id_or_serial_number/sku
func (c *Client) SetDeviceSKU(id string, hardware types.DeviceHardware) error {
	_, e := c.Device(id).SKU().Post(hardware).Send()
	return e
}

// POST /device/:device_id_or_serial_number/links
func (c *Client) SetDeviceBuild(id string, build types.DeviceBuild) error {
	_, e := c.Device(id).Build("").Post(build).Send()
	return e
}

// GET /device/:device_id_or_serial_number/location
func (c *Client) GetDeviceLocation(id string) (location types.DeviceLocation) {
	c.Device(id).Location().Receive(&location)
	return
}

// POST /device/:device_id_or_serial_number/location
func (c *Client) SetDeviceLocation(id string, location types.DeviceLocationUpdate) error {
	_, e := c.Device(id).Location().Post(location).Send()
	return e
}

// DELETE /device/:device_id_or_serial_number/location
func (c *Client) DeleteDeviceLocation(id string) error {
	_, e := c.Device(id).Location().Delete().Send()
	return e
}

// GET /device/:device_id_or_serial_number/settings
func (c *Client) GetDeviceSettings(id string) (settings types.DeviceSettings) {
	c.Device(id).Settings("").Receive(&settings)
	return
}

// POST /device/:device_id_or_serial_number/settings
func (c *Client) SetDeviceSettings(id string, settings types.DeviceSettings) error {
	_, e := c.Device(id).Settings("").Post(settings).Send()
	return e
}

// GET /device/:device_id_or_serial_number/settings/:key
func (c *Client) GetDeviceSettingByName(id, name string) (setting types.DeviceSetting) {
	c.Device(id).Settings(name).Receive(&setting)
	return
}

func (c *Client) GetDeviceTags(id string) (tags types.DeviceSettings) {
	c.Device(id).Settings("").Receive(&tags)
	return
}

func (c *Client) GetDeviceTagByName(id, name string) (tag types.DeviceSetting) {
	name = fmt.Sprintf("tag_%s", name)
	c.Device(id).Settings(name).Receive(&tag)
	return
}

func (c *Client) SetDeviceTag(id, name, value string) error {
	name = fmt.Sprintf("tag_%s", name)
	_, e := c.Device(id).Settings(name).Post(value).Send()
	return e
}

func (c *Client) DeleteDeviceTag(id, name string) error {
	name = fmt.Sprintf("tag_%s", name)
	_, e := c.Device(id).Settings(name).Delete().Send()
	return e
}

// POST /device/:device_id_or_serial_number/settings/:key
func (c *Client) SetDeviceSetting(id, name, value string) error {
	_, e := c.Device(id).Settings(name).Post(value).Send()
	return e
}

// DELETE /device/:device_id_or_serial_number/settings/:key
func (c *Client) DeleteDeviceSetting(id, name string) error {
	_, e := c.Device(id).Settings(name).Delete().Send()
	return e
}

// POST /device/:device_id_or_serial_number/validation/:validation_id
func (c *Client) RunValidationForDevice(device, validation string, report types.DeviceReport) (results types.ValidationResults) {
	c.Device(device).Validation(validation).Post(report).Receive(&results)
	return
}

// GET /device/:device_id_or_serial_number/validation_state?status=<pass|fail|error>&status=...
func (c *Client) GetDeviceValidationStates(id string, states ...string) (validations types.ValidationStateWithResults) {
	c.Device(id).ValidationStates(states...).Receive(&validations)
	return
}

// GET /device/:device_id_or_serial_number/interface
func (c *Client) GetDeviceInterfaces(id string) (nics types.DeviceNics) {
	c.Device(id).Interface("").Receive(&nics)
	return
}

// GET /device/:device_id_or_serial_number/interface
func (c *Client) GetDeviceInterfaceByName(id, name string) (nic types.DeviceNic) {
	c.Device(id).Interface(name).Receive(&nic)
	return
}

// GET /device/:device_id_or_serial_number/interface
func (c *Client) GetDeviceInterfaceField(id, name, field string) (nicField types.DeviceNicField) {
	c.Device(id).Interface(name).Field(field).Receive(&nicField)
	return
}
