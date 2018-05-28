package backends

import (
	"fmt"

	"github.com/camptocamp/creds-unsealer/config"
)

type Backend interface {
	GetName() string
	ListCredentials(string) ([]string, error)
	GetSecret(string, interface{}) error
}

func GetBackend(c *config.Config) (b Backend, err error) {
	switch c.Backend {
	case "pass":
		b = NewPassBackend(c)
	default:
		err = fmt.Errorf("Unknown backend: %s", c.Backend)
	}
	return
}
