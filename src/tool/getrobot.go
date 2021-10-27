package tool

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

func GetRobot(url string) {
	resp, err := http.Get(url + "robots.txt")
	if err != nil {
		fmt.Printf("%v", err)
	}

	if resp != nil && resp.StatusCode != http.StatusNotFound {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v", err)
		}
		sb := string(body)
		color.Green(sb)

	} else {
		color.Yellow("----- Sorry get 404 status code for this robots.txt ----- ")
	}
}
