package cli

import (
	"errors"
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/conch/types"
)

func racksCmd(cmd *cli.Cmd) {
	var conch *conch.Client

	cmd.Before = func() {
		config.requireAuth()
		config.requireSysAdmin()
		conch = config.ConchClient()
	}

	cmd.Command("create", "Create a new rack", func(cmd *cli.Cmd) {
		var (
			nameOpt      = cmd.StringOpt("name", "", "Name of the rack")
			roomAliasOpt = cmd.StringOpt("room", "", "Alias of the datacenter room")
			roleNameOpt  = cmd.StringOpt("role", "", "Name of the role")
			buildNameOpt = cmd.StringOpt("build", "", "Build for the rack")
			phaseOpt     = cmd.StringOpt("phase", "", "Optional phase for the rack")
		)

		cmd.Spec = "--name --room --role [OPTIONS]"
		cmd.Action = func() {
			var (
				roomID  types.UUID
				roleID  types.UUID
				buildID types.UUID
			)

			// The user can be very silly and supply something like
			// `--name ""` which will pass the cli lib's requirement
			// check but is still crap
			if *nameOpt == "" {
				fatalIf(errors.New("--name is required"))
			}

			if *roomAliasOpt == "" {
				fatalIf(errors.New("--room is required"))
			} else {
				room, e := conch.GetRoomByAlias(*roomAliasOpt)
				fatalIf(e)

				if (room == types.DatacenterRoomDetailed{}) {
					fatalIf(errors.New("could not find room"))
				}
				roomID = room.ID
			}

			if *roleNameOpt == "" {
				fatalIf(errors.New("--role is required"))
			} else {
				role, e := conch.GetRackRoleByName(*roleNameOpt)
				if e != nil {
					fatalIf(e)
				}
				if (role == types.RackRole{}) {
					fatalIf(errors.New("could not find rack role"))
				}
				roleID = role.ID
			}

			if *buildNameOpt == "" {
				fatalIf(errors.New("--build is required"))
			} else {
				build, e := conch.GetBuildByName(*buildNameOpt)
				if e != nil {
					fatalIf(e)
				}
				buildID = build.ID
			}
			conch.CreateRack(types.RackCreate{
				Name:             types.MojoRelaxedPlaceholder(*nameOpt),
				DatacenterRoomID: roomID,
				RackRoleID:       roleID,
				Phase:            types.DevicePhase(*phaseOpt),
				BuildID:          buildID,
			})
		}
	})
	cmd.Command("import", "Import a new rack from json", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file that defines the layout. '-' indicates STDIN")

			input, e := getInputReader(*filePathArg)
			fatalIf(e)

			rackCreate := conch.ReadRackCreate(input)
			conch.CreateRack(rackCreate)
		}
	})
}

func rackCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{}, error)

	var rack types.Rack

	idArg := cmd.StringArg(
		"UUID",
		"",
		"The UUID of the rack. Short UUIDs are *not* accepted, unless you are a Conch sysadmin",
	)

	cmd.Spec = "UUID"

	cmd.Before = func() {
		config.requireAuth()
		conch = config.ConchClient()
		display = config.Renderer()

		var e error
		rack, e = conch.GetRackByName(*idArg)
		if e != nil {
			fatalIf(e)
		}
		if (rack == types.Rack{}) {
			fatalIf(errors.New("could not find the rack"))
		}
	}

	cmd.Command("get", "Get a single rack", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(rack, nil) }
	})

	cmd.Command("update", "Update information about a single rack", func(cmd *cli.Cmd) {
		var (
			nameOpt      = cmd.StringOpt("name", "", "Name of the rack")
			roomAliasOpt = cmd.StringOpt("room", "", "Alias of the datacenter room")
			roleNameOpt  = cmd.StringOpt("role", "", "Name of the role")
			phaseOpt     = cmd.StringOpt("phase", "", "Phase for the rack")

			serialNumberOpt = cmd.StringOpt("serial-number", "", "Serial number of the rack")
			clearSerialOpt  = cmd.BoolOpt("clear-serial-number", false, "Delete the serial number. Overrides --serial-number")

			assetTagOpt      = cmd.StringOpt("asset-tag", "", "Asset Tag of the rack")
			clearAssetTagOpt = cmd.BoolOpt("clear-asset-tag", false, "Delete the asset tag. Overrides --asset-tag")
		)

		cmd.Action = func() {
			var (
				roomID   types.UUID
				roleID   types.UUID
				serial   *string
				assetTag *string
			)

			if *roomAliasOpt != "" {
				room, e := conch.GetRoomByAlias(*roomAliasOpt)
				if e != nil {
					fatalIf(e)
				}
				if (room == types.DatacenterRoomDetailed{}) {
					fatalIf(errors.New("could not find room"))
				}
				roomID = room.ID
			}
			if *roleNameOpt != "" {
				role, e := conch.GetRackRoleByName(*roomAliasOpt)
				if e != nil {
					fatalIf(e)
				}
				if (role == types.RackRole{}) {
					fatalIf(errors.New("could not find rack role"))
				}
				roleID = role.ID
			}

			empty := ""

			if *clearSerialOpt {
				serial = nil
			} else if *serialNumberOpt != "" {
				serial = serialNumberOpt
			} else {
				serial = &empty
			}

			if *clearAssetTagOpt {
				assetTag = nil
			} else if *assetTagOpt != "" {
				assetTag = assetTagOpt
			} else {
				assetTag = &empty
			}

			conch.UpdateRack(rack.ID, types.RackUpdate{
				Name:             types.MojoRelaxedPlaceholder(*nameOpt),
				DatacenterRoomID: roomID,
				RackRoleID:       roleID,
				Phase:            types.DevicePhase(*phaseOpt),
				SerialNumber:     serial,
				AssetTag:         assetTag,
			})
		}
	})

	cmd.Command("delete rm", "Delete a rack", func(cmd *cli.Cmd) {
		cmd.Before = func() {
			config.requireAuth()
			config.requireSysAdmin()
		}

		cmd.Action = func() {
			conch.DeleteRack(rack.ID)
			fmt.Println("OK")
		}
	})

	cmd.Command("layout", "The layout of the rack", func(cmd *cli.Cmd) {
		cmd.Command("get", "Get the layout of a rack", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				display(conch.GetRackLayout(rack.ID))
			}
		})

		cmd.Command("export", "Export the layout of the rack as JSON", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				l, e := conch.GetRackLayout(rack.ID)
				if e != nil {
					fatalIf(e)
				}
				fmt.Println(renderJSON(l))
			}
		})

		cmd.Command("import", "Import the layout of this rack (using the same format as 'export')", func(cmd *cli.Cmd) {
			filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file that defines the layout. '-' indicates STDIN")
			overwriteOpt := cmd.BoolOpt("overwrite", false, "If the rack has an existing layout, *overwrite* it. This is a destructive action")

			cmd.Action = func() {
				layout, e := conch.GetRackLayout(rack.ID)
				if e != nil {
					fatalIf(e)
				}
				if len(layout) > 0 {
					if !*overwriteOpt {
						fatalIf(errors.New("rack already has a layout. use --overwrite to force"))
					}
				}

				input, e := getInputReader(*filePathArg)
				fatalIf(e)

				update := conch.ReadRackLayoutUpdate(input)
				conch.UpdateRackLayout(rack.ID, update)
				fmt.Println("OK")
			}
		})
	})

	cmd.Command("assign", "Assign devices to rack slots, using the `--json` output from 'assignments'", func(cmd *cli.Cmd) {
		filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file to use as the data source. '-' indicates STDIN")
		cmd.Action = func() {
			input, err := getInputReader(*filePathArg)
			if err != nil {
				fatalIf(err)
			}
			update := conch.ReadRackAssignmentUpdate(input)
			conch.UpdateRackAssignments(rack.ID, update)
		}
	})

	cmd.Command("assignments", "The devices assigned to the rack", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			display(conch.GetRackAssignments(rack.ID))
		}
	})
}
