package conch

import (
	"encoding/json"
	"io"

	"github.com/joyent/kosh/conch/types"
)

// GetHardwareProducts (GET /hardware_product) returns a list of known hardware
// products
func (c *Client) GetHardwareProducts() (products types.HardwareProducts, e error) {
	_, e = c.HardwareProduct().Receive(&products)
	return
}

// CreateHardwareProduct (POST /hardware_product) creates a new hardware
// product
func (c *Client) CreateHardwareProduct(product types.HardwareProductCreate) error {
	_, e := c.HardwareProduct().Post(product).Send()
	return e
}

// ReadHardwareProduct takes an io reader and returns a HardwareProductCreate
// struct suitable for CreateHardwareProduct
func (c *Client) ReadHardwareProduct(r io.Reader) (create types.HardwareProductCreate) {
	json.NewDecoder(r).Decode(&create)
	return
}

// GetHardwareProductByID (GET /hardware_product/:hardware_product_id_or_other)
// returns a hardware product by the given id string
func (c *Client) GetHardwareProductByID(id string) (products types.HardwareProduct, e error) {
	_, e = c.HardwareProduct(id).Receive(&products)
	return
}

// UpdateHardwareProduct (POST /hardware_product/:hardware_product_id_or_other)
// updates the given hardware product information
func (c *Client) UpdateHardwareProduct(id string, update types.HardwareProductUpdate) error {
	_, e := c.HardwareProduct(id).Post(update).Send()
	return e
}

// DeleteHardwareProduct (DELETE /hardware_product/:hardware_product_id_or_other)
// removes a hardware product
func (c *Client) DeleteHardwareProduct(id types.UUID) error {
	_, e := c.HardwareProduct(id.String()).Delete().Send()
	return e
}

// UpdateHardwareProductSpecification (PUT /hardware_product/:hardware_product_id_or_other/specification?path=:path_to_data)
// updates the Hardware Product Specification information at the given path for
// the given hardware product
func (c *Client) UpdateHardwareProductSpecification(id, path string, update types.HardwareProductSpecification) error {
	_, e := c.HardwareProduct(id).Specification(path).Put(update).Send()
	return e
}

// DeleteHardwareProductSpecification (DELETE /hardware_product/:hardware_product_id_or_other/specification?path=:path_to_data)
// removes the specification at teh given path for the given hardware product
func (c *Client) DeleteHardwareProductSpecification(id, path string) error {
	_, e := c.HardwareProduct(id).Specification(path).Delete().Send()
	return e
}
