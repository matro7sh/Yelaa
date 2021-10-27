package tool

import (
	"fmt"
	"os/exec"
)

func Sslscan(url string) {
	args := "-u"
	args2 := url
	out, err := exec.Command("testssl", args, args2).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	fmt.Println("testssl output ", output)
}
