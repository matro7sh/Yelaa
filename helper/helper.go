package helper

import (
	"net/http"
	"os"
)

var YelaaPath = GetHome() + "/.yelaa"

func GetHome() (home string) {
	home, _ = os.UserHomeDir()
	return
}

func GetHttpTransport() (*http.Transport) {
    var proxy = os.Getenv("HTTP_PROXY")

    if proxy != "" {
        return &http.Transport{
			DisableKeepAlives: true,
            Proxy: http.ProxyFromEnvironment,
        }
    }
    return &http.Transport{}
}
