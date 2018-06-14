package providers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/justwatchcom/gopass/utils/fsutil"
	log "github.com/sirupsen/logrus"
)

// File represents a provider
type File struct {
	*BaseProvider
}

// FileConfig stores a file configuration
type FileConfig struct {
	Path    string      `yaml:"path"`
	Content string      `yaml:"content"`
	Group   string      `yaml:"group,omitempty"`
	Mode    os.FileMode `yaml:"mode,omitempty"`
	Owner   string      `yaml:"owner,omitempty"`
}

// GetName returns the provider's name
func (f *File) GetName() string {
	return "File"
}

// Unseal unseals a secret from the backend and add it to the config file
func (f *File) Unseal(cred string) (err error) {
	var secret FileConfig

	err = f.backend.GetSecret(f.inputPath+"/"+cred, &secret)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %s", err)
	}
	if secret.Mode == os.FileMode(0000) {
		secret.Mode = os.FileMode(0666)
	}

	outputFile := fsutil.CleanPath(os.ExpandEnv(secret.Path))
	log.WithFields(log.Fields{
		"provider": f.GetName(),
	}).Debugf("Output path: %s", outputFile)

	dir, err := filepath.Abs(filepath.Dir(outputFile))
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials' absolute dir path: %s", err)
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("failed to create dir: %s", err)
	}
	log.WithFields(log.Fields{
		"provider": f.GetName(),
	}).Debugf("Write file %s with mode %s", outputFile, secret.Mode)
	err = ioutil.WriteFile(outputFile, []byte(secret.Content), secret.Mode)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %s", err)
	}
	err = os.Chmod(outputFile, secret.Mode)
	if err != nil {
		return fmt.Errorf("failed to change file mode: %s", err)
	}

	return
}
