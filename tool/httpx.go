package tool

import (
	"fmt"
	"os/exec"
)

func Httpx(ipsFile string) {
	out, err := exec.Command("httpx", "-l", ipsFile, "-title", "-content-length", "-content-type",
		"-status-code", "-tech-detect", "-vhost", "-websocket", "-follow-redirects",
		"-ports", "25,80,81,135,389,443,1080,3000,3306,8080,8443,8888,9090,8089",
		"-retries", "2", "-timeout", "8", "-threads", "50").Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	fmt.Println(output)
}
