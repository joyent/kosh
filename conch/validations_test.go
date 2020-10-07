package conch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/conch"
	"github.com/stretchr/testify/assert"
)

func TestValidationRoutes(t *testing.T) {
	tests := []struct {
		URL    string
		Method string
		Do     func(c *conch.Client)
	}{
		{
			URL:    "/validation/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetValidations() },
		},
		{
			URL:    "/validation/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetValidationByID("foo") },
		},
		{
			URL:    "/validation_plan/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetValidationPlans() },
		},
		{
			URL:    "/validation_plan/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetValidationPlanByID("foo") },
		},
		{
			URL:    "/validation_plan/foo/validation/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetValidationPlanValidations("foo") },
		},
		{
			URL:    "/validation_state/foo/",
			Method: "GET",
			Do:     func(c *conch.Client) { _ = c.GetValidationStateByID("foo") },
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
			test.Do(conch.New(ts.URL, "token", &logger{}))
			assert.True(t, seen, "saw the correct post to conch")
		})
	}
}
