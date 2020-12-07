package conch

import (
	"encoding/json"
	"io"

	"github.com/joyent/kosh/conch/types"
)

// ReadRackCreate takes an io.Reader and returns a RackCreate
// struct suitable for CreateRack
func (c *Client) ReadRackCreate(r io.Reader) (rackCreate types.RackCreate, e error) {
	e = json.NewDecoder(r).Decode(&rackCreate)
	return
}

// CreateRack (POST /rack) creates a new rack
func (c *Client) CreateRack(rack types.RackCreate) error {
	_, e := c.Rack().Post(rack).Send()
	return e
}

// GetRackByName (GET /rack/:rack_id_or_name) returns a rack with the given
// name
func (c *Client) GetRackByName(name string) (rack types.Rack, e error) {
	_, e = c.Rack(name).Receive(&rack)
	return
}

// GetRackByID (GET /rack/:rack_id_or_name) returns a rack with the given
// UUID
func (c *Client) GetRackByID(id types.UUID) (rack types.Rack, e error) {
	_, e = c.Rack(id.String()).Receive(&rack)
	return
}

// UpdateRack (POST /rack/:rack_id_or_name) updates a rack with the given UUID
func (c *Client) UpdateRack(id types.UUID, rack types.RackUpdate) error {
	_, e := c.Rack(id.String()).Post(rack).Send()
	return e
}

// DeleteRack (DELETE /rack/:rack_id_or_name) removes the rack with the given UUID
func (c *Client) DeleteRack(id types.UUID) error {
	_, e := c.Rack(id.String()).Delete().Send()
	return e
}

// GetRackLayout (GET /rack/:rack_id_or_name/layout) retrieves the rack layout
// for the rack with the given UUID
func (c *Client) GetRackLayout(id types.UUID) (rack types.RackLayouts, e error) {
	_, e = c.Rack(id.String()).Layout().Receive(&rack)
	return
}

// UpdateRackLayout (POST /rack/:rack_id_or_name/layout) updates the rack
// layout for the rack with the given UUID
func (c *Client) UpdateRackLayout(id types.UUID, layout []types.RackLayoutUpdate) error {
	_, e := c.Rack(id.String()).Layout().Post(layout).Send()
	return e
}

// ReadRackLayoutUpdate takes an io.Reader and returns a RackLayoutUpdate
// struct suitable for UpdateRackLayout
func (c *Client) ReadRackLayoutUpdate(r io.Reader) (update []types.RackLayoutUpdate, e error) {
	e = json.NewDecoder(r).Decode(&update)
	return
}

// GetRackAssignments (GET /rack/:rack_id_or_name/assignment) gets the current
// assignments for the rack
func (c *Client) GetRackAssignments(id types.UUID) (rack types.RackAssignments, e error) {
	_, e = c.Rack(id.String()).Assignment().Receive(&rack)
	return
}

// UpdateRackAssignments (POST /rack/:rack_id_or_name/assignment) updates the
// assignments for the rack
func (c *Client) UpdateRackAssignments(id types.UUID, rack types.RackAssignmentUpdates) error {
	_, e := c.Rack(id.String()).Assignment().Post(rack).Send()
	return e
}

// ReadRackAssignmentUpdate takes an io reader and returns a RackAssignmentUpdates
// struct suitable for UpdateRackAssignments
func (c *Client) ReadRackAssignmentUpdate(r io.Reader) (update types.RackAssignmentUpdates, e error) {
	e = json.NewDecoder(r).Decode(&update)
	return
}

// DeleteRackAssignments (DELETE /rack/:rack_id_or_name/assignment) removes all
// of the current assignments for the rack
func (c *Client) DeleteRackAssignments(id types.UUID, deletes types.RackAssignmentDeletes) error {
	_, e := c.Rack(id.String()).Assignment().Delete(deletes).Send()
	return e
}

// UpdateRackPhase (POST /rack/:rack_id_or_name/phase?rack_only=<0|1>) updates
// the rack phase and by default all the devices in the rack
// BUG(perigrin): rackOnly currently isn't implemented
func (c *Client) UpdateRackPhase(id types.UUID, phase types.RackPhase, rackOnly bool) error {
	_, e := c.Rack(id.String()).Phase().Post(phase).Send()
	return e
}

// UpdateRackLinks (POST /rack/:rack_id_or_name/links) updates the links
// associated with the rack
func (c *Client) UpdateRackLinks(id types.UUID, links types.RackLinks) error {
	_, e := c.Rack(id.String()).Links().Post(links).Send()
	return e
}

// DeleteRackLinks (DELETE /rack/:rack_id_or_name/links) removes the links
// associated with the rack
func (c *Client) DeleteRackLinks(id types.UUID, phase types.RackLinks) error {
	_, e := c.Rack(id.String()).Links().Delete().Send()
	return e
}

// GetSingleRackLayoutByRU (GET /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// returns a single layout for the given RU
func (c *Client) GetSingleRackLayoutByRU(id types.UUID, ru string) (rack types.RackLayout, e error) {
	_, e = c.Rack(id.String()).Layout(ru).Receive(&rack)
	return
}

// GetSingleRackLayoutByID (GET /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// returns a single layout for the given UUID
func (c *Client) GetSingleRackLayoutByID(rackID, layoutID types.UUID) (rack types.RackLayout, e error) {
	_, e = c.Rack(rackID.String()).Layout(layoutID.String()).Receive(&rack)
	return
}

// UpdateSingleRackLayout (POST /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// updates a single rack layout with the given UUID
func (c *Client) UpdateSingleRackLayout(rackID, layoutID types.UUID, update types.RackLayoutUpdate) error {
	_, e := c.Rack(rackID.String()).Layout(layoutID.String()).Post(update).Send()
	return e
}

// DeleteSingleRackLayout (DELETE /rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// removes the rack layout with the given UUID
func (c *Client) DeleteSingleRackLayout(rackID types.UUID, layoutID types.UUID) error {
	_, e := c.Rack(rackID.String()).Layout(layoutID.String()).Delete().Send()
	return e
}
