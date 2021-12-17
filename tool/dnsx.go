package tool

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Dnsx(subdomainsFile string, ipsFile string) {

	out, err := exec.Command("dnsx", "-l", subdomainsFile, "-resp", "-a").Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	bytesReader := bytes.NewReader(out)
	scanner := bufio.NewScanner(bytesReader)

	ip_list := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		color.Green(line)

		line = strings.Replace(line, "[", "", -1)
		line = strings.Replace(line, "]", "", -1)
		ip := strings.Split(line, " ")[1]

		if contains(ip_list, ip) {
			continue
		}

		ip_list = append(ip_list, ip)

	}

	output := strings.Join(ip_list, "\n")

	os.WriteFile(ipsFile, []byte(output), 0644)

	if err != nil {
		fmt.Printf("%s", err)
	}
}
