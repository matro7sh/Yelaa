package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CMEPW/Yelaa/helper"
	internal_runner "github.com/CMEPW/Yelaa/tool/override/nuclei/runner"
	"github.com/fatih/color"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolinit"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

type Nuclei struct {
    proxy       string
	rateLimiter int32
}

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

func (n *Nuclei) Configure(c interface{}) {
	outputDir := helper.YelaaPath + "/nuclei"
	n.rateLimiter = c.(map[string]interface{})["rateLimiter"].(int32)

    proxy := c.(map[string]interface{})["proxy"].(string)
    n.proxy = proxy

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0750); err != nil {
			fmt.Println(err)
		}
	}
}

func (*Nuclei) Info(website string) {
	color.Cyan("Running Nuclei on %s", website)
}

func (*Nuclei) Run(website string) {
	outputFile := helper.YelaaPath + "/nuclei/scan_log_nuclei-" + time.Now().Format("2006-01-02_15-04-05") + ".txt"
	templates_path := filepath.Join(helper.GetHome(), "nuclei-templates/")
	opts := types.Options{
		TemplatesDirectory: templates_path,
		Targets:            []string{website},
		NoInteractsh:       true,
		Output:             outputFile,
	}

	installTemplatesIfNeeded(templates_path)
	runWrapper(opts)
}
