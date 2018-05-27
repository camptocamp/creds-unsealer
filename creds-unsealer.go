package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/camptocamp/creds-unsealer/config"
	"github.com/camptocamp/creds-unsealer/providers"
)

var version string = "undefined"
var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig(version)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
}

func main() {

	providers, _ := providers.List(cfg)

	for _, p := range providers {
		log.Infof("Using provider %s", p.GetName())
		err := p.Unseal()
		if err != nil {
			log.Errorf("failed to unseal: %s", err)
		}
	}
}
