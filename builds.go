package main

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/template"
	"github.com/olekukonko/tablewriter"
)

type Builds struct {
	*Conch
}

type Build struct {
	ID            uuid.UUID    `json:"id" faker:"uuid"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Admins        UserAndRoles `json:"admins"`
	Created       time.Time    `json:"created" faker:"-"`
	Started       time.Time    `json:"started" faker:"-"`
	Completed     time.Time    `json:"completed" faker:"-"`
	CompletedUser UserAndRole  `json:"completed_user" faker:"-"`
}

func (b Build) String() string {
	if API.JsonOnly {
		return API.AsJSON(b)
	}

	t, err := template.NewTemplate().Parse(buildTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, b); err != nil {
		panic(err)
	}

	return buf.String()
}

type BuildList []Build

func (bl BuildList) Len() int {
	return len(bl)
}

func (bl BuildList) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}

func (bl BuildList) Less(i, j int) bool {
	return bl[i].Name < bl[j].Name
}

func (bl BuildList) String() string {
	sort.Sort(bl)
	if API.JsonOnly {
		return API.AsJSON(bl)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Name",
		"Description",
		"Started",
		"Completed",
		"Completed By",
	})

	for _, b := range bl {
		table.Append([]string{
			b.Name,
			b.Description,
			b.Started.String(),
			b.Completed.String(),
			b.CompletedUser.Email,
		})
	}

	table.Render()
	return tableString.String()
}

func (c *Conch) Builds() *Builds {
	return &Builds{c}
}

var BuildRoleList = []string{"admin", "rw", "ro"}

func prettyBuildRoleList() string {
	return strings.Join(BuildRoleList, ", ")
}

func okBuildRole(role string) bool {
	for _, b := range BuildRoleList {
		if role == b {
			return true
		}
	}
	return false
}


func (b *Builds) GetAll() (list BuildList) {
	res := b.Do(b.Sling().Get("/build"))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}
	return
}

func (b *Builds) Get(ID uuid.UUID) (build Build) {
	uri := fmt.Sprintf("/build/%s", url.PathEscape(ID.String()))
	res := b.Do(b.Sling().Get(uri))
	if ok := res.Parse(&build); !ok {
		panic(res)
	}
	return
}

func (b *Builds) GetByName(name string) (build Build) {
	uri := fmt.Sprintf("/build/%s", url.PathEscape(name))
	res := b.Do(b.Sling().Get(uri))
	if ok := res.Parse(&build); !ok {
		panic(res)
	}
	return
}

func (b *Builds) Create(name, description string, admins []map[string]string) (build Build) {
	payload := make(map[string]interface{})
	payload["name"] = name
	payload["admins"] = admins
	if description != "" {
		payload["description"] = description
	}

	res := b.Do(b.Sling().New().Post("/build").
		Set("Content-Type", "application/json").
		BodyJSON(payload),
	)

	if ok := res.Parse(&build); !ok {
		panic(res)
	}

	return
}

func (b *Builds) GetUsers(ID uuid.UUID) (list UserAndRoles) {
	res := b.Do(b.Sling().Get(fmt.Sprintf("/build/%s/user", ID.String())))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}
	return
}

func (b *Builds) AddUser(ID uuid.UUID, email, role string, sendEmail bool) {
	uri := fmt.Sprintf("/build/%s/user", url.PathEscape(ID.String()))

	payload := make(map[string]string)
	payload["email"] = email
	payload["role"] = role

	send := 0
	if sendEmail {
		send = 1
	}
	q := struct {
		SendEmail int `url:"send_mail"`
	}{send}

	_ = b.Do(
		b.Sling().Post(uri).
			Set("Content-Type", "application/json").
			QueryStruct(q).
			BodyJSON(payload),
	)
}

// userID is a string because it may be a UUID or an Email, the API accepts both
func (b *Builds) RemoveUser(ID uuid.UUID, userID string, sendEmail bool) bool {
	uri := fmt.Sprintf(
		"/build/%s/user/%s",
		url.PathEscape(ID.String()),
		url.PathEscape(userID),
	)

	send := 0
	if sendEmail {
		send = 1
	}
	q := struct {
		SendEmail int `url:"send_mail"`
	}{send}

	res := b.Do(b.Sling().Delete(uri).QueryStruct(q))

	return res.StatusCode() == 204
}

func (b *Builds) GetOrgs(ID uuid.UUID) (list OrgAndRoles) {
	res := b.Do(b.Sling().Get(fmt.Sprintf("/build/%s/organization", ID.String())))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}
	return
}

func (b *Builds) AddOrg(ID uuid.UUID, orgID, role string, sendEmail bool) {
	uri := fmt.Sprintf("/build/%s/organization", url.PathEscape(ID.String()))

	payload := make(map[string]string)
	payload["organization_id"] = orgID
	payload["role"] = role

	send := 0
	if sendEmail {
		send = 1
	}
	q := struct {
		SendEmail int `url:"send_mail"`
	}{send}

	_ = b.Do(
		b.Sling().Post(uri).
			Set("Content-Type", "application/json").
			QueryStruct(q).
			BodyJSON(payload),
	)
}

// userID is a string because it may be a UUID or an Email, the API accepts both
func (b *Builds) RemoveOrg(ID uuid.UUID, orgID string, sendEmail bool) bool {
	uri := fmt.Sprintf(
		"/build/%s/organization/%s",
		url.PathEscape(ID.String()),
		url.PathEscape(orgID),
	)

	send := 0
	if sendEmail {
		send = 1
	}
	q := struct {
		SendEmail int `url:"send_mail"`
	}{send}

	res := b.Do(b.Sling().Delete(uri).QueryStruct(q))

	return res.StatusCode() == 204
}

func (b *Builds) GetDevices(ID uuid.UUID) (list DeviceList) {
	res := b.Do(b.Sling().Get(fmt.Sprintf("/build/%s/device", ID.String())))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}
	return
}

func (b *Builds) CreateDevice(ID uuid.UUID, deviceID, sku string) {
	type BDC struct {
		Serial string `json:"serial_number"`
		SKU    string `json:"sku"`
	}

	type BDCList []BDC

	list := BDCList{{deviceID, sku}}

	uri := fmt.Sprintf("/build/%s/device", url.PathEscape(ID.String()))

	_ = b.Do(
		b.Sling().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(list),
	)
}

func (b *Builds) AddDevice(ID uuid.UUID, deviceID string) {
	uri := fmt.Sprintf("/build/%s/device/%s", url.PathEscape(ID.String()), url.PathEscape(deviceID))

	_ = b.Do(b.Sling().Post(uri))
}

func (b *Builds) RemoveDevice(ID uuid.UUID, deviceID string) bool {
	uri := fmt.Sprintf(
		"/build/%s/device/%s",
		url.PathEscape(ID.String()),
		url.PathEscape(deviceID),
	)

	res := b.Do(b.Sling().Delete(uri))

	return res.StatusCode() == 204
}

func (b *Builds) GetRacks(ID uuid.UUID) (list RackList) {
	res := b.Do(b.Sling().Get(fmt.Sprintf("/build/%s/rack", ID.String())))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}
	return
}

func (b *Builds) AddRack(ID uuid.UUID, rackID string) {
	uri := fmt.Sprintf("/build/%s/rack/%s", url.PathEscape(ID.String()), url.PathEscape(rackID))

	_ = b.Do(b.Sling().Post(uri))
}

func (b *Builds) RemoveRack(ID uuid.UUID, rackID string) bool {
	uri := fmt.Sprintf(
		"/build/%s/rack/%s",
		url.PathEscape(ID.String()),
		url.PathEscape(rackID),
	)

	res := b.Do(b.Sling().Delete(uri))

	return res.StatusCode() == 204
}

func init() {

	App.Command("builds", "Work with builds", func(cmd *cli.Cmd) {
		cmd.Command("get ls", "Get a list of all builds", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Builds().GetAll())
			}
		})

		cmd.Command("create", "Create a new build", func(cmd *cli.Cmd) {
			nameArg := cmd.StringArg("NAME", "", "Name of the new build")

			descOpt := cmd.StringOpt("description", "", "A description of the build")
			adminEmailArg := cmd.StringOpt(
				"admin",
				"",
				"Email address for the (initial) admin user for the build. This does *not* create the user.",
			)

			cmd.Spec = "NAME [OPTIONS]"
			cmd.Action = func() {
				API.Builds().Create(
					*nameArg,
					*descOpt,
					[]map[string]string{{"email": *adminEmailArg}},
				)
			}
		})

	})

	App.Command("build", "Work with a specific build", func(cmd *cli.Cmd) {
		var b Build
		buildNameArg := cmd.StringArg("NAME", "", "Name or ID of the build")

		cmd.Spec = "NAME"
		cmd.Before = func() {
			b = API.Builds().GetByName(*buildNameArg)
			// TODO(sungo): should we verify that the build exists?
		}

		cmd.Command("get", "Get information about a single build by its name", func(cmd *cli.Cmd) {

			cmd.Action = func() {
				fmt.Println(b)
			}
		})

		cmd.Command("users", "Manage users in a specific build", func(cmd *cli.Cmd) {

			cmd.Command("get ls", "Get a list of users in an build", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Builds().GetUsers(b.ID))
				}
			})

			cmd.Command("add", "Add a user to an build", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email of the user to add to the build. Does *not* create the user",
				)

				roleOpt := cmd.StringOpt(
					"role",
					"ro",
					"The role for the user. One of: "+prettyBuildRoleList(),
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target user, notifying them of the change",
				)

				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					if !okBuildRole(*roleOpt) {
						panic(fmt.Errorf(
							"'role' value must be one of: %s",
							prettyBuildRoleList(),
						))
					}
					API.Builds().AddUser(
						b.ID,
						*userEmailArg,
						*roleOpt,
						*sendEmailOpt,
					)
					fmt.Println(API.Builds().GetUsers(b.ID))
				}

			})

			cmd.Command("remove rm", "remove a user from an build", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email or ID of the user to modify",
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target user, notifying them of the change",
				)
				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					API.Builds().RemoveUser(
						b.ID,
						*userEmailArg,
						*sendEmailOpt,
					)
					fmt.Println(API.Builds().GetUsers(b.ID))
				}
			})
		})

		cmd.Command("organizations orgs", "Manage organizations in a specific build", func(cmd *cli.Cmd) {

			cmd.Command("get ls", "Get a list of organizations in an build", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Builds().GetOrgs(b.ID))
				}
			})

			cmd.Command("add", "Add a organization to an build", func(cmd *cli.Cmd) {
				orgNameArg := cmd.StringArg(
					"NAME",
					"",
					"The name of the organization to add to the build. Does *not* create the organization",
				)

				roleOpt := cmd.StringOpt(
					"role",
					"ro",
					"The role for the organization. One of: "+prettyBuildRoleList(),
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the organization admins, notifying them of the change",
				)

				cmd.Spec = "NAME [OPTIONS]"
				cmd.Action = func() {
					if !okBuildRole(*roleOpt) {
						panic(fmt.Errorf(
							"'role' value must be one of: %s",
							prettyBuildRoleList(),
						))
					}
					API.Builds().AddOrg(
						b.ID,
						*orgNameArg,
						*roleOpt,
						*sendEmailOpt,
					)
					fmt.Println(API.Builds().GetOrgs(b.ID))
				}

			})

			cmd.Command("remove rm", "remove an organization from a build", func(cmd *cli.Cmd) {
				orgNameArg := cmd.StringArg(
					"NAME",
					"",
					"The name or ID of the organization to modify",
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target organization admins, notifying them of the change",
				)
				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					API.Builds().RemoveOrg(
						b.ID,
						*orgNameArg,
						*sendEmailOpt,
					)
					fmt.Println(API.Builds().GetOrgs(b.ID))
				}
			})
		})

		cmd.Command("devices ds", "Manage devices in a specific build", func(cmd *cli.Cmd) {

			cmd.Command("get ls", "Get a list of devices in an build", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Builds().GetDevices(b.ID))
				}
			})

			cmd.Command("add", "Add a device to an build", func(cmd *cli.Cmd) {
				deviceIDArg := cmd.StringArg(
					"ID",
					"",
					"The ID or serial number of the device to add to the build. Does *not* create the device",
				)

				cmd.Spec = "ID [OPTIONS]"
				cmd.Action = func() {
					API.Builds().AddDevice(
						b.ID,
						*deviceIDArg,
					)
					fmt.Println(API.Builds().GetDevices(b.ID))
				}

			})

			cmd.Command("remove rm", "remove a device from a build", func(cmd *cli.Cmd) {
				deviceIDArg := cmd.StringArg(
					"ID",
					"",
					"The ID or serial number of the device to add to the build. Does *not* create the device",
				)

				cmd.Spec = "ID [OPTIONS]"
				cmd.Action = func() {
					API.Builds().RemoveDevice(
						b.ID,
						*deviceIDArg,
					)
					fmt.Println(API.Builds().GetDevices(b.ID))
				}
			})
		})

		cmd.Command("racks", "Manage racks in a specific build", func(cmd *cli.Cmd) {

			cmd.Command("get ls", "Get a list of racks in an build", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Builds().GetRacks(b.ID))
				}
			})

			cmd.Command("add", "Add a rack to an build", func(cmd *cli.Cmd) {
				rackIDArg := cmd.StringArg(
					"ID",
					"",
					"The ID of the rack to add to the build. Does *not* create the rack",
				)

				cmd.Spec = "ID [OPTIONS]"
				cmd.Action = func() {
					API.Builds().AddRack(
						b.ID,
						*rackIDArg,
					)
					fmt.Println(API.Builds().GetRacks(b.ID))
				}

			})

			cmd.Command("remove rm", "remove a rack from a build", func(cmd *cli.Cmd) {
				rackIDArg := cmd.StringArg(
					"ID",
					"",
					"The ID of the rack to add to the build. Does *not* create the device",
				)

				cmd.Spec = "ID [OPTIONS]"
				cmd.Action = func() {
					API.Builds().RemoveRack(
						b.ID,
						*rackIDArg,
					)
					fmt.Println(API.Builds().GetRacks(b.ID))
				}
			})
		})
	})
}
