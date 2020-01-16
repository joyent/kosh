package main

import (
	"testing"
)

func TestBuildAPIintegration(t *testing.T) {
	defer errorHandler()

	setupAPIClient()
	r := setupRecorder("fixtures/conch-v3/builds")
	defer r() // Make sure recorder is stopped once done with it

	var build Build
	t.Run("create a build", func(t *testing.T) {
		defer errorHandler()
		fake := newTestBuild()
		build = API.Builds().Create(
			fake.Name,
			fake.Description,
			[]map[string]string{{"email": "conch@example.com"}},
		)
	})

	t.Run("get all builds", func(t *testing.T) {
		defer errorHandler()
		_ = API.Builds().GetAll()
	})

	t.Run("get a specific build", func(t *testing.T) {
		defer errorHandler()
		_ = API.Builds().Get(build.ID)
	})
}
