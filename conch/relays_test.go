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

func TestRelayRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/relay/foo/register/",
			Method: "POST",
			Do:     func(c *conch.Client) { _ = c.RegisterRelay("foo", types.RegisterRelay{}) },
		},
		{
			URL:    "/relay/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetAllRelays() },
		},
		{
			URL:    "/relay/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRelayBySerial("foo") },
		},
		{
			URL:    "/relay/00000000-0000-0000-0000-000000000000/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetRelayByID(types.UUID{}) },
		},

		{
			URL:    "/relay/foo/",
			Method: "DELETE",
			Do:     func(c *conch.Client) { _ = c.DeleteRelay("foo") },
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
