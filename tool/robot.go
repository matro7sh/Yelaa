package tool

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CMEPW/Yelaa/helper"
	"github.com/fatih/color"
)

type Robot struct{}

func (s *Robot) Info(website string) {
	color.Cyan("Looking for robots.txt on: %s", website)
}

func (g *Robot) Configure(c interface{}) {
}

func (g *Robot) Run(domain string) {

	domain = strings.TrimSuffix(domain, "/")

    transport := helper.GetHttpTransport()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

    client := &http.Client{
        Transport: transport,
    }

	for _, u := range getUrls(domain) {

        req, err := http.NewRequest("GET", fmt.Sprint(u, "/robots.txt"), nil)
        req.Header.Add("User-Agent", helper.GetUserAgent())

        resp, err := client.Do(req)

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
			color.Yellow("----- Sorry, got 404 status code for robots.txt ----- ")
		}
	}

}

var _ ToolInterface = (*Robot)(nil)
