package tool

import (
	"fmt"
	"os"
	"os/exec"
)

func Subfinder(domain string, filename string) {
	out, err := exec.Command("subfinder", "-d", domain).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	err = os.WriteFile(filename, out, 0644)

	if err != nil {
		fmt.Printf("%s", err)
	}
}
