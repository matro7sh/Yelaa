package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type Folder struct {
	name     string
	children []Folder
}

var client string
var exludedType string
var path string
var baseDirectory string

func createDirectory(base string, folders []Folder) {
	for _, folder := range folders {
		if folder.name == exludedType {
			continue
		}
		current := fmt.Sprintf("%s/%s", base, folder.name)
		if e := os.Mkdir(current, 0775); e != nil {
			fmt.Println(e)
		}
		if len(folder.children) != 0 {
			createDirectory(current, folder.children)
		}
	}
}

func copyCherryTreeAndTargets() {
	os.Link("./trace.ctb", baseDirectory+"/Web-Penetration-Test/trace.ctb")
	os.Link("./targets.txt", baseDirectory+"/targets.txt")
}

func main() {

	var createDirectories = &cobra.Command{
		Use:   "create [string to print]",
		Short: "It will create all directories to work",
		Long:  `Obtain a clean-cut architecture at the launch of a mission`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("User input: " + client)

			if path == "" {
				path = "."
			}
			baseDirectory = fmt.Sprintf("%s/%s", path, client)
			_ = os.MkdirAll(baseDirectory, 0775)
			createDirectory(baseDirectory, []Folder{
				Folder{
					name: "Infrastructure-Penetration-Test",
				},
				Folder{
					name: "Web-Penetration-Test",
					children: []Folder{
						Folder{
							name: "nmap",
						},
						Folder{
							name: "nessus",
						},
						Folder{
							name: "report",
						},
						Folder{
							name: "screenshot",
						},
						Folder{
							name: "ssl",
						},
					},
				},
			})
			copyCherryTreeAndTargets()
			out, _ := exec.Command("tree", baseDirectory).Output()
			fmt.Println(string(out))
		},
	}

	var rootCmd = createDirectories
	rootCmd.Flags().StringVarP(&client, "client", "c", "", "Client name")
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "path to shared folder")
	rootCmd.Flags().StringVarP(&exludedType, "excludedType", "e", "", "EXCLUDED TYPE")
	rootCmd.MarkFlagRequired("client")
	rootCmd.Execute()

}
