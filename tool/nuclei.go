package tool

import (
	"fmt"
	"os/exec"
)

func Nuclei(url string) {
	args := "-u"
	args2 := url
	out, err := exec.Command("nuclei", args, args2).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	fmt.Println("nuclei output ", output)
}
