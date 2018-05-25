package providers

import (
	"fmt"
	"os"

	"github.com/camptocamp/creds-unsealer/backends"
)

type OVH struct {
	Backend    backends.Backend
	InputPath  string
	OutputPath string
}

func (o *OVH) GetName() string {
	return "OVH"
}

func (o *OVH) GetOutputPath() string {
	return os.ExpandEnv("$HOME/.ovh.cfg")
}

func (o *OVH) Unseal() (err error) {
	fmt.Println(o.Backend.GetName())
	return
}
