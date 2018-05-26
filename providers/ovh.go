package providers

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/camptocamp/creds-unsealer/backends"
)

type OVH struct {
	Backend    backends.Backend
	InputPath  string
	OutputPath string
}

type OVHConfig struct {
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
}

func (o *OVH) GetName() string {
	return "OVH"
}

func (o *OVH) GetOutputPath() string {
	return os.ExpandEnv("$HOME/.ovh.cfg")
}

func (o *OVH) Unseal() (err error) {

	// o.buildConfig(o.OuputPath)

	log.Info("Using provider OVH")

	_, err = o.Backend.GetCredentials(o.InputPath)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}
	return
}
