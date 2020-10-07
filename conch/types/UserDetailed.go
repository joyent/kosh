package types

import (
	"strings"

	"github.com/joyent/kosh/template"
)

func (u UserDetailed) JSON() ([]byte, error) {
	return AsJSON(u)
}

func (u UserDetailed) String() (string, error) {

	t, err := template.NewTemplate().Parse(detailedUserTemplate)
	if err != nil {
		return "", err
	}

	buf := &strings.Builder{}

	if err := t.Execute(buf, u); err != nil {
		return "", err
	}

	return buf.String(), nil
}

const detailedUserTemplate = `
ID: {{ .ID }}
Name: {{ .Name }}
Email: {{ .Email }}
System Admin: {{ if $.IsAdmin }}Yes{{ else }}No{{ end }}

Created: {{ TimeStr .Created }}
Last Login: {{ if $.LastLogin.IsZero }}Never/Unknown{{ else }}{{ TimeStr .LastLogin }}{{ end }}


Workspaces:
{{ .Workspaces }}

Organizations:
{{ .Organizations }}
`
