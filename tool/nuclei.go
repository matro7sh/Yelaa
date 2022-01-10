package tool

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

type Nuclei struct{}

func (*Nuclei) Configure(c interface{}) {}

func (*Nuclei) Info(website string) {
	color.Cyan("Running Nuclei on %s", website)
}

func (*Nuclei) Run(website string) {
	args := "-u"
	args2 := website
	out, err := exec.Command("nuclei", args, args2).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	fmt.Println("nuclei output ", output)
}
