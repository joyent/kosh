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

// Config is the interface for the configuration object for the full app
type Config interface {
	GetVersion() string
	GetGitRev() string

	SetURL(string)

	SetToken(string)

	SetLogger(logger.Logger)

	GetOutputJSON() bool
	SetOutputJSON(bool)

	ConchClient() *conch.Client
	Renderer() func(interface{})

	conch.Config
}

// DefaultConfig is the default configuration struct
type DefaultConfig struct {
	Version string
	GitRev  string

	ConchURL         string
	ConchToken       string
	ConchEnvironment string

	OutputJSON bool

	logger.Logger
}

// GetVersion returns the current app version
func (c *DefaultConfig) GetVersion() string { return c.Version }

// GetGitRev returns the current app version
func (c *DefaultConfig) GetGitRev() string { return c.GitRev }

// GetURL returns the current conch API url
func (c *DefaultConfig) GetURL() string { return c.ConchURL }

// GetToken returns the current conch API token
func (c *DefaultConfig) GetToken() string { return c.ConchToken }

// GetLogger returns the current logger instance
func (c *DefaultConfig) GetLogger() logger.Logger { return c.Logger }

// SetURL updates the current conch API url
func (c *DefaultConfig) SetURL(URL string) { c.ConchURL = URL }

// SetToken updates the current conch API token being used
func (c *DefaultConfig) SetToken(token string) { c.ConchToken = token }

// GetOutputJSON returns whether JSON is used for output or not
func (c *DefaultConfig) GetOutputJSON() bool { return c.OutputJSON }

// SetOutputJSON sets whether to use JSOn for output or not
func (c *DefaultConfig) SetOutputJSON(p bool) { c.OutputJSON = p }

// SetLogger sets the configured logger instance
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

// String returns a string implementation of the config object
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
	userAgent := fmt.Sprintf("kosh %s", c.GitRev)
	return conch.New(c).UserAgent(userAgent)
}

// Renderer returns a function that will render to STDOUT
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

// RenderTo returns a function tha renders to a given io.Writer based on the
// configuraton and datatype
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
