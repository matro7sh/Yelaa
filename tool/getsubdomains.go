package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

type Response []struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
}

func GetSubdomains(_url string) {
	var result Response

	parsed_url, err := url.Parse(_url)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	resp, err := http.Get("https://crt.sh/?q=" + parsed_url.Host + "&output=json")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	if err := json.Unmarshal(body, &result); err != nil {
		color.Red("Can not unmarshal JSON:")
		color.Red(string(body))
		return
	}

	subdomains := make([]string, 0)

	for _, subdomain := range result {
		if contains(subdomains, subdomain.CommonName) {
			continue
		}

		subdomains = append(subdomains, subdomain.CommonName)
	}

	if len(subdomains) == 0 {
		color.Yellow("No subdomains found")
		return
	}

	color.Cyan("------------------")
	color.Green(strings.Join(subdomains, "\n"))
	color.Cyan("------------------")
}
