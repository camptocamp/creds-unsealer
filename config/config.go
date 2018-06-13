package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jessevdk/go-flags"
)

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version         bool     `short:"V" long:"version" description:"Display version."`
	LogLevel        string   `short:"l" long:"loglevel" description:"Set loglevel ('debug', 'info', 'warn', 'error', 'fatal', 'panic')." env:"BIVAC_LOG_LEVEL" default:"info"`
	Manpage         bool     `short:"m" long:"manpage" description:"Output manpage."`
	Backend         string   `short:"b" long:"backend" description:"Backend to use." env:"CREDS_BACKEND" default:"pass"`
	Providers       []string `short:"p" long:"providers" description:"Providers to use." env:"CREDS_PROVIDERS" default:"ovh" default:"aws" default:"openstack"`
	OutputKeyPrefix string   `long:"output-key-prefix" description:"String to prepend to key of the secret"`

	// Backends configuration
	Pass struct {
		Path string `long:"backend-pass-path" description:"Path to password-store." env:"CREDS_BACKEND_PASS_PATH"`
	} `group:"Pass backend options"`

	// Providers configuration
	OVH struct {
		InputPath string `long:"provider-ovh-input-path" description:"OVH Provider input path" default:"ovh"`
	} `group:"OVH Provider options"`
	AWS struct {
		InputPath string `long:"provider-aws-input-path" description:"AWS Provider input path" default:"aws"`
	} `group:"AWS Provider options"`
	Openstack struct {
		InputPath string `long:"provider-openstack-input-path" description:"Openstack Provider input path" default:"openstack"`
	} `group:"Openstack Provider options"`
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) *Config {
	var c Config
	parser := flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	if c.Version {
		fmt.Printf("Creds-unsealer %v\n", version)
		os.Exit(0)
	}

	if c.Manpage {
		os.Exit(0)
	}

	err := c.setupLogLevel()
	if err != nil {
		log.Errorf("failed to setup log level: %s", err)
		os.Exit(1)
	}
	return &c
}

func (c *Config) setupLogLevel() (err error) {
	switch c.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		err = fmt.Errorf("Wrong log level '%v'", c.LogLevel)
	}

	return
}
