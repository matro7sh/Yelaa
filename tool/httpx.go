package tool

import (
	"encoding/json"

	"github.com/fatih/color"
	// "github.com/projectdiscovery/httpx/common/httpx"

	"github.com/projectdiscovery/httpx/runner"
)

type HttpxConfiguration struct {
	SubdomainsFilename string `json:"subdomainsFilename"`
}

type Httpx struct {
	runnerInstance *runner.Runner
}

func (h *Httpx) Info(_ string) {
	color.Cyan("Running httpx on subdomains")
}

func (h *Httpx) Configure(config interface{}) {
	b, _ := json.Marshal(config.(map[string]interface{}))
	var httpxconfiguration HttpxConfiguration
	_ = json.Unmarshal(b, &httpxconfiguration)

	h.runnerInstance, _ = runner.New(&runner.Options{
		InputFile:         httpxconfiguration.SubdomainsFilename,
		ExtractTitle:      true,
		ContentLength:     true,
		OutputContentType: true,
		StatusCode:        true,
		TechDetect:        true,
		VHost:             true,
		OutputWebSocket:   true,
		FollowRedirects:   true,
		Retries:           2,
		Threads:           50,
	})
}

func (h *Httpx) Run(_ string) {
	h.runnerInstance.RunEnumeration()
}
