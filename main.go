package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	baseDirectory string
	client string
	excludedType string
	path string
)

type folder struct {
	name     string
	children []folder
}

func createDirectory(base string, folders []folder) {
	for _, f := range folders {
		if f.name == excludedType {
			continue
		}
		current := fmt.Sprintf("%s/%s", base, f.name)
		if e := os.Mkdir(current, 0775); e != nil {
			fmt.Println(e)
		}
		if len(f.children) != 0 {
			createDirectory(current, f.children)
		}
	}
}

func copyCherryTreeAndTargets() {
	if e := os.Link("./trace.ctb", baseDirectory+"/Web-Penetration-Test/trace.ctb"); e != nil {
		fmt.Println(e)
	}
	if e := os.Link("./targets.txt", baseDirectory+"/targets.txt"); e != nil {
		fmt.Println(e)
	}
}

func folderNameFactory(names ...string) []folder {
	f := make([]folder, len(names))
	for _, name := range names {
		f = append(f, folder{name: name})
	}

	return f
}

func main() {

	var createDirectories = &cobra.Command{
		Use:   "create -c [client name]",
		Short: "It will create all directories to work",
		Long:  "Obtain a clean-cut architecture at the launch of a mission",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("User input: " + client)

			if path == "" {
				path = "."
			}
			baseDirectory = fmt.Sprintf("%s/%s", path, client)
			_ = os.MkdirAll(baseDirectory, 0775)
			createDirectory(baseDirectory, []folder{
				{
					name: "Infrastructure-Penetration-Test",
				},
				{
					name: "Web-Penetration-Test",
					children: folderNameFactory("nmap", "nessus", "report", "screenshot", "ssl"),
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
	rootCmd.Flags().StringVarP(&excludedType, "excludedType", "e", "", "EXCLUDED TYPE")
	if err := rootCmd.MarkFlagRequired("client"); err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
