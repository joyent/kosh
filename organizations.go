package main

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Organizations struct {
	*Conch
}

func (c *Conch) Organizations() *Organizations {
	return &Organizations{c}
}

type Org struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Created     time.Time         `json:"created"`
	Admins      DetailedUsers     `json:"admins"`
	Workspaces  WorkspaceAndRoles `json:"workspaces"`
}

type OrgAndRole struct {
	Org
	Role string `json:role`
}

type OrgAndRoles []OrgAndRole

func (o OrgAndRoles) Len() int {
	return len(o)
}

func (o OrgAndRoles) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o OrgAndRoles) Less(i, j int) bool {
	return o[i].Name < o[j].Name
}

func (o OrgAndRoles) String() string {
	sort.Sort(o)
	if API.JsonOnly {
		return API.AsJSON(o)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Name",
		"Role",
		"Description",
	})

	for _, org := range o {
		table.Append([]string{
			org.Name,
			org.Role,
			org.Description,
		})
	}

	table.Render()
	return tableString.String()
}

type Orgs []Org

func (o Orgs) Len() int {
	return len(o)
}

func (o Orgs) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o Orgs) Less(i, j int) bool {
	return o[i].Name < o[j].Name
}

func (o Orgs) String() string {
	sort.Sort(o)
	if API.JsonOnly {
		return API.AsJSON(o)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Name",
		"Role",
		"Description",
	})

	for _, org := range o {
		table.Append([]string{
			org.Name,
			org.Description,
		})
	}

	table.Render()
	return tableString.String()
}

func (o Org) String() string {
	if API.JsonOnly {
		return API.AsJSON(o)
	}

	t, err := NewTemplate().Parse(organizationTemplate)
	if err != nil {
		panic(err)
	}

	buf := &strings.Builder{}

	if err := t.Execute(buf, o); err != nil {
		panic(err)
	}

	return buf.String()
}

func (o *Organizations) GetAll() Orgs {
	var list Orgs
	res := o.Do(o.Sling().Get("/organization"))
	if ok := res.Parse(&list); !ok {
		panic(res)
	}
	return list
}

func (o *Organizations) Get(ID uuid.UUID) Org {
	var org Org
	uri := fmt.Sprintf("/organization/%s", url.PathEscape(ID.String()))
	res := o.Do(o.Sling().Get(uri))
	if ok := res.Parse(&org); !ok {
		panic(res)
	}
	return org
}

func (o *Organizations) GetByName(name string) Org {
	var org Org
	uri := fmt.Sprintf("/organization/%s", url.PathEscape(name))
	res := o.Do(o.Sling().Get(uri))
	if ok := res.Parse(&org); !ok {
		panic(res)
	}
	return org
}

func (o *Organizations) Create(name, description string, admins []map[string]string) (org Org) {
	payload := make(map[string]interface{})
	payload["name"] = name
	payload["admins"] = admins
	if description != "" {
		payload["description"] = description
	}

	res := o.Do(o.Sling().New().Post("/organization").
		Set("Content-Type", "application/json").
		BodyJSON(payload),
	)

	if ok := res.Parse(&org); !ok {
		panic(fmt.Sprintf("%v", res))
	}

	return
}

func (o *Organizations) Delete(ID uuid.UUID) {
	uri := fmt.Sprintf("/organization/%s", url.PathEscape(ID.String()))
	_ = o.Do(o.Sling().Delete(uri))
}

type OrganizationUser struct {
	Email string    `json:"email"`
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}

func (o *Organizations) GetUsers(ID uuid.UUID) (users []OrganizationUser) {
	uri := fmt.Sprintf("/organization/%s/user", url.PathEscape(ID.String()))
	res := o.Do(o.Sling().Get(uri))
	if ok := res.Parse(&users); !ok {
		panic(res)
	}
	return
}

func (o *Organizations) AddUser(orgID uuid.UUID, email, role string, sendEmail bool) {
	uri := fmt.Sprintf("/organization/%s/user", url.PathEscape(orgID.String()))

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

	fmt.Printf("%v\n", payload)
	_ = o.Do(
		o.Sling().Post(uri).
			Set("Content-Type", "application/json").
			QueryStruct(q).
			BodyJSON(payload),
	)
}

