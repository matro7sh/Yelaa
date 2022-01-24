package tool

import (
    "path/filepath"
    "os"

    "github.com/CMEPW/Yelaa/helper"
    internal_runner "github.com/CMEPW/Yelaa/tool/override/nuclei/runner"
    "github.com/fatih/color"
    "github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolinit"
    "github.com/projectdiscovery/nuclei/v2/pkg/types"
)

var (
    nuclei_template_path string
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

func installTemplatesIfNeeded() {
    if _, err := os.Stat(nuclei_template_path); os.IsNotExist(err) {
        color.Cyan("Installing nuclei templates...")
        opts := types.Options {
            Targets:        []string{"-update-templates"},
            NoInteractsh:   true,
        }
        runWrapper(opts)
    }
}

func (*Nuclei) Configure(c interface{}) {}

func (*Nuclei) Info(website string) {
    color.Cyan("Running Nuclei on %s", website)
}

func (*Nuclei) Run(website string) {
    nuclei_template_path := filepath.Join(helper.GetHome(), "nuclei-templates")

    opts := types.Options{
        TemplatesDirectory: nuclei_template_path,
        Targets:            []string{website},
        NoInteractsh:       true,
    }
    installTemplatesIfNeeded()
    runWrapper(opts)
}
