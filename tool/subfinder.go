package tool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func Subfinder(domain string, filename string) {
	arg := "-d"
	out, err := exec.Command("subfinder", arg, domain).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	color.Green(output)

	err = os.WriteFile(filename, out, 0644)

	if err != nil {
		fmt.Printf("%s", err)
	}
}
