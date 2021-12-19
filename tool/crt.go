package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

type ResponseItem struct {
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

func Crt(url string, getSubDomainCrt string) {
	var result []ResponseItem

	domain := "parseDomain(url)"

	resp, err := http.Get("https://crt.sh/?q=" + domain + "&output=json")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	if resp.StatusCode == 502 {
		color.Yellow("[!] Carrefull crt.sh doesn't respond correctly, you must wait before you retry")
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		fmt.Printf("Status code : %v\nBody: %v", resp.StatusCode, body)
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

	err = os.WriteFile(getSubDomainCrt, []byte(strings.Join(subdomains, "\n")+"\n"), 0644)

	if err != nil {
		fmt.Printf("%s", err)
	}

	color.Cyan("------------------")
	color.Green(strings.Join(subdomains, "\n"))
	color.Cyan("------------------")
}
