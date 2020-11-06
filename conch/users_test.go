package conch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/user/me/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetCurrentUser() },
		},
		{
			URL:    "/user/me/revoke/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.RevokeCurrentUserCredentials() },
		},
		{
			URL:    "/user/me/password/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.ChangeCurrentUserPassword(types.UserSetting("password"))
			},
		},
		{
			URL:    "/user/me/settings/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetCurrentUserSettings() },
		},
		{
			URL:    "/user/me/settings/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetCurrentUserSettingByName("foo") },
		},
		{
			URL:    "/user/me/settings/foo/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.SetCurrentUserSettingByName("foo", "bar") },
		},
		{
			URL:    "/user/me/settings/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteCurrentUserSetting("foo") },
		},
		{
			URL:    "/user/me/token/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetCurrentUserTokens() },
		},
		{
			URL:    "/user/me/token/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.CreateCurrentUserToken(types.NewUserToken{"foo"}) },
		},
		{
			URL:    "/user/me/token/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetCurrentUserTokenByName("foo") },
		},
		{
			URL:    "/user/me/token/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteCurrentUserToken("foo") },
		},
		{
			URL:    "/user/foo@example.com/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetUserByEmail("foo@example.com") },
		},
		{
			URL:    "/user/foo/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.UpdateUser("foo", types.UpdateUser{}, false) },
		},
		{
			URL:    "/user/foo/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.UpdateUser("foo", types.UpdateUser{}, true) },
		},
		{
			URL:    "/user/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteUser("foo") },
		},
		{
			URL:    "/user/foo/revoke/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.RevokeUserCredentials("foo") },
		},
		{
			URL:    "/user/foo/password/",
			Method: "POST",
			Do: func(c *conch.Client) {
				c.ChangeUserPassword("foo", types.UserSetting("password"))
			},
		},
		{
			URL:    "/user/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetAllUsers() },
		},
		{
			URL:    "/user/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.CreateUser(types.NewUser{}, false) },
		},
		{
			URL:    "/user/",
			Method: "POST",
			Do:     func(c *conch.Client) { c.CreateUser(types.NewUser{}, true) },
		},

		{
			URL:    "/user/foo/token/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetUserTokens("foo") },
		},
		{
			URL:    "/user/foo/token/bar/",
			Method: "GET",
			Do:     func(c *conch.Client) { c.GetUserTokenByName("foo", "bar") },
		},
		{
			URL:    "/user/foo/token/bar/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { c.DeleteUserToken("foo", "bar") },
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s %s", test.Method, test.URL)
		t.Run(name, func(t *testing.T) {
			seen := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.URL, r.URL.String())
				assert.Equal(t, test.Method, r.Method)
				seen = true

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()
			test.Do(conch.New(conch.API(ts.URL)))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}

func TestUserAPIIntergration(t *testing.T) {
	c := NewTestClient("fixtures/conch-v3/user")

	t.Run("me", func(t *testing.T) {
		c.GetCurrentUser()
	})

	t.Run("me-settings", func(t *testing.T) {
		c.GetCurrentUserSettings()
	})
}
