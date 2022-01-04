package tool

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

func Gowitness(urls string) {
	currentTime := time.Now().Local()

	UserHomeDir, _ := os.UserHomeDir()
	args := "file"
	args2 := "-f"
	DomainsFile := urls

	args4 := "--screenshot-path"
	destinationPath := fmt.Sprintf("%s-%d-%d-%d_%d-%d", UserHomeDir+"/.yelaa/screenshots", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute())
	_, err := exec.Command("gowitness", args, args2, DomainsFile, args4, destinationPath).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	color.Yellow("Screenshot stored in " + destinationPath)
}
