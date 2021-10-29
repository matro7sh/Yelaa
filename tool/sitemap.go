package tool

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

var GlobalHeaders = []string{"Server", "X-XSS-Protection", "Access-Control-Allow-Credentials", "Content-Security-Policy", "X-Powered-By", "Strict-Transport-Security"}

func GetSitemap(url string) {
	resp, err := http.Get(url + "sitemap.xml")

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
				//		fmt.Println("Header + " + headerName + "Found : " + headerValue)
				fmt.Printf("Found Header: %s | %s \n", headerName, headerValue)
			} else {
				// fmt.Println("sorry no headers in global headers variable")
			}
		}
		sb := string(body)
		color.Green(sb)

	} else {
		color.Yellow("-----  Sorry get 404 status code for this sitemap.xml -----")
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
