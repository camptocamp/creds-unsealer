package providers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Openstack represents a provider
type Openstack struct {
	*BaseProvider
}

// OpenstackConfigAuth stores authentication informations
type OpenstackConfigAuth struct {
	AuthURL        string `yaml:"auth_url,omitempty"`
	Password       string `yaml:"password,omitempty"`
	ProjectID      string `yaml:"project_id,omitempty"`
	ProjectName    string `yaml:"project_name,omitempty"`
	Username       string `yaml:"username,omitempty"`
	UserDomainName string `yaml:"user_domain_name,omitempty"`
}

// OpenstackConfig stores Openstack's cloud config
type OpenstackConfig struct {
	Auth               *OpenstackConfigAuth
	RegionName         string `yaml:"region_name,omitempty"`
	IdentityAPIVersion string `yaml:"identity_api_version,omitempty"`
	Interface          string `yaml:"interface,omitempty"`
	Cert               string `yaml:"cert,omitempty"`
	Key                string `yaml:"key,omitempty"`
}

// OpenstackClouds represents the Openstack's clouds.yaml
type OpenstackClouds struct {
	Clouds map[string]*OpenstackConfig `yaml:"clouds"`
}

// GetName returns the provider's name
func (o *Openstack) GetName() string {
	return "Openstack"
}

// Unseal unseals a secret from the backend and add it to the config file
func (o *Openstack) Unseal(cred string) (err error) {
	var secret OpenstackConfig
	err = o.backend.GetSecret(o.inputPath+"/"+cred, &secret)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}

	// Expand certificates path
	secret.Cert, _ = homedir.Expand(secret.Cert)
	secret.Key, _ = homedir.Expand(secret.Key)

	err = o.writeSecret(cred, secret)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %s", err)
	}

	return
}

func (o *Openstack) writeSecret(name string, config OpenstackConfig) (err error) {
	f, err := os.OpenFile(o.outputPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %s", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read config file: %s", err)
	}

	var t *OpenstackClouds
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for cloudKey := range t.Clouds {
		if cloudKey == o.outputKeyPrefix+name {
			t.Clouds[cloudKey] = &config
			return
		}
	}

	t.Clouds[o.outputKeyPrefix+name] = &config

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	_, err = f.WriteAt(d, 0)
	if err != nil {
		return fmt.Errorf("failed to write config file: %s", err)
	}
	return
}
