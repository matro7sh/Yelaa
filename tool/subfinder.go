package tool

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Subfinder(domain string, filename string) {

	var myrealString = strings.TrimPrefix(domain, "https://")
	if domain == "" {
		myrealString = domain
	}
	out, err := exec.Command("subfinder", "-d", myrealString).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	err = os.WriteFile(filename, out, 0644)

	if err != nil {
		fmt.Printf("%s", err)
	}

}
