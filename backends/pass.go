package backends

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/yaml.v2"

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

func (p *Pass) GetCredentials(inputPath string) (credentials interface{}, err error) {
	secret, err := p.decryptSecret(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %s", err)
	}

	err = yaml.Unmarshal(secret, &credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %s", err)
	}
	return
}

func (p *Pass) decryptSecret(path string) ([]byte, error) {
	// Using gpg binary.
	// The golang package for opengpg still doesn't support well GPG 2.1.
	// I don't want to use gopass packages because this backend should work with the old pass.

	secretPath := p.Path + "/" + path + ".gpg"
	out, err := exec.Command("gpg", "--decrypt", secretPath).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %s", err)
	}
	re := regexp.MustCompile(`(?ms)---(.*)`)

	content := re.Find(out)
	if content == nil {
		return nil, fmt.Errorf("found empty content")
	}
	return content, nil
}
