package tool

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func Gowitness(urls string) {
	args := "file"
	args2 := "-f"
	args3 := urls

	args4 := "--screenshot-path"
	args5 := "/tmp/yelaa/screenshots"
	out, err := exec.Command("gowitness", args, args2, args3, args4, args5).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}
	color.Yellow("Screenshot stored in /tmp/yelaa/screenshots")
	output := string(out[:])
	fmt.Println("gowitness output ", output)
}
