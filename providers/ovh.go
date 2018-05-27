package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"honnef.co/go/augeas"

	"github.com/camptocamp/creds-unsealer/backends"
)

type OVH struct {
	Backend    backends.Backend
	InputPath  string
	OutputPath string
}

type OVHConfig struct {
	ApplicationKey    string `json:"application_key,omitempty"`
	ApplicationSecret string `json:"application_secret,omitempty"`
	ConsumerKey       string `json:"consumer_key,omitempty"`
}

func (o *OVH) GetName() string {
	return "OVH"
}

func (o *OVH) GetOutputPath() string {
	return os.ExpandEnv("$HOME/.ovh.cfg")
}

func (o *OVH) Unseal() (err error) {
	iCreds, err := o.Backend.GetCredentials(o.InputPath)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}

	bCreds, err := json.Marshal(iCreds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %s", err)
	}

	var config OVHConfig
	err = json.Unmarshal(bCreds, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal credentials: %s", err)
	}

	o.buildConfig(o.GetOutputPath(), &config)

	return
}

func (o *OVH) buildConfig(path string, config *OVHConfig) (err error) {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		return fmt.Errorf("failed to load augeas: %s", err)
	}

	err = aug.Set("/augeas/load/IniFile/lens", "Puppet.lns")
	if err != nil {
		return fmt.Errorf("failed to set augeas lens: %s", err)
	}
	err = aug.Set("/augeas/load/IniFile/incl", path)
	if err != nil {
		return fmt.Errorf("failed to set augeas incl: %s", err)
	}

	keys := strings.Split(o.InputPath, "/")
	resourceKey := keys[len(keys)-1]
	matches, err := aug.Match("/files" + path + "/" + resourceKey)
	if err != nil {
		return fmt.Errorf("failed to list augeas resources: %s", err)
	}
	if len(matches) == 0 {
		// Create node
	}
	return
}
