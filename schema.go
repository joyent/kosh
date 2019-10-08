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

	res := s.Do(s.Sling().Get(uri))
	if ok := res.Parse(&rs); !ok {
		panic(res)
	}
	return rs
}
