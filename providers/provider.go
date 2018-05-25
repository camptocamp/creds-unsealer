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

	for _, p := range cfg.Providers {
		provider, _ := getProvider(backend, p)
		providers = append(providers, provider)
	}
	return
}

func getProvider(backend backends.Backend, provider string) (p Provider, err error) {
	switch provider {
	case "ovh":
		p = &OVH{
			Backend: backend,
		}
	}
	return
}
