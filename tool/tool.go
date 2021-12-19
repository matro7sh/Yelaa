package tool

import (
	"fmt"
	"net/url"
	"strings"
)

type ToolInterface interface {
	Configure(interface{})
	Info(string)
	Run(string)
}

func parseDomain(website string) string {
	if strings.HasPrefix(website, "http") {
		parsedUrl, err := url.Parse(website)
		if err != nil {
			fmt.Printf("%s", err)
		}

		website = parsedUrl.Host
	}

	return website
}
