package conch

import "github.com/joyent/kosh/conch/types"

// POST /relay/:relay_serial_number/register
func (c *Client) RegisterRelay(id string, relay types.RegisterRelay) error {
	_, e := c.Relay(id).Register().Post(relay).Send()
	return e
}

// GET /relay
func (c *Client) GetAllRelays() (relays types.Relays) {
	c.Relay().Receive(&relays)
	return
}

// GET /relay/:relay_id_or_serial_number
func (c *Client) GetRelayBySerial(serial string) (relay types.Relay) {
	c.Relay(serial).Receive(&relay)
	return
}

func (c *Client) GetRelayByID(id types.UUID) (relay types.Relay) {
	c.Relay(id.String()).Receive(&relay)
	return
}

// DELETE /relay/:relay_id_or_serial_number
func (c *Client) DeleteRelay(id string) error {
	_, e := c.Relay(id).Delete().Send()
	return e
}
