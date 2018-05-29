package providers

import (
	"fmt"

	"github.com/raphink/narcissus"
	log "github.com/sirupsen/logrus"
	"honnef.co/go/augeas"

	"github.com/camptocamp/creds-unsealer/backends"
)

type OVH struct {
	Backend    backends.Backend
	InputPath  string
	OutputPath string
}

type OVHConfig struct {
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

	err = aug.Transform("IniFile.lns_loose", o.OutputPath, false)
	if err != nil {
		return fmt.Errorf("failed to set up Augeas transform: %s", err)
	}

	err = aug.Load()
	if err != nil {
		return fmt.Errorf("failed to load Augeas tree: %s", err)
	}

	augErr, _ := aug.Get("/augeas/files" + o.OutputPath + "/message")
	if augErr != "" {
		return fmt.Errorf("failed to load file with Augeas: %s", augErr)
	}

	n := narcissus.New(&aug)
	configs := OVHConfigs{
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
