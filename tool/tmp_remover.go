package tool

import (
	"fmt"
	"os"
	"path/filepath"
)

func getAllPath() []string {
	files, err := filepath.Glob("/tmp/yelaa-*")
	if err != nil {
		fmt.Printf("%s", err)
	}
	return files
}

func TmpRemover() {
	allFiles := getAllPath()
	for _, file := range allFiles {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("%s", err)
		}
	}
	fmt.Print("All temporary file has been succesfully removed")
}
