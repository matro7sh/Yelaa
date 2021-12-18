package tool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func Gowitness(urls string) {
	UserHomeDir, _ := os.UserHomeDir()
	args := "file"
	args2 := "-f"
	DomainsFile := urls

	args4 := "--screenshot-path"
	destinationPath := UserHomeDir + "/.yelaa/screenshots"
	_, err := exec.Command("gowitness", args, args2, DomainsFile, args4, destinationPath).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}
	color.Yellow("Screenshot stored in ~/.yelaa/screenshots")

}
