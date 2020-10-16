package conch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/logger"
	"github.com/stretchr/testify/assert"
)

type config struct {
	url    string
	token  string
	logger logger.Logger
}

func (c config) GetURL() string           { return c.url }
func (c config) GetToken() string         { return c.token }
func (c config) GetLogger() logger.Logger { return c.logger }

func newConfig(URL string) config {
	return config{
		URL,
		"token",
		logger.New(),
	}
}

func TestDefaultRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/ping/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.Ping() },
		},
		{
			URL:    "/version/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.Version() },
		},
		{
			URL:    "/login/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.Login("foo", "bar") },
		},
		{
			URL:    "/logout/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.Logout() },
		},
		{
			URL:    "/refresh_token/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.RefreshToken() },
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
			test.Do(conch.New(newConfig(ts.URL)))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}
