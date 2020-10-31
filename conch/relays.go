package conch

import "github.com/joyent/kosh/v3/conch/types"

// RegisterRelay (POST /relay/:relay_serial_number/register)
// registers a relay with the given serial number
func (c *Client) RegisterRelay(serial string, relay types.RegisterRelay) error {
	_, e := c.Relay(serial).Register().Post(relay).Send()
	return e
}

// GetAllRelays (GET /relay) returns a list of all relays
func (c *Client) GetAllRelays() (relays types.Relays) {
	c.Relay().Receive(&relays)
	return
}

// GetRelayBySerial (GET /relay/:relay_id_or_serial_number) retrieves a relay
// with the given serial number
func (c *Client) GetRelayBySerial(serial string) (relay types.Relay) {
	c.Relay(serial).Receive(&relay)
	return
}

// GetRelayByID (GET /relay/:relay_id_or_serial_number) retrieves a relay
// with the given UUID
func (c *Client) GetRelayByID(id types.UUID) (relay types.Relay) {
	c.Relay(id.String()).Receive(&relay)
	return
}

// DeleteRelay (DELETE /relay/:relay_id_or_serial_number) removes a relay
// BUG(perigrin): id here should be a UUID not a string
func (c *Client) DeleteRelay(id string) error {
	_, e := c.Relay(id).Delete().Send()
	return e
}
