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

var GlobalHeaders = []string{"Server", "X-XSS-Protection", "Access-Control-Allow-Credentials", "Content-Security-Policy", "X-Powered-By", "Strict-Transport-Security"}

type Sitemap struct{}

func (s *Sitemap) Info(website string) {
	color.Cyan("Looking for sitemap.xml on: %s ", website)
}

func (s *Sitemap) Configure(c interface{}) {}

func (s *Sitemap) Run(domain string) {
	domain = strings.TrimSuffix(domain, "/")

    transport := helper.GetHttpTransport()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

    client := &http.Client{
        Transport: transport,
    }

	for _, u := range getUrls(domain) {
        req, err := http.NewRequest("GET", fmt.Sprint(u, "/sitemap.xml"), nil)
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

			for headerName, headerValue := range resp.Header {
				if contains(GlobalHeaders, headerName) {
					fmt.Printf("Found Header: %s | %s \n", headerName, headerValue)
				}
			}
			sb := string(body)
			color.Green(sb)

		} else {
			color.Yellow("-----  Sorry, got 404 status code for sitemap.xml -----")
		}
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

var _ ToolInterface = (*Sitemap)(nil)
