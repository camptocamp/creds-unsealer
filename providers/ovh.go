package providers

import (
	"fmt"

	"github.com/raphink/narcissus"
	log "github.com/sirupsen/logrus"
	"honnef.co/go/augeas"

	"github.com/camptocamp/creds-unsealer/backends"
)

// OVH stores data used to manage OVH credentials
type OVH struct {
	Backend    backends.Backend
	InputPath  string
	OutputPath string
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

// UnsealAll unseals all secrets from the backend and add them to the config file
func (o *OVH) UnsealAll() (err error) {
	secrets, err := o.Backend.ListSecrets(o.InputPath)

	log.WithFields(log.Fields{
		"provider": o.GetName(),
	}).Debugf("Retrieved secrets: %v", secrets)

	for _, secret := range secrets {
		err = o.Unseal(secret)
		if err != nil {
			return fmt.Errorf("failed to unseal secret: %s", err)
		}
	}
	return
}

// Unseal unseals a secret from the backend and add it to the config file
func (o *OVH) Unseal(cred string) (err error) {
	var secret OVHConfig
	err = o.Backend.GetSecret(o.InputPath+cred, &secret)
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
		augeasFile: o.OutputPath,
		augeasPath: "/files" + o.OutputPath,
	}
	configs.Configs = make(map[string]OVHConfig)
	configs.Configs[name] = config

	err = n.Write(&configs)
	if err != nil {
		return fmt.Errorf("failed to write configs: %s", err)
	}

	return
}
