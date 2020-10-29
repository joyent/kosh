package conch

import (
	"github.com/qri-io/jsonschema"
)

// GetSchema (GET /json_schema) retieves the json-schema defined with the given
// path (/common/:name, /request/:name, /response/:name)
func (c *Client) GetSchema(path string) (schema jsonschema.Schema) {
	c.Schema(path).Receive(&schema)
	return
}
