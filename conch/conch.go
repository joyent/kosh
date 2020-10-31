/*
Package conch provivides a wrapper around the  Conch-API web service. For more
inforamtion about the Conch-API, including documentation for the endpoints that
this library wraps, please see: https://joyent.github.io/conch-api/
*/
package conch

import "github.com/joyent/kosh/v3/conch/types"

// Ping ( GET /ping ) checks to see if the API server is online and reachable.
// See also https://joyent.github.io/conch-api/modules/Conch::Route#get-ping
func (c *Client) Ping() (ping types.Ping) {
	c.Path("ping").Receive(&ping)
	return
}

// Version ( GET /version ) returns the conch-api version See also
// https://joyent.github.io/conch-api/modules/Conch::Route#get-version
func (c *Client) Version() (version types.Version) {
	c.Path("version").Receive(&version)
	return
}

// Login ( POST /login ) takes a username and password and returns a session
// token from the server. See also also
// https://joyent.github.io/conch-api/modules/Conch::Route#post-login
func (c *Client) Login(user, pass string) (token types.LoginToken) {
	c.Path("login").Post(types.Login{Email: types.EmailAddress(user), Password: types.NonEmptyString(pass)}).Receive(&token)
	return
}

// Logout ( POST /logout ) terminates a login token on the server. See also
// https://joyent.github.io/conch-api/modules/Conch::Route#post-logout
func (c *Client) Logout() error {
	_, e := c.Path("logout").Post("").Send()
	return e
}

// RefreshToken ( POST /refresh_token ) uses the current login token to
// generate a *new* login token. See also
// https://joyent.github.io/conch-api/modules/Conch::Route#post-refresh_token
func (c *Client) RefreshToken() (token types.LoginToken) {
	c.Path("refresh_token").Post("").Receive(&token)
	return
}

// IsSysAdmin is a helper function that uses GetCurrentUser to figure out if
// the current user is a system admin or not.
func (c *Client) IsSysAdmin() bool {
	return c.GetCurrentUser().IsAdmin
}
