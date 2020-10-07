package types

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

type UUID struct{ uuid.UUID }

// see https://metacpan.org/pod/Mojolicious::Guides::Routing#Relaxed-placeholders
type MojoRelaxedPlaceholder string

// see https://metacpan.org/pod/Mojolicious::Guides::Routing#Standard-placeholders
type MojoStandardPlaceholder string

func AsJSON(d interface{}) ([]byte, error) {
	return json.MarshalIndent(d, "", "    ")
}

type Link string

type DeviceSerialNumber string

type DeviceSerialNumberEmbedded0 string

type Macaddr string

type NonEmptyString string

type EmailAddress string
