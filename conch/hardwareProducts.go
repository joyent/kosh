package conch

import "github.com/joyent/kosh/conch/types"

// GET /hardware_product
func (c *Client) GetHardwareProducts() (products types.HardwareProducts) {
	c.HardwareProduct().Receive(products)
	return
}

// POST /hardware_product
func (c *Client) CreateHardwareProduct(product types.HardwareProductCreate) error {
	_, e := c.HardwareProduct().Post(product).Send()
	return e
}

// GET /hardware_product/:hardware_product_id_or_other
func (c *Client) GetHardwareProductByID(id string) (products types.HardwareProduct) {
	c.HardwareProduct(id).Receive(products)
	return
}

// POST /hardware_product/:hardware_product_id_or_other
func (c *Client) UpdateHardwareProduct(id string, update types.HardwareProductUpdate) error {
	_, e := c.HardwareProduct(id).Post(update).Send()
	return e
}

// DELETE /hardware_product/:hardware_product_id_or_other
func (c *Client) DeleteHardwareProduct(id string) error {
	_, e := c.HardwareProduct(id).Delete().Send()
	return e
}

// PUT /hardware_product/:hardware_product_id_or_other/specification?path=:path_to_data
func (c *Client) UpdateHardwareProductSpecification(id, path string, update types.HardwareProductSpecification) error {
	_, e := c.HardwareProduct(id).Specification(path).Put(update).Send()
	return e
}

// DELETE /hardware_product/:hardware_product_id_or_other/specification?path=:path_to_data
func (c *Client) DeleteHardwareProductSpecification(id, path string) error {
	_, e := c.HardwareProduct(id).Specification(path).Delete().Send()
	return e
}
