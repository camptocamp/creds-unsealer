package providers

import (
	"os"

	"github.com/camptocamp/creds-unsealer/backends"
	"github.com/camptocamp/creds-unsealer/config"
)

type Provider interface {
	GetName() string
	GetOutputPath() string
	UnsealAll() error
	Unseal(string) error
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
			var inputPath string
			if inputPath = cfg.Provider.InputPath; inputPath == "" {
				inputPath = "/ovh"
			}
			p = &OVH{
				Backend:    backend,
				InputPath:  inputPath,
				OutputPath: os.ExpandEnv("$HOME/.ovh.conf"),
			}
		}
		providers = append(providers, p)
	}
	return
}
