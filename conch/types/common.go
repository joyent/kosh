package types

import (
	"github.com/gofrs/uuid"
)

// UUID is a Universally Unique ID
type UUID struct{ uuid.UUID }

// MojoRelaxedPlaceholder is a string
type MojoRelaxedPlaceholder string

// MojoStandardPlaceholder is a string
type MojoStandardPlaceholder string

// Link is a string
type Link string

// DeviceSerialNumber is a string
type DeviceSerialNumber string

// DeviceSerials is a slice of DeviceSerialNumbers
type DeviceSerials []DeviceSerialNumber

// DeviceSerialNumberEmbedded0 is a string
type DeviceSerialNumberEmbedded0 string

// Macaddr is a string
type Macaddr string

// NonEmptyString is a string
type NonEmptyString string

// EmailAddress is a string
type EmailAddress string

// NewDeviceLinks takes a list of strings and returns them as a single
// DeviceLinks slice
func NewDeviceLinks(uris ...string) DeviceLinks {
	links := []Link{}
	for _, u := range uris {
		links = append(links, Link(u))
	}
	return DeviceLinks{links}
}

// IntOrStringyInt is an integer that may be presented as a json string
type IntOrStringyInt interface{}

// DiskSizeItem is an int
type DiskSizeItem int

// DeviceHealth is a string
// corresponds to device_health_enum in the database
type DeviceHealth string

// DevicePhase corresponds to device_phase_enum in the database (also used for racks)
type DevicePhase string

// NonNegativeInteger is an int
type NonNegativeInteger int

// PositiveInteger is an int
type PositiveInteger int

// Ipaddr is a string
type Ipaddr string
