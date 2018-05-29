package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/camptocamp/creds-unsealer/config"
	"github.com/camptocamp/creds-unsealer/providers"
)

var version = "undefined"
var cfg *config.Config

func init() {
	cfg = config.LoadConfig(version)
}

func main() {

	providers, _ := providers.List(cfg)

	for _, p := range providers {
		log.Infof("Using provider %s", p.GetName())
		err := p.UnsealAll()
		if err != nil {
			log.Errorf("failed to unseal: %s", err)
		}
	}
}
