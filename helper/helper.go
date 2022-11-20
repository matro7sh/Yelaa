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

func GetUserAgent() string {
    return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
}
