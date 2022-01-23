package tool

import (
	"path/filepath"

	"github.com/CMEPW/Yelaa/helper"
	internal_runner "github.com/CMEPW/Yelaa/tool/override/nuclei/runner"
	"github.com/fatih/color"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolinit"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

type Nuclei struct{}

func (*Nuclei) Configure(c interface{}) {}

func (*Nuclei) Info(website string) {
	color.Cyan("Running Nuclei on %s", website)
}

func (*Nuclei) Run(website string) {
	opts := types.Options{
		TemplatesDirectory: filepath.Join(helper.GetHome(), "nuclei-templates"),
		Targets:            []string{website},
		NoInteractsh:       true,
	}
	protocolinit.Init(&opts)
	r, err := internal_runner.New(&opts)

	if err != nil {
		color.Yellow("%s", err)
	}

	if err := r.RunEnumeration(); err != nil {
		color.Red("Could not run nuclei: %s\n", err)
	}
	r.Close()
}
