package cli

import (
	"fmt"
	"os/user"
	"runtime"

	"github.com/joyent/kosh/conch"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
}

type Config struct {
	Version string
	GitRev  string

	ConchURL         string
	ConchToken       string
	ConchEnvironment string

	OutputJSON  bool
	StrictParse bool
	DevMode     bool

	Logger
}

func NewConfig(Version, GitRev string) Config {
	return Config{
		Version: Version,
		GitRev:  GitRev,
	}
}

func BuildUserAgent(GitRev string) map[string]string {
	var isRoot bool
	if current, err := user.Current(); err == nil {
		if current.Uid == "0" {
			isRoot = true
		}
	}

	agentBits := make(map[string]string)
	agent := fmt.Sprintf(
		"%s (%s; %s; r=%v)",
		GitRev,
		runtime.GOOS,
		runtime.GOARCH,
		isRoot,
	)

	agentBits["Kosh"] = agent
	return agentBits
}

func (c Config) ConchClient() *conch.Client {
	return conch.New(c.ConchURL, c.ConchToken, c.Logger)
}
