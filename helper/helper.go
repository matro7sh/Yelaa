package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var YelaaPath = GetHome() + "/.yelaa"

func GetHome() (home string) {
	home, _ = os.UserHomeDir()
	return
}

func GetHttpTransport() *http.Transport {
	var proxy = os.Getenv("HTTP_PROXY")
	url, err := url.Parse(proxy)

	if proxy != "" && err == nil {
		return &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(url),
		}
	}
	return &http.Transport{}
}

func GetUserAgent() string {
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
}

func GetCurrentIP() string {
	ht := GetHttpTransport()
	ua := GetUserAgent()

	cli := &http.Client{
		Transport: ht,
	}

	req, err := http.NewRequest("GET", "http://icanhazip.com", nil)
	if err != nil {
		fmt.Printf("[!] Could not query public IP address: %s\n", err.Error())
		return "127.0.0.1" // if the service is offline or not reachable, we should be able to keep going
	}

	req.Header.Add("User-Agent", ua)
	resp, err := cli.Do(req)

	if err != nil || resp == nil {
		fmt.Printf("[!] Could not query public IP address: %s\n", err.Error())
		return ""
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[!] Could not IP address: %s\n", err.Error())
		return ""
	}

	return string(result)
}
