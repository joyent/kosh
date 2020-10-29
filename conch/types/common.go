package types

import (
	"github.com/gofrs/uuid"
)

type UUID struct{ uuid.UUID }

type MojoRelaxedPlaceholder string

type MojoStandardPlaceholder string

type Link string

type DeviceSerialNumber string

type DeviceSerialNumberEmbedded0 string

type Macaddr string

type NonEmptyString string

type EmailAddress string

func NewDeviceLinks(uris ...string) DeviceLinks {
	links := []Link{}
	for _, u := range uris {
		links = append(links, Link(u))
	}
	return DeviceLinks{links}
}
