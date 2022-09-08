package helper

import "os"

var YelaaPath = GetHome() + "/.yelaa"

func GetHome() (home string) {
	home, _ = os.UserHomeDir()
	return
}
