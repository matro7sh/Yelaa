package tool

import (
	"fmt"
	"os"
	"os/exec"
)

func run(url string) ([]byte, error) {
    proxy := os.Getenv("HTTP_PROXY")

    if proxy != "" {
        proxyCmd := fmt.Sprintf("--proxy=%s", proxy)
        return exec.Command("dirsearch", "-u", url, proxyCmd).Output()
    }
    return exec.Command("dirsearch", "-u", url).Output()
}

func Dirsearch(url string) {
    out, err := run(url)
	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	fmt.Println("dirsearch output ", output)
}
