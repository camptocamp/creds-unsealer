package config

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version   bool     `short:"V" long:"version" description:"Display version."`
	Loglevel  string   `short:"l" long:"loglevel" description:"Set loglevel ('debug', 'info', 'warn', 'error', 'fatal', 'panic')." env:"BIVAC_LOG_LEVEL" default:"info"`
	Manpage   bool     `short:"m" long:"manpage" description:"Output manpage."`
	Backend   string   `short:"b" long:"backend" description:"Backend to use." env:"CREDS_BACKEND" default:"pass"`
	Providers []string `short:"p" long:"providers" description:"Providers to use." env:"CREDS_PROVIDERS" default:"ovh"`
	Pass      struct {
		Path string `long:"pass-path" description:"Path to password-store." env:"CREDS_PASS_PATH" default:"$HOME/.password-store"`
	} `group:"Pass backend options"`

	Provider struct {
		InputPath string `long:"provider-input-path" description:"Provider input path"`
	} `group:"Provider options"`
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) (*Config, error) {
	var c Config
	parser := flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse config: %s", err)
	}

	if c.Version {
		fmt.Printf("Creds-unsealer %v\n", version)
		os.Exit(0)
	}

	if c.Manpage {
		os.Exit(0)
	}

	return &c, nil
}
