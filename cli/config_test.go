package cli_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joyent/kosh/v3/cli"
	"github.com/joyent/kosh/v3/conch"
	"github.com/joyent/kosh/v3/conch/types"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	buffer := bytes.NewBufferString("")
	display := cli.NewConfig("test", "test").RenderTo(buffer)

	// TODO replace with fixtures
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	conch := conch.New(conch.API(ts.URL))

	tests := []struct {
		Name string
		Do   func()
	}{
		{
			Name: "GetCurrentUser()",
			Do:   func() { display(conch.GetCurrentUser()) },
		},

		{
			Name: "display(conch.FindDevicesByField(\"hostname\", hostname))",
			Do:   func() { display(conch.FindDevicesByField("hostname", "foo")) },
		},

		{
			Name: "display(conch.FindDevicesBySetting(key, value))",
			Do:   func() { display(conch.FindDevicesBySetting("foo", "bar")) },
		},

		{
			Name: "display(conch.FindDevicesByTag(key, value))",
			Do:   func() { display(conch.FindDevicesByTag("foo", "bar")) },
		},

		{
			Name: "	display(conch.GetAllBuildDevices(*buildNameArg))",
			Do: func() { display(conch.GetAllBuildDevices("build-000")) },
		},

		{
			Name: "	display(conch.GetAllBuildOrganizations(*buildNameArg))",
			Do: func() { display(conch.GetAllBuildOrganizations("build-000")) },
		},

		{
			Name: "	display(conch.GetAllBuilds())",
			Do: func() { display(conch.GetAllBuilds()) },
		},

		{
			Name: "	display(conch.GetAllDatacenterRooms(dc.ID))",
			Do: func() { display(conch.GetAllDatacenterRooms(types.UUID{})) },
		},

		{
			Name: "	display(conch.GetAllDatacenters())",
			Do: func() { display(conch.GetAllDatacenters()) },
		},

		{
			Name: "	display(conch.GetAllRackRoles())",
			Do: func() { display(conch.GetAllRackRoles()) },
		},

		{
			Name: "	display(conch.GetAllRelays())",
			Do: func() { display(conch.GetAllRelays()) },
		},

		{
			Name: "	display(conch.GetAllRoomRacks(room.ID))",
			Do: func() { display(conch.GetAllRoomRacks(types.UUID{})) },
		},

		{
			Name: "	display(conch.GetAllRooms())",
			Do: func() { display(conch.GetAllRooms()) },
		},

		{
			Name: "	display(conch.GetAllValidationPlans())",
			Do: func() { display(conch.GetAllValidationPlans()) },
		},

		{
			Name: "	display(conch.GetBuildRacks(*buildNameArg))",
			Do: func() { display(conch.GetBuildRacks("build-001")) },
		},

		{
			Name: "	display(conch.GetBuildUsers(*buildNameArg))",
			Do: func() { display(conch.GetBuildUsers("build-001")) },
		},

		{
			Name: "	display(conch.GetCurrentUser())",
			Do: func() { display(conch.GetCurrentUser()) },
		},

		{
			Name: "	display(conch.GetCurrentUserSettingByName(setting))",
			Do: func() { display(conch.GetCurrentUserSettingByName("foo")) },
		},

		{
			Name: "	display(conch.GetCurrentUserSettings())",
			Do: func() { display(conch.GetCurrentUserSettings()) },
		},

		{
			Name: "	display(conch.GetDeviceBySerial(*id))",
			Do: func() { display(conch.GetDeviceBySerial("serial")) },
		},

		{
			Name: "	display(conch.GetDeviceInterfaceByName(*id, name))",
			Do: func() { display(conch.GetDeviceInterfaceByName("foo", "device")) },
		},

		{
			Name: "	display(conch.GetDeviceLocation(*id))",
			Do: func() { display(conch.GetDeviceLocation("foo")) },
		},

		{
			Name: "	display(conch.GetDevicePhase(*id))",
			Do: func() { display(conch.GetDevicePhase("foo")) },
		},

		{
			Name: "	display(conch.GetDeviceSettingByName(*id, key))",
			Do: func() { display(conch.GetDeviceSettingByName("foo", "bar")) },
		},

		{
			Name: "	display(conch.GetDeviceSettings(*id))",
			Do: func() { display(conch.GetDeviceSettings("foo")) },
		},

		{
			Name: "	display(conch.GetDeviceTagByName(*id, name))",
			Do: func() { display(conch.GetDeviceTagByName("foo", "bar")) },
		},

		{
			Name: "	display(conch.GetDeviceTags(*id))",
			Do: func() { display(conch.GetDeviceTags("foo")) },
		},

		{
			Name: "	display(conch.GetDeviceValidationStates(*id))",
			Do: func() { display(conch.GetDeviceValidationStates("foo")) },
		},

		{
			Name: "	display(conch.GetHardwareProductByID(*name))",
			Do: func() { display(conch.GetHardwareProductByID("foo")) },
		},

		{
			Name: "	display(conch.GetHardwareProducts())",
			Do: func() { display(conch.GetHardwareProducts()) },
		},

		{
			Name: "	display(conch.GetHardwareVendors())",
			Do: func() { display(conch.GetAllHardwareVendors()) },
		},

		{
			Name: "	display(conch.GetOrganizationByID(o.ID))",
			Do: func() { display(conch.GetOrganizationByID(types.UUID{})) },
		},

		{
			Name: "	display(conch.GetOrganizations())",
			Do: func() { display(conch.GetAllOrganizations()) },
		},

		{
			Name: "	display(conch.GetRackAssignments(rack.ID))",
			Do: func() { display(conch.GetRackAssignments(types.UUID{})) },
		},

		{
			Name: "	display(conch.GetRackLayout(rack.ID))",
			Do: func() { display(conch.GetRackLayout(types.UUID{})) },
		},

		{
			Name: "	display(conch.GetRoomByID(room.ID))",
			Do: func() { display(conch.GetRoomByID(types.UUID{})) },
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			defer ts.Close()
			test.Do()
			assert.NotEmpty(t, buffer.String())
			buffer.Reset()
		})
	}
}
