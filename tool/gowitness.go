package tool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func Gowitness(urls string) {
	UserHomeDir, err := os.UserHomeDir()
	args := "file"
	args2 := "-f"
	args3 := urls

	args4 := "--screenshot-path"
	args5 := UserHomeDir + "/.yelaa/screenshots"
	out, err := exec.Command("gowitness", args, args2, args3, args4, args5).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}
	color.Yellow("Screenshot stored in ~/.yelaa/screenshots")
	output := string(out[:])
	fmt.Println("gowitness output ", output)
}
