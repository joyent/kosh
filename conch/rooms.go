package conch

import "github.com/joyent/kosh/conch/types"

// GET /room
func (c *Client) GetAllRooms() (rooms types.DatacenterRoomsDetailed) {
	c.Room().Receive(rooms)
	return
}

// POST /room
func (c *Client) CreateRoom(room types.DatacenterRoomCreate) error {
	_, e := c.Room().Post(room).Send()
	return e
}

// GET /room/:datacenter_room_id_or_alias
func (c *Client) GetRoomByAlias(alias string) (room types.DatacenterRoomDetailed) {
	c.Room(alias).Receive(room)
	return
}

func (c *Client) GetRoomByID(id types.UUID) (room types.DatacenterRoomDetailed) {
	c.Room(id.String()).Receive(room)
	return
}

// POST /room/:datacenter_room_id_or_alias
func (c *Client) UpdateRoom(id types.UUID, room types.DatacenterRoomUpdate) error {
	_, e := c.Room(id.String()).Post(room).Send()
	return e
}

// DELETE /room/:datacenter_room_id_or_alias
func (c *Client) DeleteRoom(id types.UUID) error {
	_, e := c.Room(id.String()).Delete().Send()
	return e
}

// GET /room/:datacenter_room_id_or_alias/rack
func (c *Client) GetAllRoomRacks(id types.UUID) (racks types.Racks) {
	c.Room(id.String()).Rack().Receive(racks)
	return
}

// GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name
func (c *Client) GetRoomRackByName(id types.UUID, name string) (rack types.Rack) {
	c.Room(id.String()).Rack(name).Receive(rack)
	return
}

func (c *Client) GetRoomRackByID(roomID, rackID types.UUID) (rack types.Rack) {
	c.Room(roomID.String()).Rack(rackID.String()).Receive(rack)
	return
}

// POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name
func (c *Client) UpdateRoomRack(roomID, rackID types.UUID, rack types.RackUpdate) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Post(rack).Send()
	return e
}

// DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name
func (c *Client) DeleteRoomRack(roomID, rackID types.UUID) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Delete().Send()
	return e
}

// GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout
func (c *Client) GetRoomRackLayout(roomID, rackID types.UUID) (rack types.RackLayouts) {
	c.Room(roomID.String()).Rack(roomID.String()).Layout().Receive(rack)
	return
}

// POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout
func (c *Client) UpdateRoomRackLayout(roomID, rackID types.UUID, layout types.RackLayoutUpdate) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Layout().Post(layout).Send()
	return e
}

// GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/assignment
func (c *Client) GetRoomRackAssignments(roomID, rackID types.UUID) (rack types.RackAssignments) {
	c.Room(roomID.String()).Rack(rackID.String()).Assignment().Receive(rack)
	return
}

// POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/assignment
func (c *Client) UpdateRoomRackAssignments(roomID, rackID types.UUID, rack types.RackAssignmentUpdates) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Assignment().Post(rack).Send()
	return e
}

// DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/assignment
func (c *Client) DeleteRoomRackAssignments(roomID, rackID types.UUID, deletes types.RackAssignmentDeletes) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Assignment().Delete(deletes).Send()
	return e
}

// POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/phase?rack_only=<0|1>
func (c *Client) UpdateRoomRackPhase(roomID, rackID types.UUID, phase types.RackPhase, rackOnly bool) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Phase().Post(phase).Send()
	return e
}

// POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/links
func (c *Client) UpdateRoomRackLinks(roomID, rackID types.UUID, links types.RackLinks) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Links().Post(links).Send()
	return e
}

// DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/links
func (c *Client) DeleteRoomRackLinks(roomID, rackID types.UUID, phase types.RackLinks) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Links().Delete().Send()
	return e
}

// GET /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start
func (c *Client) GetSingleRoomRackLayoutByRU(roomID, rackID types.UUID, ru string) (rack types.RackLayout) {
	c.Room(roomID.String()).Rack(rackID.String()).Layout(ru).Receive(rack)
	return
}

func (c *Client) GetSingleRoomRackLayoutByID(roomID, rackID, layoutID types.UUID) (rack types.RackLayout) {
	c.Room(rackID.String()).Rack(rackID.String()).Layout(layoutID.String()).Receive(rack)
	return
}

// POST /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start
func (c *Client) UpdateSingleRoomRackLayout(roomID, rackID, layoutID types.UUID, update types.RackLayoutUpdate) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Layout(layoutID.String()).Post(update).Send()
	return e
}

// DELETE /room/:datacenter_room_id_or_alias/rack/:rack_id_or_name/layout/:layout_id_or_rack_unit_start
func (c *Client) DeleteSingleRoomRackLayout(roomID, rackID, layoutID types.UUID) error {
	_, e := c.Room(roomID.String()).Rack(rackID.String()).Layout(layoutID.String()).Delete().Send()
	return e
}
