package conch

import (
	"encoding/json"
	"io"

	"github.com/joyent/kosh/conch/types"
)

// POST /rack
func (c *Client) CreateRack(rack types.RackCreate) error {
	_, e := c.Rack().Post(rack).Send()
	return e
}

// GET /rack/:rack_id_or_name
func (c *Client) GetRackByName(name string) (rack types.Rack) {
	c.Rack(name).Receive(&rack)
	return
}

func (c *Client) GetRackByID(id types.UUID) (rack types.Rack) {
	c.Rack(id.String()).Receive(&rack)
	return
}

// POST /rack/:rack_id_or_name
func (c *Client) UpdateRack(id types.UUID, rack types.RackUpdate) error {
	_, e := c.Rack(id.String()).Post(rack).Send()
	return e
}

// DELETE /rack/:rack_id_or_name
func (c *Client) DeleteRack(id types.UUID) error {
	_, e := c.Rack(id.String()).Delete().Send()
	return e
}

// GET /rack/:rack_id_or_name/layout
func (c *Client) GetRackLayout(id types.UUID) (rack types.RackLayouts) {
	c.Rack(id.String()).Layout().Receive(&rack)
	return
}

// POST /rack/:rack_id_or_name/layout
func (c *Client) UpdateRackLayout(id types.UUID, layout types.RackLayoutUpdate) error {
	_, e := c.Rack(id.String()).Layout().Post(layout).Send()
	return e
}

func (c *Client) ReadRackLayoutUpdate(r io.Reader) (update types.RackLayoutUpdate) {
	json.NewDecoder(r).Decode(&update)
	return
}

// GET /rack/:rack_id_or_name/assignment
func (c *Client) GetRackAssignments(id types.UUID) (rack types.RackAssignments) {
	c.Rack(id.String()).Assignment().Receive(&rack)
	return
}

// POST /rack/:rack_id_or_name/assignment
func (c *Client) UpdateRackAssignments(id types.UUID, rack types.RackAssignmentUpdates) error {
	_, e := c.Rack(id.String()).Assignment().Post(rack).Send()
	return e
}

func (c *Client) ReadRackAssignmentUpdate(r io.Reader) (update types.RackAssignmentUpdates) {
	json.NewDecoder(r).Decode(&update)
	return
}

// DELETE /rack/:rack_id_or_name/assignment
func (c *Client) DeleteRackAssignments(id types.UUID, deletes types.RackAssignmentDeletes) error {
	_, e := c.Rack(id.String()).Assignment().Delete(deletes).Send()
	return e
}

// POST /rack/:rack_id_or_name/phase?rack_only=<0|1>
func (c *Client) UpdateRackPhase(id types.UUID, phase types.RackPhase, rackOnly bool) error {
	_, e := c.Rack(id.String()).Phase().Post(phase).Send()
	return e
}

// POST /rack/:rack_id_or_name/links
func (c *Client) UpdateRackLinks(id types.UUID, links types.RackLinks) error {
	_, e := c.Rack(id.String()).Links().Post(links).Send()
	return e
}

// DELETE /rack/:rack_id_or_name/links
func (c *Client) DeleteRackLinks(id types.UUID, phase types.RackLinks) error {
	_, e := c.Rack(id.String()).Links().Delete().Send()
	return e
}

// GET /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start
func (c *Client) GetSingleRackLayoutByRU(id types.UUID, ru string) (rack types.RackLayout) {
	c.Rack(id.String()).Layout(ru).Receive(&rack)
	return
}

func (c *Client) GetSingleRackLayoutByID(rackID, layoutID types.UUID) (rack types.RackLayout) {
	c.Rack(rackID.String()).Layout(layoutID.String()).Receive(&rack)
	return
}

// POST /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start
func (c *Client) UpdateSingleRackLayout(rackID, layoutID types.UUID, update types.RackLayoutUpdate) error {
	_, e := c.Rack(rackID.String()).Layout(layoutID.String()).Post(update).Send()
	return e
}

// DELETE /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start
func (c *Client) DeleteSingleRackLayout(rackID types.UUID, layoutID types.UUID) error {
	_, e := c.Rack(rackID.String()).Layout(layoutID.String()).Delete().Send()
	return e
}
