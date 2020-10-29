package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joyent/kosh/conch"
	"github.com/joyent/kosh/logger"
	"github.com/joyent/kosh/tables"
	"github.com/joyent/kosh/template"
)

// DefaultConfig is the default configuration struct
type Config struct {
	Version string
	GitRev  string

	ConchURL         string
	ConchToken       string
	ConchEnvironment string

	OutputJSON bool

	logger.Logger
}

// NewConfig takes a Version and a GitRev and returns a Config object
func NewConfig(Version, GitRev string) *Config {
	return &Config{
		Version: Version,
		GitRev:  GitRev,
		Logger:  logger.New(),
	}
}

const configTemplate = `
---
# Config

* Version: {{ .Version }}
* GitRev: {{ .GitRev }}

* ConchURL: {{ .ConchURL }}
* ConchToken: {{ .ConchToken }}
* ConchEnvironment: {{ .ConchEnvironment }}

* OutputJSON: {{ .OutputJSON }}
---
`

// String returns a string implementation of the config object
func (c *Config) String() string {
	t, err := template.NewTemplate().Parse(configTemplate)
	if err != nil {
		log.Fatal(err)
	}

	buf := &strings.Builder{}
	if e := t.Execute(buf, c); e != nil {
		log.Fatal(e)
	}
	return buf.String()
}

// ConchClient returns a configured client for the Conch API
func (c *Config) ConchClient() *conch.Client {
	userAgent := fmt.Sprintf("kosh %s", c.GitRev)
	return conch.New(
		conch.API(c.ConchURL),
		conch.AuthToken(c.ConchToken),
		conch.UserAgent(userAgent),
		conch.Logger(c.Logger),
	)
}

// Renderer returns a function that will render to STDOUT
func (c *Config) Renderer() func(interface{}) {
	return c.RenderTo(os.Stdout)
}

func renderJSON(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		fatal(e)
	}
	return string(b)
}

// RenderTo returns a function tha renders to a given io.Writer based on the
// configuraton and datatype
func (c *Config) RenderTo(w io.Writer) func(interface{}) {
	return func(i interface{}) {
		if c.OutputJSON {
			fmt.Fprintln(w, renderJSON(i))
		}
		switch t := i.(type) {
		case template.Templated:
			s, e := template.Render(t)
			if e != nil {
				fatal(e)
			}
			fmt.Fprintln(w, s)
		case tables.Tabulable:
			fmt.Fprintln(w, tables.Render(t))
		default:
			fmt.Fprintln(w, renderJSON(t))
		}
	}
}
