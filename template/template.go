package template

import (
	"bytes"
	"html/template"
	"regexp"
	"time"

	"github.com/joyent/kosh/tables"
)

const dateFormat = "2006-01-02 15:04:05 -0700 MST"

// CutUUID - trims a UUID down to a short readable version
func CutUUID(id string) string {
	re := regexp.MustCompile("^(.+?)-")
	bits := re.FindStringSubmatch(id)
	if len(bits) > 0 {
		return bits[1]
	}
	return id
}

// TimeStr fformats a time value into something human readable
func TimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(dateFormat)
}

func Table(t tables.Tabulable) string {
	return tables.Render(t)
}

// NewTemplate returns a new template instance
func NewTemplate() *template.Template {
	return template.New("wat").Funcs(template.FuncMap{
		"CutUUID": func(id string) string { return CutUUID(id) },
		"TimeStr": func(t time.Time) string { return TimeStr(t) },
		"Table":   Table,
	})
}

type Templated interface {
	Template() string
}

// Render takes a struct, and a template string and returns a string
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
