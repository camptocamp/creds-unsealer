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

	pp, _ := providers.List(cfg)

	for _, p := range pp {
		log.Infof("Using provider %s", p.GetName())
		err := providers.UnsealAll(p)
		if err != nil {
			log.Errorf("failed to unseal: %s", err)
		}
	}
}
