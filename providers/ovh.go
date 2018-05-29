package providers

import (
	"fmt"

	"github.com/raphink/narcissus"
	"honnef.co/go/augeas"
)

// OVH stores data used to manage OVH credentials
type OVH struct {
	*BaseProvider
}

// OVHConfig represents an entry in the OVH's configuration file
type OVHConfig struct {
	ApplicationKey    string `yaml:"application_key,omitempty" narcissus:"application_key"`
	ApplicationSecret string `yaml:"application_secret,omitempty" narcissus:"application_secret"`
	ConsumerKey       string `yaml:"consumer_key,omitempty" narcissus:"consumer_key"`
}

// OVHConfigs is used by Narcissus to manage entries on the OVH's configuration file
type OVHConfigs struct {
	augeasFile string
	augeasLens string `default:"IniFile.lns_loose"`
	augeasPath string
	Configs    map[string]OVHConfig `narcissus:"section"`
}

// GetName returns the provider's name
func (o *OVH) GetName() string {
	return "OVH"
}

// Unseal unseals a secret from the backend and add it to the config file
func (o *OVH) Unseal(cred string) (err error) {
	var secret OVHConfig
	err = o.backend.GetSecret(o.inputPath+"/"+cred, &secret)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}

	err = o.writeSecret(cred, secret)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %s", err)
	}

	return
}

func (o *OVH) writeSecret(name string, config OVHConfig) (err error) {
	aug, err := augeas.New("/", "", augeas.NoModlAutoload)
	defer aug.Close()
	if err != nil {
		return fmt.Errorf("failed to initialize Augeas: %s", err)
	}

	n := narcissus.New(&aug)
	configs := OVHConfigs{
		augeasFile: o.outputPath,
		augeasPath: "/files" + o.outputPath,
	}
	configs.Configs = make(map[string]OVHConfig)
	configs.Configs[name] = config

	err = n.Write(&configs)
	if err != nil {
		return fmt.Errorf("failed to write configs: %s", err)
	}

	return
}
