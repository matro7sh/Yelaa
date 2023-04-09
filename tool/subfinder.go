package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type SubfinderConfiguration struct {
	Filename string `json:"filename"`
}

type Subfinder struct {
	filename       string
	runnerInstance *runner.Runner
	proxy          string
}

func (s *Subfinder) Info(_ string) {
	color.Cyan("Searching for subdomains with subfinder")
	color.Yellow("[!] Subfinder only run passive recon on domain, it may not find all the subdomains !")
}

func (s *Subfinder) Configure(conf interface{}) {
	b, _ := json.Marshal(conf.(map[string]interface{}))
	var subfinderconfiguration SubfinderConfiguration
	_ = json.Unmarshal(b, &subfinderconfiguration)

	s.filename = subfinderconfiguration.Filename

	s.runnerInstance, _ = runner.NewRunner(&runner.Options{
		Threads:            3,
		Timeout:            30,
		MaxEnumerationTime: 10,
		Proxy:              conf.(map[string]interface{})["proxy"].(string),
		YAMLConfig: runner.ConfigFile{
			Resolvers:  resolve.DefaultResolvers,
			Sources:    passive.DefaultSources,
			AllSources: passive.DefaultAllSources,
			Recursive:  passive.DefaultRecursiveSources,
		},
	})
}

func (s *Subfinder) Run(website string) {
	website = parseDomain(website)

	buf := bytes.Buffer{}
	err := s.runnerInstance.EnumerateSingleDomain(context.Background(), website, []io.Writer{&buf})
	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile(s.filename, buf.Bytes(), 0755)
}

var _ ToolInterface = (*Subfinder)(nil)
