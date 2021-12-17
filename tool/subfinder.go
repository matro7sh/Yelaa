package tool

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func parseDomain(domain string) string {
	if strings.HasPrefix(domain, "http") {
		parsedUrl, err := url.Parse(domain)
		if err != nil {
			fmt.Printf("%s", err)
		}

		domain = parsedUrl.Host
	}

	return domain
}

func Subfinder(domain string, filename string) {

	domain = parseDomain(domain)

	out, err := exec.Command("subfinder", "-d", domain).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	err = os.WriteFile(filename, out, 0644)

	if err != nil {
		fmt.Printf("%s", err)
	}

}
