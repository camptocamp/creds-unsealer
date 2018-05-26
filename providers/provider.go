package providers

import (
	"github.com/camptocamp/creds-unsealer/backends"
	"github.com/camptocamp/creds-unsealer/config"
)

type Provider interface {
	GetName() string
	GetOutputPath() string
	Unseal() error
}

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
				Backend:   backend,
				InputPath: cfg.Provider.InputPath,
			}
		}
		providers = append(providers, p)
	}
	return
}
