package tool

import (
	"os"
	"path/filepath"

	"github.com/CMEPW/Yelaa/helper"
	internal_runner "github.com/CMEPW/Yelaa/tool/override/nuclei/runner"
	"github.com/fatih/color"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolinit"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

type Nuclei struct{}

func runWrapper(opts types.Options) {
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

func installTemplatesIfNeeded(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		color.Cyan("Installing nuclei templates...")
		opts := types.Options{
			Targets:      []string{"-update-templates"},
			NoInteractsh: true,
		}
		runWrapper(opts)
		color.Cyan("Restarting nuclei")
	}
}

func (*Nuclei) Configure(c interface{}) {}

func (*Nuclei) Info(website string) {
	color.Cyan("Running Nuclei on %s", website)
}

func (*Nuclei) Run(website string) {
	templates_path := filepath.Join(helper.GetHome(), "nuclei-templates/")
	opts := types.Options{
		TemplatesDirectory: templates_path,
		Targets:            []string{website},
		NoInteractsh:       true,
		Output:             "scan_log_nuclei.txt",
	}

	installTemplatesIfNeeded(templates_path)
	runWrapper(opts)
}
