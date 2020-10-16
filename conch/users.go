package conch

import "github.com/joyent/kosh/conch/types"

// GET /user/me
func (c *Client) GetCurrentUser() (me types.UserDetailed) {
	c.Debug("GetCurrentUser()")
	me = c.GetUserByID("me")
	return
}

// POST /user/me/revoke
func (c *Client) RevokeCurrentUserCredentials() error {
	return c.RevokeUserCredentials("me")
}

// POST /user/me/password
func (c *Client) ChangeCurrentUserPassword(setting types.UserSetting) error {
	return c.ChangeUserPassword("me", setting)
}

// GET /user/me/settings
func (c *Client) GetCurrentUserSettings() (settings types.UserSettings) {
	c.User("me").Settings("").Receive(&settings)
	return
}

// POST /user/me/setting
func (c *Client) SetCurrentUserSettings(settings types.UserSettings) error {
	_, e := c.User("me").Settings().Post(settings).Send()
	return e
}

// GET /user/me/setting/:name
func (c *Client) GetCurrentUserSettingByName(name string) (setting types.UserSetting) {
	c.User("me").Settings(name).Receive(&setting)
	return
}

// POST /user/me/setting/:name
func (c *Client) SetCurrentUserSettingByName(name string, setting types.UserSetting) error {
	_, e := c.User("me").Settings(name).Post(setting).Send()
	return e
}

// DELETE /user/me/setting/:name
func (c *Client) DeleteCurrentUserSetting(name string) error {
	_, e := c.User("me").Settings(name).Delete().Send()
	return e
}

// GET /user/me/token
func (c *Client) GetCurrentUserTokens() (tokens types.UserTokens) {
	tokens = c.GetUserToken("me")
	return
}

// POST /user/me/token
func (c *Client) CreateCurrentUserToken(newToken types.NewUserToken) (token types.NewUserToken) {
	c.User("me").Token().Post(newToken).Receive(&token)
	return
}

// GET /user/me/token/:token_name
func (c *Client) GetCurrentUserTokenByName(name string) (token types.UserToken) {
	token = c.GetUserTokenByName("me", name)
	return
}

// DELETE /user/me/token/:token_name
func (c *Client) DeleteCurrentUserToken(name string) error {
	return c.DeleteUserToken("me", name)
}

// GET /user/:target_user_id_or_email
func (c *Client) GetUserByID(id string) (user types.UserDetailed) {
	c.User(id).Receive(&user)
	return
}

// POST /user/:target_user_id_or_email?send_mail=<1|0>
func (c *Client) UpdateUser(id string, update types.UpdateUser) error {
	_, e := c.User(id).Post(update).Send()
	return e
}

// DELETE /user/:target_user_id_or_email
func (c *Client) DeleteUser(id string) error {
	_, e := c.User(id).Delete().Send()
	return e
}

// POST /user/:target_user_id_or_email/revoke
func (c *Client) RevokeUserCredentials(id string) error {
	_, e := c.User(id).Revoke().Post("").Send()
	return e
}

// DELETE /user/:target_user_id_or_email/password
func (c *Client) ChangeUserPassword(id string, setting types.UserSetting) error {
	_, e := c.User(id).Password().Post(setting).Send()
	return e
}

// GET /user
func (c *Client) GetUsers() (me types.UsersDetailed) {
	c.User("").Receive(&me)
	return
}

// POST /user?send_mail=<1|0>
func (c *Client) CreateUser(newUser types.NewUser) (user types.NewUser) {
	c.User().Post(newUser).Receive(&user)
	return
}

// GET /user/:target_user_id_or_email/token
func (c *Client) GetUserToken(id string) (tokens types.UserTokens) {
	c.User(id).Token().Receive(&tokens)
	return
}

// GET /user/:target_user_id_or_email/token/:token_name
func (c *Client) GetUserTokenByName(id, name string) (token types.UserToken) {
	c.User(id).Token(name).Receive(&token)
	return
}

// DELETE /user/:target_user_id_or_email/token/:token_name
func (c *Client) DeleteUserToken(id, name string) error {
	_, e := c.User(id).Token(name).Delete().Send()
	return e
}