// userID is a string because it may be a UUID or an Email, the API accepts both
func (o *Organizations) RemoveUser(orgID uuid.UUID, userID string, sendEmail bool) bool {
	uri := fmt.Sprintf(
		"/organization/%s/user/%s",
		url.PathEscape(orgID.String()),
		url.PathEscape(userID),
	)

	send := 0
	if sendEmail {
		send = 1
	}
	q := struct {
		SendEmail int `url:"send_mail"`
	}{send}

	res := o.Do(o.Sling().Delete(uri).QueryStruct(q))

	return res.StatusCode() == 204
}

func init() {

	App.Command("organizations orgs", "Work with organizations", func(cmd *cli.Cmd) {
		cmd.Command("get", "Get a list of all organizations", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Organizations().GetAll())
			}
		})

		cmd.Command("create", "Create a new subworkspace", func(cmd *cli.Cmd) {
			nameArg := cmd.StringArg("NAME", "", "Name of the new organization")

			descOpt := cmd.StringOpt("description", "", "A description of the organization")
			adminEmailArg := cmd.StringOpt(
				"admin",
				"",
				"Email address for the (initial) admin user for the organization. This does *not* create the user.",
			)

			cmd.Spec = "NAME [OPTIONS]"
			cmd.Action = func() {
				API.Organizations().Create(
					*nameArg,
					*descOpt,
					[]map[string]string{{"email": *adminEmailArg}},
				)
			}
		})

	})

	App.Command("organization org", "Work with a specific organization", func(cmd *cli.Cmd) {
		var o Org
		organizationNameArg := cmd.StringArg("NAME", "", "Name or ID of the Organization")

		cmd.Spec = "NAME"
		cmd.Before = func() {
			o = API.Organizations().GetByName(*organizationNameArg)
			// TODO(sungo): should we verify that the organization exists?
		}

		cmd.Command("get", "Get information about a single organization by its name", func(cmd *cli.Cmd) {

			cmd.Action = func() {
				fmt.Println(o)
			}
		})

		cmd.Command("delete rm", "Remove a specific organization", func(cmd *cli.Cmd) {

			cmd.Action = func() {
				API.Organizations().Delete(o.ID)
			}
		})

		cmd.Command("users", "Manage users in a specific organization", func(cmd *cli.Cmd) {

			cmd.Command("get ls", "Get a list of users in an organization", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Organizations().GetUsers(o.ID))
				}
			})

			cmd.Command("add", "Add a user to an organization", func(cmd *cli.Cmd) {
				userEmailArg := cmd.StringArg(
					"EMAIL",
					"",
					"The email of the user to add to the workspace. Does *not* create the user",
				)

				roleOpt := cmd.StringOpt(
					"role",
					"ro",
					"The role for the user. One of: "+prettyWorkspaceRoleList(),
				)

				sendEmailOpt := cmd.BoolOpt(
					"send-email",
					true,
					"Send email to the target user, notifying them of the change",
				)

				cmd.Spec = "EMAIL [OPTIONS]"
				cmd.Action = func() {
					if !okWorkspaceRole(*roleOpt) {
						panic(fmt.Errorf(
							"'role' value must be one of: %s",
							prettyWorkspaceRoleList(),
						))
					}
					fmt.Printf("%v\n", *userEmailArg)
					API.Organizations().AddUser(
						o.ID,
						*userEmailArg,
						*roleOpt,
						*sendEmailOpt,
					)
					fmt.Println(API.Organizations().GetUsers(o.ID))
				}

			})

			cmd.Command("remove rm", "remove a user from an organization", func(cmd *cli.Cmd) {
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
					API.Organizations().RemoveUser(
						o.ID,
						*userEmailArg,
						*sendEmailOpt,
					)
					fmt.Println(API.Organizations().GetUsers(o.ID))
				}
			})
		})
	})
}
