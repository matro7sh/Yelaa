package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/fatih/color"

	customport "github.com/projectdiscovery/httpx/common/customports"
	"github.com/projectdiscovery/httpx/runner"
)

type HttpxConfiguration struct {
	Input  string `json:"input"`
	Output string `json:"output"`
    proxy          string
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
    proxy := config.(map[string]interface{})["proxy"].(string)

	h.configuration = httpxconfiguration
	customPorts := customport.CustomPorts{}
	customPorts.Set("25,80,81,135,389,443,1080,3000,3306,8080,8443,8888,9090,8089")

	opts := runner.Options{
		InputFile:         httpxconfiguration.Input,
		CustomPorts:       customPorts,
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
		Timeout:           8,
        RandomAgent:       true,
	}

    if strings.HasPrefix(proxy, "http") {
        opts.HTTPProxy = proxy
    }

    if strings.HasPrefix(proxy, "socks") {
        opts.SocksProxy = proxy
    }

	if httpxconfiguration.Output != "" {
		opts.Output = httpxconfiguration.Output
	}

	h.runnerInstance, _ = runner.New(&opts)
}

func (h *Httpx) Run(_ string) {
	h.runnerInstance.RunEnumeration()

	body, _ := ioutil.ReadFile(h.configuration.Output)
	fmt.Println(string(body))
	ioutil.WriteFile(h.configuration.Output, regexp.MustCompile(`( \[.+)`).ReplaceAll(body, []byte{}), 0755)
}
