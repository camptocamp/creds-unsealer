package providers

import (
	"fmt"
	"os"

	"github.com/camptocamp/creds-unsealer/backends"
	"github.com/camptocamp/creds-unsealer/config"
	log "github.com/sirupsen/logrus"
)

// Provider is an interface used to abstract the provider
type Provider interface {
	GetName() string
	GetBackend() backends.Backend
	GetInputPath() string
	GetOutputPath() string
	Unseal(string) error
}

// BaseProvider implements a base Provider
type BaseProvider struct {
	backend         backends.Backend
	inputPath       string
	outputPath      string
	outputKeyPrefix string
}

// List returns providers declared in config
func List(cfg *config.Config) (providers []Provider, err error) {
	backend, err := backends.GetBackend(cfg)
	if err != nil {
		return providers, err
	}

	var p Provider
	for _, provider := range cfg.Providers {
		switch provider {
		case "ovh":
			p = &OVH{
				BaseProvider: &BaseProvider{
					backend:         backend,
					inputPath:       cfg.OVH.InputPath,
					outputPath:      os.ExpandEnv(cfg.OutputPathBaseDir + "/.ovh.conf"),
					outputKeyPrefix: cfg.OutputKeyPrefix,
				},
			}
		case "aws":
			p = &AWS{
				BaseProvider: &BaseProvider{
					backend:         backend,
					inputPath:       cfg.AWS.InputPath,
					outputPath:      os.ExpandEnv(cfg.OutputPathBaseDir + "/.aws/credentials"),
					outputKeyPrefix: cfg.OutputKeyPrefix,
				},
			}
		case "openstack":
			p = &Openstack{
				BaseProvider: &BaseProvider{
					backend:         backend,
					inputPath:       cfg.Openstack.InputPath,
					outputPath:      os.ExpandEnv(cfg.OutputPathBaseDir + "/.config/openstack/clouds.yaml"),
					outputKeyPrefix: cfg.OutputKeyPrefix,
				},
			}
		case "file":
			p = &File{
				BaseProvider: &BaseProvider{
					backend:   backend,
					inputPath: cfg.File.InputPath,
				},
			}
		}
		providers = append(providers, p)
	}
	return
}

// UnsealAll unseals all secrets from the backend and add them to the config file
func UnsealAll(p Provider) (err error) {
	secrets, err := p.GetBackend().ListSecrets(p.GetInputPath())

	log.WithFields(log.Fields{
		"provider": p.GetName(),
	}).Debugf("Retrieved secrets: %v", secrets)

	for _, secret := range secrets {
		log.Debugf("Unsealing secret '%s'", secret)
		err = p.Unseal(secret)
		if err != nil {
			return fmt.Errorf("failed to unseal secret: %s", err)
		}
	}
	return
}

// GetName return the provider's name
func (p *BaseProvider) GetName() string {
	return "Base"
}

// GetBackend return the provider's backend
func (p *BaseProvider) GetBackend() backends.Backend {
	return p.backend
}

// GetInputPath return the provider's input path
func (p *BaseProvider) GetInputPath() string {
	return p.inputPath
}

// GetOutputPath return the provider's output path
func (p *BaseProvider) GetOutputPath() string {
	return p.outputPath
}

// Unseal unseals a secret
func (p *BaseProvider) Unseal(s string) error {
	return nil
}
