package conch

import (
	"github.com/qri-io/jsonschema"
)

func (c *Client) GetSchema(name string) (schema jsonschema.Schema) {
	c.Schema(name).Receive(&schema)
	return
}
