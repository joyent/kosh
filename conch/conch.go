package conch

import "github.com/joyent/kosh/conch/types"

// GET /ping
func (c *Client) Ping() (ping types.Ping) {
	c.Path("ping").Receive(ping)
	return
}

// GET /version
func (c *Client) Version() (version types.Version) {
	c.Path("version").Receive(version)
	return
}

// POST /login
func (c *Client) Login(user, pass string) (token types.LoginToken) {
	c.Path("login").Post(types.Login{Email: types.EmailAddress(user), Password: types.NonEmptyString(pass)}).Receive(token)
	return
}

// POST /logout
func (c *Client) Logout() error {
	_, e := c.Path("logout").Post("").Send()
	return e
}

// POST /refresh_token
func (c *Client) RefreshToken() (token types.LoginToken) {
	c.Path("refresh_token").Post("").Receive(token)
	return
}
