package tool

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func Subfinder(domain string) {
	arg := "-d"
	out, err := exec.Command("subfinder", arg, domain).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	color.Green(output)
}
