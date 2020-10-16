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

type Config interface {
	GetVersion() string

	SetURL(string)

	SetToken(string)

	SetLogger(logger.Logger)

	GetOutputJSON() bool
	SetOutputJSON(bool)

	ConchClient() *conch.Client
	Renderer() func(interface{})

	conch.Config
}

// Config struct
type DefaultConfig struct {
	Version string
	GitRev  string

	ConchURL         string
	ConchToken       string
	ConchEnvironment string

	OutputJSON bool

	logger.Logger
}

func (c *DefaultConfig) GetVersion() string       { return c.Version }
func (c *DefaultConfig) GetURL() string           { return c.ConchURL }
func (c *DefaultConfig) GetToken() string         { return c.ConchToken }
func (c *DefaultConfig) GetLogger() logger.Logger { return c.Logger }

func (c *DefaultConfig) SetURL(URL string)         { c.ConchURL = URL }
func (c *DefaultConfig) SetToken(token string)     { c.ConchToken = token }
func (c *DefaultConfig) GetOutputJSON() bool       { return c.OutputJSON }
func (c *DefaultConfig) SetOutputJSON(p bool)      { c.OutputJSON = p }
func (c *DefaultConfig) SetLogger(l logger.Logger) { c.Logger = l }

// NewConfig takes a Version and a GitRev and returns a Config object
func NewConfig(Version, GitRev string) *DefaultConfig {
	return &DefaultConfig{
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

func (c *DefaultConfig) String() string {
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
func (c *DefaultConfig) ConchClient() *conch.Client {
	return conch.New(c)
}

func (c *DefaultConfig) Renderer() func(interface{}) {
	return c.RenderTo(os.Stdout)
}

func renderJSON(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		fatal(e)
	}
	return string(b)
}

func (c *DefaultConfig) RenderTo(w io.Writer) func(interface{}) {
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
