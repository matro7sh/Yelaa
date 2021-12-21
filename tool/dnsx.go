package tool

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

type DnsxConfiguration struct {
	SubdomainsFilename string `json:"subdomainsFilename"`
	IpsFilename        string `json:"ipsFilename"`
}

type Dnsx struct {
	opts               *dnsx.Options
	subdomainsFilename string
	ipsFilename        string
}

func (d *Dnsx) Info(_ string) {
	color.Cyan("Running dnsx on subdomains to find IP address")
}

func (d *Dnsx) Configure(config interface{}) {
	d.opts = &dnsx.DefaultOptions

	b, _ := json.Marshal(config.(map[string]interface{}))
	var dnsxConfig DnsxConfiguration
	_ = json.Unmarshal(b, &dnsxConfig)

	d.subdomainsFilename = dnsxConfig.SubdomainsFilename
	d.ipsFilename = dnsxConfig.IpsFilename
}

func (d *Dnsx) Run(_ string) {
	runner, err := dnsx.New(*d.opts)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Open(d.subdomainsFilename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	ip_list := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		data, err := runner.QueryOne(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}

		if data == nil {
			continue
		}

		color.Green(strings.Join(data.A, "\n"))
		ip_list = append(ip_list, data.A...)
	}

	output := strings.Join(ip_list, "\n")
	os.WriteFile(d.ipsFilename, []byte(output), 0644)
}

var _ ToolInterface = (*Dnsx)(nil)
