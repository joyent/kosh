package main

import (
	"math/rand"
	"reflect"

	"github.com/bxcodec/faker"
	"github.com/gofrs/uuid"
)

// currently ValidationPlans are immutable in every API
// changing this is part of the plans for v3.1
// but for now we just inline one of the plan IDs
const validationPlanID = "a30ab8b2-8a9e-4e51-8bb0-92862abd8b54"

func init() {
	faker.AddProvider("uuid", func(v reflect.Value) (interface{}, error) {
		return uuid.NewV4()
	})

	faker.AddProvider("rack_unit_size", func(v reflect.Value) (interface{}, error) {
		return rand.Intn(2) + 1, nil
	})

	lastRSU := 0
	faker.AddProvider("rack_unit_start", func(v reflect.Value) (interface{}, error) {
		if lastRSU >= 60 {
			lastRSU = 0
		}
		lastRSU += 1
		return lastRSU, nil
	})

	faker.AddProvider("rack_size", func(v reflect.Value) (interface{}, error) {
		return 60, nil
	})
}

type Fixture struct {
	dc             Datacenter
	room           Room
	role           RackRole
	rack           Rack
	rackLayout     RackLayout
	hardwareVendor HardwareVendor
	switchProduct  HardwareProduct
	serverProduct  HardwareProduct
	validationPlan ValidationPlan
	build          Build

	// reset function
	reset func()
}

func (f *Fixture) addReset(reset func()) *Fixture {
	next := f.reset
	f.reset = func() {
		reset()
		if next != nil {
			next()
		}
	}
	return f
}

func (f *Fixture) setupBuild() *Fixture {
	if f.build.ID != (uuid.UUID{}) {
		return f
	}

	mock := newTestBuild()
	f.build = API.Builds().Create(
		mock.Name,
		mock.Description,
		[]map[string]string{{"email": "conch@example.com"}},
	)

	return f
}

func (f *Fixture) setupHardwareVendor() *Fixture {
	if f.hardwareVendor != (HardwareVendor{}) {
		return f
	}

	//mock := newTestHardwareVendor()
	f.hardwareVendor = API.Hardware().FindOrCreateVendor("MyBigVendor")

	return f.addReset(func() {
		API.Hardware().DeleteVendor(f.hardwareVendor.ID)
	})
}

func (f *Fixture) setupValidationPlan() *Fixture {
	if f.validationPlan != (ValidationPlan{}) {
		return f
	}

	id, _ := uuid.FromString(validationPlanID)
	f.validationPlan = API.Validations().GetPlan(id)

	return f
}

func (f *Fixture) setupHardwareProducts() *Fixture {
	f.setupHardwareVendor()
	f.setupValidationPlan()

	mswp := newTestHardwareProduct()
	f.switchProduct = API.Hardware().Create(
		mswp.Name,
		mswp.Alias,
		f.hardwareVendor.ID,
		mswp.SKU,
		mswp.RackUnitSize,
		f.validationPlan.ID,
		mswp.Purpose,
		mswp.BiosFirmware,
		mswp.CpuType,
	)

	msvp := newTestHardwareProduct()
	f.serverProduct = API.Hardware().Create(
		msvp.Name,
		msvp.Alias,
		f.hardwareVendor.ID,
		msvp.SKU,
		msvp.RackUnitSize,
		f.validationPlan.ID,
		msvp.Purpose,
		msvp.BiosFirmware,
		msvp.CpuType,
	)

	return f.addReset(func() {
		API.Hardware().Delete(f.serverProduct.ID)
		API.Hardware().Delete(f.switchProduct.ID)
	})
}

func (f *Fixture) setupDatacenter() *Fixture {
	if f.dc != (Datacenter{}) {
		return f
	}

	mock := newTestDatacenter()
	f.dc = API.Datacenters().CreateFromStruct(mock)
	return f.addReset(func() {
		API.Datacenters().Delete(f.dc.ID)
	})
}

func (f *Fixture) setupRoom() *Fixture {
	if f.room != (Room{}) {
		return f
	}

	f.setupDatacenter()

	mockRoom := newTestRoom()
	mockRoom.DatacenterID = f.dc.ID
	f.room = API.Rooms().CreateFromStruct(mockRoom)

	return f.addReset(func() {
		API.Rooms().Delete(f.room.ID)
	})
}

