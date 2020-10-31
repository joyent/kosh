package conch

import "github.com/joyent/kosh/conch/types"

// GetAllRooms (GET /room) returns a list of all datacenter rooms
func (c *Client) GetAllRooms() (rooms types.DatacenterRoomsDetailed) {
	c.Room().Receive(rooms)
	return
}

// CreateRoom (POST /room) creates a new datacenter room
func (c *Client) CreateRoom(room types.DatacenterRoomCreate) error {
	_, e := c.Room().Post(room).Send()
	return e
}

// GetRoomByAlias (GET /room/:datacenter_room_id_or_alias) retrieves the
// datacenter room with the given alias
func (c *Client) GetRoomByAlias(alias string) (room types.DatacenterRoomDetailed) {
	c.Room(alias).Receive(room)
	return
}

// GetRoomByID (GET /room/:datacenter_room_id_or_alias) retrieves the
// datacenter room with the given UUID
func (c *Client) GetRoomByID(id types.UUID) (room types.DatacenterRoomDetailed) {
	c.Room(id.String()).Receive(room)
	return
}

// UpdateRoom (POST /room/:datacenter_room_id_or_alias)
// updates the datacenter room with the given UUID
func (c *Client) UpdateRoom(id types.UUID, room types.DatacenterRoomUpdate) error {
	_, e := c.Room(id.String()).Post(room).Send()
	return e
}

// DeleteRoom (DELETE /room/:datacenter_room_id_or_alias)
// removes the datacenter room with the given UUID
func (c *Client) DeleteRoom(id types.UUID) error {
	_, e := c.Room(id.String()).Delete().Send()
	return e
}

// GetAllRoomRacks (GET /room/:datacenter_room_id_or_alias/rack)
// returns a list of all racks in the room with the given UUID
func (c *Client) GetAllRoomRacks(id types.UUID) (racks types.Racks) {
	c.Room(id.String()).Rack().Receive(racks)
	return
}

// GetRoomRackByName (GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name)
// returns the specific rack with the given name in the room with the given UUID
func (c *Client) GetRoomRackByName(id types.UUID, name string) (rack types.Rack) {
	c.Room(id.String()).Rack(name).Receive(rack)
	return
}

// GetRoomRackByID (GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name)
// returns the specific rack with the given UUID in the room with the given UUID
func (c *Client) GetRoomRackByID(roomID, rackID types.UUID) (rack types.Rack) {
	c.Room(roomID.String()).Rack(rackID.String()).Receive(rack)
	return
}

// UpdateRoomRack (POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name)
// update the rack with the given UUID in the room with the given UUID
func (c *Client) UpdateRoomRack(roomID, rackID types.UUID, rack types.RackUpdate) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Post(rack).Send()
	return e
}

// DeleteRoomRack (DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name)
// remove the rack with the given UUID in the room with the given UUID
func (c *Client) DeleteRoomRack(roomID, rackID types.UUID) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Delete().Send()
	return e
}

// GetRoomRackLayout (GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout)
// get the rack layout for the rack in the given room
func (c *Client) GetRoomRackLayout(roomID, rackID types.UUID) (rack types.RackLayouts) {
	c.Room(roomID.String()).Rack(roomID.String()).Layout().Receive(rack)
	return
}

// UpdateRoomRackLayout (POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout)
// update the rack layout for the rack in the room
func (c *Client) UpdateRoomRackLayout(roomID, rackID types.UUID, layout types.RackLayoutUpdate) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Layout().Post(layout).Send()
	return e
}

// GetRoomRackAssignments (GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/assignment)
// retrieve the rack assignments for the rack in the given room
func (c *Client) GetRoomRackAssignments(roomID, rackID types.UUID) (rack types.RackAssignments) {
	c.Room(roomID.String()).Rack(rackID.String()).Assignment().Receive(rack)
	return
}

// UpdateRoomRackAssignments (POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/assignment)
// update the rack assignments for the rack in the given room
func (c *Client) UpdateRoomRackAssignments(roomID, rackID types.UUID, rack types.RackAssignmentUpdates) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Assignment().Post(rack).Send()
	return e
}

// DeleteRoomRackAssignments (DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/assignment)
// remove the rack assignments for the rack in the given room
func (c *Client) DeleteRoomRackAssignments(roomID, rackID types.UUID, deletes types.RackAssignmentDeletes) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Assignment().Delete(deletes).Send()
	return e
}

// UpdateRoomRackPhase (POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/phase?rack_only=<0|1>)
// update the rack phase for the rack in the given room
func (c *Client) UpdateRoomRackPhase(roomID, rackID types.UUID, phase types.RackPhase, rackOnly bool) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Phase().Post(phase).Send()
	return e
}

// UpdateRoomRackLinks (POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/links)
// update the rack links for the rack in the given room
func (c *Client) UpdateRoomRackLinks(roomID, rackID types.UUID, links types.RackLinks) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Links().Post(links).Send()
	return e
}

// DeleteRoomRackLinks (DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/links)
// remove the rack links for the rack in the given room
func (c *Client) DeleteRoomRackLinks(roomID, rackID types.UUID, phase types.RackLinks) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Links().Delete().Send()
	return e
}

// GetSingleRoomRackLayoutByRU (GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// get the rack layout for a single named RU in a rack in a room
func (c *Client) GetSingleRoomRackLayoutByRU(roomID, rackID types.UUID, ru string) (rack types.RackLayout) {
	c.Room(roomID.String()).Rack(rackID.String()).Layout(ru).Receive(rack)
	return
}

// GetSingleRoomRackLayoutByID (GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// get the rack layout for given UUID in a rack in a room
func (c *Client) GetSingleRoomRackLayoutByID(roomID, rackID, layoutID types.UUID) (rack types.RackLayout) {
	c.Room(rackID.String()).Rack(rackID.String()).Layout(layoutID.String()).Receive(rack)
	return
}

// UpdateSingleRoomRackLayout (POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// update the rack layout for given UUID in a rack in a room
func (c *Client) UpdateSingleRoomRackLayout(roomID, rackID, layoutID types.UUID, update types.RackLayoutUpdate) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Layout(layoutID.String()).Post(update).Send()
	return e
}

// DeleteSingleRoomRackLayout (DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start)
// remove the rack layout for given UUID in a rack in a room
func (c *Client) DeleteSingleRoomRackLayout(roomID, rackID, layoutID types.UUID) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Layout(layoutID.String()).Delete().Send()
	return e
}
