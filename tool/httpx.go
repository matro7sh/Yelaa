package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"

	"github.com/projectdiscovery/httpx/runner"
)

type HttpxConfiguration struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Httpx struct {
	configuration  HttpxConfiguration
	runnerInstance *runner.Runner
}

func (h *Httpx) Info(_ string) {
	color.Cyan("Running httpx on subdomains")
}

func (h *Httpx) Configure(config interface{}) {
	b, _ := json.Marshal(config.(map[string]interface{}))
	var httpxconfiguration HttpxConfiguration
	_ = json.Unmarshal(b, &httpxconfiguration)
	h.configuration = httpxconfiguration

	opts := runner.Options{
		InputFile:         httpxconfiguration.Input,
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
	}

	if httpxconfiguration.Output != "" {
		opts.Output = httpxconfiguration.Output
	}

	h.runnerInstance, _ = runner.New(&opts)
}

func (h *Httpx) Run(_ string) {
	h.runnerInstance.RunEnumeration()

	body, _ := ioutil.ReadFile(h.configuration.Output)
	fmt.Println(ioutil.WriteFile(h.configuration.Output, regexp.MustCompile(`( \[.+)`).ReplaceAll(body, []byte{}), 0755))
}
