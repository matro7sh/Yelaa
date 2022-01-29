package tool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/CMEPW/Yelaa/helper"
)

func Assetfinder(url string) {
	out, err := exec.Command("assetfinder", url).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	destinationPath := helper.YelaaPath + "/assetfinder.txt"
	output := string(out[:])
	fmt.Println("assetfinder output ", output)
	os.WriteFile(destinationPath, []byte(output), 0644)
}
