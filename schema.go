package main

import (
	"fmt"

	"github.com/qri-io/jsonschema"
)

type Schema struct {
	*Conch
}

func (c *Conch) Schema() *Schema {
	return &Schema{c}
}

func (s *Schema) Get(name string) *jsonschema.RootSchema {
	uri := fmt.Sprintf("/schema/%s", name)
	rs := &jsonschema.RootSchema{}

	// for now we completely skip our internal HTTP handling and just use sling directly
	_, err := s.Sling().Get(uri).Receive(&rs, nil)
	if err != nil {
		panic(err)
	}
	return rs
}
