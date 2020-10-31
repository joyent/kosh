package conch

import "github.com/joyent/kosh/v3/conch/types"

// GetCurrentUser (GET /user/me) retrieves the user associated with the current
// authentication
func (c *Client) GetCurrentUser() (me types.UserDetailed) {
	me = c.GetUserByEmail("me")
	return
}

// RevokeCurrentUserCredentials (POST /user/me/revoke) reokes the
// authentication (login) credentials for the current user.
// DOES NOT AFFECT API TOKENS
func (c *Client) RevokeCurrentUserCredentials() error {
	return c.RevokeUserCredentials("me")
}

// ChangeCurrentUserPassword (POST /user/me/password)
// updates the password for the current user
func (c *Client) ChangeCurrentUserPassword(setting types.UserSetting) error {
	return c.ChangeUserPassword("me", setting)
}

// GetCurrentUserSettings (GET /user/me/settings) gets the settings for the
// current user
func (c *Client) GetCurrentUserSettings() (settings types.UserSettings) {
	c.User("me").Settings("").Receive(&settings)
	return
}

// SetCurrentUserSettings (POST /user/me/setting) updates the settings for the
// current user
func (c *Client) SetCurrentUserSettings(settings types.UserSettings) error {
	_, e := c.User("me").Settings().Post(settings).Send()
	return e
}

// GetCurrentUserSettingByName (GET /user/me/setting/:name) retrieves a single
// user setting
func (c *Client) GetCurrentUserSettingByName(name string) (setting types.UserSetting) {
	c.User("me").Settings(name).Receive(&setting)
	return
}

// SetCurrentUserSettingByName (POST /user/me/setting/:name) sets a single
// users setting
func (c *Client) SetCurrentUserSettingByName(name string, setting types.UserSetting) error {
	_, e := c.User("me").Settings(name).Post(setting).Send()
	return e
}

// DeleteCurrentUserSetting (DELETE /user/me/setting/:name) removes a single
// user setting by name
func (c *Client) DeleteCurrentUserSetting(name string) error {
	_, e := c.User("me").Settings(name).Delete().Send()
	return e
}

// GetCurrentUserTokens (GET /user/me/token) returns the list of API tokens for the current user
func (c *Client) GetCurrentUserTokens() (tokens types.UserTokens) {
	tokens = c.GetUserTokens("me")
	return
}

// CreateCurrentUserToken (POST /user/me/token) creates a new API token for the
// current user. This is the only time the actual token string will be readable
func (c *Client) CreateCurrentUserToken(newToken types.NewUserToken) (token types.NewUserToken) {
	c.User("me").Token().Post(newToken).Receive(&token)
	return
}

// GetCurrentUserTokenByName (GET /user/me/token/:token_name) returns the
// information for a single API token for the current user. The token string
// itself is not readable.
func (c *Client) GetCurrentUserTokenByName(name string) (token types.UserToken) {
	token = c.GetUserTokenByName("me", name)
	return
}

// DeleteCurrentUserToken (DELETE /user/me/token/:token_name) removes the token
// with the given name for the current user
func (c *Client) DeleteCurrentUserToken(name string) error {
	return c.DeleteUserToken("me", name)
}

// GetUserByEmail (GET /user/:target_user_id_or_email) retrieves the user with
// the given email
func (c *Client) GetUserByEmail(email string) (user types.UserDetailed) {
	c.User(email).Receive(&user)
	return
}

// GetUserByID (GET /user/:target_user_id_or_email) retrieves the user with
// the given UUID
func (c *Client) GetUserByID(id types.UUID) (user types.UserDetailed) {
	c.User(id.String()).Receive(&user)
	return
}

// UpdateUser (POST /user/:target_user_id_or_email?send_mail=<1|0>) will update the
// user with the given email. Optionally notify the user via email.
// BUG(perigrin): sendEmail is currently not implemented
func (c *Client) UpdateUser(email string, update types.UpdateUser, sendEmail bool) error {
	_, e := c.User(email).Post(update).Send()
	return e
}

// DeleteUser (DELETE /user/:target_user_id_or_email) will remove the user with the
// given email
func (c *Client) DeleteUser(email string) error {
	_, e := c.User(email).Delete().Send()
	return e
}

// RevokeUserCredentials (POST /user/:target_user_id_or_email/revoke) will
// revoke the authentication credentials for the user with hte given email.
// DOES NOT AFFECT API TOKENS FOR THE USER
func (c *Client) RevokeUserCredentials(email string) error {
	_, e := c.User(email).Revoke().Post("").Send()
	return e
}

// ChangeUserPassword (DELETE /user/:target_user_id_or_email/password) triggers
// the password change mechanisim for the user with the given email.
func (c *Client) ChangeUserPassword(email string, setting types.UserSetting) error {
	_, e := c.User(email).Password().Post(setting).Send()
	return e
}

// GetAllUsers (GET /user) retrieves a list of all users
func (c *Client) GetAllUsers() (me types.UsersDetailed) {
	c.User("").Receive(&me)
	return
}

// CreateUser (POST /user?send_mail=<1|0>) create a new user in teh system and
// optionally send them an email notification.
// BUG(perigrin): sendEmail isn't implemented
func (c *Client) CreateUser(newUser types.NewUser, sendEmail bool) (user types.NewUser) {
	c.User().Post(newUser).Receive(&user)
	return
}

// GetUserTokens (GET /user/:target_user_id_or_email/token) retrieves the list
// of API tokens for the given user.
func (c *Client) GetUserTokens(email string) (tokens types.UserTokens) {
	c.User(email).Token().Receive(&tokens)
	return
}

// GetUserTokenByName (GET /user/:target_user_id_or_email/token/:token_name)
// retrieves a single named API token for the given user
func (c *Client) GetUserTokenByName(email, name string) (token types.UserToken) {
	c.User(email).Token(name).Receive(&token)
	return
}

// DeleteUserToken (DELETE /user/:target_user_id_or_email/token/:token_name)
// removes a named API token for the given user
func (c *Client) DeleteUserToken(email, name string) error {
	_, e := c.User(email).Token(name).Delete().Send()
	return e
}
