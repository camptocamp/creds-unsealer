package providers

import (
	"fmt"
	"os"

	"github.com/camptocamp/creds-unsealer/backends"
	"github.com/raphink/narcissus"
	"honnef.co/go/augeas"
)

type OVH struct {
	Backend    backends.Backend
	InputPath  string
	OutputPath string
}

type OVHConfig struct {
	Name              string `yaml:"omitempty"`
	ApplicationKey    string `yaml:"application_key,omitempty" path:"application_key"`
	ApplicationSecret string `yaml:"application_secret,omitempty" path:"application_secret"`
	ConsumerKey       string `yaml:"consumer_key,omitempty" path:"consumer_key"`
}

type OVHConfigs struct {
	augeasPath string
	Configs    map[string]OVHConfig `path:"section" purge:"false"`
}

func (o *OVH) GetName() string {
	return "OVH"
}

func (o *OVH) GetOutputPath() string {
	return os.ExpandEnv("$HOME/.ovh.cfg")
}

func (o *OVH) UnsealAll() (err error) {
	creds, err := o.Backend.ListCredentials(o.InputPath)
	for _, cred := range creds {
		err = o.Unseal(cred)
		if err != nil {
			return fmt.Errorf("failed to unseal secret: %s", err)
		}
	}
	return
}

func (o *OVH) Unseal(cred string) (err error) {
	var secret OVHConfig
	err = o.Backend.GetSecret(o.InputPath+cred, &secret)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}
	secret.Name = cred

	err = o.writeSecret(o.GetOutputPath(), secret)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %s", err)
	}

	return
}

func (o *OVH) writeSecret(path string, config OVHConfig) (err error) {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		return fmt.Errorf("failed to initialize Augeas: %s", err)
	}

	err = aug.Transform("IniFile.lns_loose", o.OutputPath, false)
	if err != nil {
		return fmt.Errorf("failed to set up Augeas transform: %s", err)
	}

	err = aug.Load()
	if err != nil {
		return fmt.Errorf("failed to load Augeas tree: %s", err)
	}

	n := narcissus.New(&aug)
	configs := OVHConfigs{
		augeasPath: "/files" + o.OutputPath,
	}
	configs.Configs = make(map[string]OVHConfig)
	configs.Configs[config.Name] = config

	err = n.Write(&configs)
	if err != nil {
		return fmt.Errorf("failed to write configs: %s", err)
	}

	return
}
