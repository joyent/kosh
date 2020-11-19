/*
Package template provies utilities for rendering the output to markdown
compatible output
*/
package template

import (
	"bytes"
	"html/template"
	"regexp"
	"time"

	"github.com/joyent/kosh/tables"
)

const dateFormat = "2006-01-02 15:04:05 -0700 MST"

// YesOrNo transforms a bool into "yes" or "no" depending on it's truth
func YesOrNo(p bool) string {
	if p {
		return "Yes"
	}
	return "No"
}

// CutUUID - trims a UUID down to a short readable version
func CutUUID(id string) string {
	re := regexp.MustCompile("^(.+?)-")
	bits := re.FindStringSubmatch(id)
	if len(bits) > 0 {
		return bits[1]
	}
	return id
}

// TimeStr formats a time value into something human readable
func TimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(dateFormat)
}

// Table formats Tabulable data using the tables package
func Table(t tables.Tabulable) string {
	return tables.Render(t)
}

// NewTemplate returns a new template instance
func NewTemplate() *template.Template {
	return template.New("wat").Funcs(template.FuncMap{
		"CutUUID": CutUUID, // func(id string) string { return CutUUID(id) },
		"TimeStr": func(t time.Time) string { return TimeStr(t) },
		"Table":   Table,
	})
}

// Templated tracks what the template is for a given data structure
type Templated interface {
	Template() string
}

// Render takes a Templated datapiecea and returns a markdown compatible string
func Render(data Templated) (string, error) {
	template := data.Template()
	t, err := NewTemplate().Parse(template)
	if err != nil {
		return "", err // TODO get logging in here
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err // TODO get logging in here
	}

	return buf.String(), nil
}
