package tool

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Httpx(ipsFile string) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	out, err := exec.Command("httpx", "-l", ipsFile, "-title", "-content-length", "-content-type",
		"-status-code", "-tech-detect", "-vhost", "-websocket", "-follow-redirects",
		"-ports", "25,80,81,135,389,443,1080,3000,3306,8080,8443,8888,9090,8089",
		"-retries", "2", "-timeout", "8", "-threads", "50", "-o", tempFile.Name()).Output()

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	output := string(out[:])
	fmt.Println(output)

	file, err := os.Open(tempFile.Name())
	if err != nil {
		fmt.Printf("%v, %+v", err, tempFile.Name())
		return
	}

	scanner := bufio.NewScanner(file)
	defer file.Close()

	domains := ""

	for scanner.Scan() {
		result := scanner.Text()
		domain := strings.Split(result, " ")[0]

		domains += domain + "\n"
	}

	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile(UserHomeDir+"checkAndScreen.txt", []byte(domains), 0644)
}
