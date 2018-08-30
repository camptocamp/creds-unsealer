package backends

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/justwatchcom/gopass/action"
	gpConfig "github.com/justwatchcom/gopass/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/camptocamp/creds-unsealer/config"
)

// Pass stores data used by gopass
type Pass struct {
	Path string
}

// NewPassBackend returns a Pass struct based on config
func NewPassBackend(c *config.Config) (p *Pass) {
	p = &Pass{
		Path: os.ExpandEnv(c.Pass.Path),
	}

	return
}

// GetName returns the backend name
func (p *Pass) GetName() string {
	return "Pass"
}

// ListSecrets returns the list of secrets from gopass
func (p *Pass) ListSecrets(inputPath string) (secrets []string, err error) {
	var gpc *gpConfig.Config
	act, err := p.initGopass(gpc)
	if err != nil {
		return nil, fmt.Errorf("failed to create gopass action: %s", err)
	}

	log.WithFields(log.Fields{
		"backend": p.GetName(),
	}).Debugf("Using path: %s", act.Store.Path())

	rootTree, err := act.Store.Tree(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gopass root tree: %s", err)
	}
	t, err := rootTree.FindFolder(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gopass tree: %s", err)
	}

	l := t.List(0)
	for _, secret := range l {
		s := strings.TrimPrefix(secret, t.String())
		s = strings.TrimPrefix(s, "/")
		log.Debugf("Found secret '%s'", s)
		secrets = append(secrets, s)
	}
	return
}

// GetSecret decrypt a gopass secret and stores it in the interface `secret`
func (p *Pass) GetSecret(inputPath string, secret interface{}) (err error) {
	s, err := p.decryptSecret(inputPath)
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %s", err)
	}
	err = yaml.Unmarshal(s, secret)
	if err != nil {
		return fmt.Errorf("failed to unmarshal secret: %s", err)
	}
	return
}

func (p *Pass) decryptSecret(path string) (content []byte, err error) {
	var gpc *gpConfig.Config
	act, err := p.initGopass(gpc)
	if err != nil {
		return nil, fmt.Errorf("failed to create gopass action: %s", err)
	}

	log.Debugf("Getting secret at %s", path)
	sec, err := act.Store.Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret at %s: %s", path, err)
	}

	body, err := sec.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %s", err)
	}

	re := regexp.MustCompile(`(?ms)---(.*)`)
	content = re.Find(body)
	if content == nil {
		return nil, fmt.Errorf("found empty content")
	}

	return
}

func (p *Pass) initGopass(gpc *gpConfig.Config) (act *action.Action, err error) {
	if p.Path == "" {
		gpc = gpConfig.Load()
	} else {
		gpc = gpConfig.New()
		s := &gpConfig.StoreConfig{
			Path: p.Path,
		}
		gpc.Root = s
	}
	act, err = action.New(context.Background(), gpc, semver.Version{})
	return
}