func (f *Fixture) setupRackRole() *Fixture {
	if f.role != (RackRole{}) {
		return f
	}

	mock := newTestRackRole()
	f.role = API.RackRoles().CreateFromStruct(mock)

	return f.addReset(func() {
		API.RackRoles().Delete(f.role.ID)
	})
}

func (f *Fixture) setupRack() *Fixture {
	if f.rack != (Rack{}) {
		return f
	}

	f.setupRoom()
	f.setupBuild()
	f.setupRackRole()

	mock := newTestRack()
	mock.RoomID = f.room.ID
	mock.RoleID = f.role.ID
	mock.Phase = "integration"
	mock.BuildID = f.build.ID
	f.rack = API.Racks().CreateFromStruct(mock)

	return f.addReset(func() {
		API.Racks().Delete(f.rack.ID)
	})
}

func (f *Fixture) setupRackLayout() *Fixture {
	if f.rackLayout != nil {
		return f
	}

	f.setupHardwareProducts()
	f.setupRack()

	rl := RackLayoutUpdates{
		{
			RU:        1,
			ProductID: f.serverProduct.ID,
		},
		{
			RU:        1 + f.serverProduct.RackUnitSize,
			ProductID: f.switchProduct.ID,
		},
	}

	rackID := f.rack.ID
	f.rackLayout = API.Racks().CreateLayout(rackID, rl)

	return f.addReset(func() {
		for _, row := range API.Racks().Layouts(rackID) {
			API.Racks().DeleteLayoutSlot(row.ID)
		}
	})
}

func newFixture() Fixture {
	return Fixture{}
}

func newTestDatacenter() (dc Datacenter) {
	err := faker.FakeData(&dc)
	if err != nil {
		panic(err)
	}
	return
}

func newTestBuildList() (list BuildList) {
	err := faker.FakeData(&list)
	if err != nil {
		panic(err)
	}
	return
}

func newTestBuild() (build Build) {
	err := faker.FakeData(&build)
	if err != nil {
		panic(err)
	}
	return
}

func newTestUser() (user UserAndRole) {
	err := faker.FakeData(&user)
	if err != nil {
		panic(err)
	}
	return
}

func newTestUserAndRoles() (list UserAndRoles) {
	err := faker.FakeData(&list)
	if err != nil {
		panic(err)
	}
	return
}

func newTestOrganization() (org Org) {
	err := faker.FakeData(&org)
	if err != nil {
		panic(err)
	}
	return
}

func newTestOrganizationUser() (ou OrganizationUser) {
	err := faker.FakeData(&ou)
	if err != nil {
		panic(err)
	}
	return
}

func newTestOrgAndRoles() (list OrgAndRoles) {
	err := faker.FakeData(&list)
	if err != nil {
		panic(err)
	}
	return
}

func newTestOrgList() (list OrgList) {
	err := faker.FakeData(&list)
	if err != nil {
		panic(err)
	}
	return
}

func newTestDeviceList() (list DeviceList) {
	err := faker.FakeData(&list)
	if err != nil {
		panic(err)
	}
	return
}

func newTestRack() (rack Rack) {
	err := faker.FakeData(&rack)
	if err != nil {
		panic(err)
	}
	return
}

func newTestRackList() (list RackList) {
	err := faker.FakeData(&list)
	if err != nil {
		panic(err)
	}
	return
}

func newTestRoom() (room Room) {
	err := faker.FakeData(&room)
	if err != nil {
		panic(err)
	}
	return
}

func newTestRelay() (relay Relay) {
	err := faker.FakeData(&relay)
	if err != nil {
		panic(err)
	}
	return
}

func newTestRackRole() (role RackRole) {
	err := faker.FakeData(&role)
	if err != nil {
		panic(err)
	}
	return
}

func newTestHardwareProduct() (hp HardwareProduct) {
	err := faker.FakeData(&hp)
	if err != nil {
		panic(err)
	}
	return
}

func newTestHardwareVendor() (hv HardwareVendor) {
	err := faker.FakeData(&hv)
	if err != nil {
		panic(err)
	}
	return
}

func newTestRackAssignmentUpdates() (ra RackAssignmentUpdates) {
	err := faker.FakeData(&ra)
	if err != nil {
		panic(err)
	}
	return
}

func newTestDevice() (d deviceCore) {
	err := faker.FakeData(&d)
	if err != nil {
		panic(err)
	}
	return
}
