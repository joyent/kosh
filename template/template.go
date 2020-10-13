package template

import (
	"html/template"
	"regexp"
	"time"
)

const DateFormat = "2006-01-02 15:04:05 -0700 MST"

func CutUUID(id string) string {
	re := regexp.MustCompile("^(.+?)-")
	bits := re.FindStringSubmatch(id)
	if len(bits) > 0 {
		return bits[1]
	}
	return id
}

func TimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(DateFormat)
}

func NewTemplate() *template.Template {
	return template.New("wat").Funcs(template.FuncMap{
		"CutUUID": func(id string) string { return CutUUID(id) },
		"TimeStr": func(t time.Time) string { return TimeStr(t) },
	})
}
