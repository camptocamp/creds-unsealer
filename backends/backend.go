package backends

import (
	"fmt"

	"github.com/camptocamp/creds-unsealer/config"
)

// Backend is an interface used to abstract the secret manager.
type Backend interface {
	GetName() string
	ListSecrets(string) ([]string, error)
	GetSecret(string, interface{}) error
}

// GetBackend returns a backend interface based on config
func GetBackend(c *config.Config) (b Backend, err error) {
	switch c.Backend {
	case "pass":
		b = NewPassBackend(c)
	default:
		err = fmt.Errorf("Unknown backend: %s", c.Backend)
	}
	return
}
