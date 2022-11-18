package helper

import (
	"net/http"
    "net/url"
	"os"
)

var YelaaPath = GetHome() + "/.yelaa"

func GetHome() (home string) {
	home, _ = os.UserHomeDir()
	return
}

func GetHttpTransport() (*http.Transport) {
    var proxy = os.Getenv("HTTP_PROXY")
    url, err := url.Parse(proxy)

    if proxy != "" && err == nil {
        return &http.Transport{
			DisableKeepAlives: true,
            Proxy: http.ProxyURL(url),
        }
    }
    return &http.Transport{}
}
