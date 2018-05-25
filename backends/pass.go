package backends

import (
	"os"

	"github.com/camptocamp/creds-unsealer/config"
)

type Pass struct {
	Path string
}

func NewPassBackend(c *config.Config) (p *Pass) {
	p = &Pass{
		Path: os.ExpandEnv(c.Pass.Path),
	}

	return
}

func (p *Pass) GetName() string {
	return "Pass"
}
