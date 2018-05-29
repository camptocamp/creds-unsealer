package providers

import (
	"os"

	"github.com/camptocamp/creds-unsealer/backends"
	"github.com/camptocamp/creds-unsealer/config"
)

// Provider is an interface used to abstract the provider
type Provider interface {
	GetName() string
	UnsealAll() error
	Unseal(string) error
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
				Backend:    backend,
				InputPath:  cfg.OVH.InputPath,
				OutputPath: os.ExpandEnv("$HOME/.ovh.conf"),
			}
		}
		providers = append(providers, p)
	}
	return
}
