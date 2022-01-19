package tool

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

type Robot struct{}

func (s *Robot) Info(website string) {
	color.Cyan("Looking for robots.txt on: %s", website)
}

func (g *Robot) Configure(c interface{}) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func (g *Robot) Run(domain string) {

	domain = strings.TrimSuffix(domain, "/")

	for _, u := range getUrls(domain) {
		resp, err := http.Get(fmt.Sprint(u, "/robots.txt"))
		if err != nil {
			fmt.Printf("%v", err)
		}

		if resp != nil && resp.StatusCode != http.StatusNotFound {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("%v", err)
			}
			sb := string(body)
			color.Green(sb)

		} else {
			color.Yellow("----- Sorry get 404 status code for this robots.txt ----- ")
		}
	}

}

var _ ToolInterface = (*Robot)(nil)
