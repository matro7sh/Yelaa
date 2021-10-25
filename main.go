package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	baseDirectory string
	client        string
	excludedType  string
	shared        string
	targetPath    string
	proxy         string
	insecure      bool
)

var GlobalHeaders = []string{"Server", "X-XSS-Protection", "Access-Control-Allow-Credentials", "Content-Security-Policy", "X-Powered-By", "Strict-Transport-Security"}

type folder struct {
	name     string
	children []folder
}

func getRobot(url string) {
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

func getSitemap(url string) {
	resp, err := http.Get(url + "sitemap.xml")

	if err != nil {
		fmt.Printf("%v", err)
	}

	if resp != nil && resp.StatusCode != http.StatusNotFound {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v", err)
		}

		for headerName, headerValue := range resp.Header {
			if contains(GlobalHeaders, headerName) {
				//		fmt.Println("Header + " + headerName + "Found : " + headerValue)
				fmt.Printf("Found Header: %s | %s \n", headerName, headerValue)
			} else {
				// fmt.Println("sorry no headers in global headers variable")
			}
		}
		sb := string(body)
		color.Green(sb)

	} else {
		color.Yellow("-----  Sorry get 404 status code for this sitemap.xml -----")
	}
}

func readFile() {
	// add check to / at the end using regex or something and check for domain/CIDR

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

	file, err := os.Open(targetPath)
	// color.Cyan("Loaded targets :")
	defer file.Close()

	if err != nil {
		fmt.Printf("%v, %+v", err, targetPath)
	}

	scanner := bufio.NewScanner(file)
	var websites []string

	for scanner.Scan() {
		// check if its ip/domain
		website := scanner.Text()
		if !strings.HasSuffix(website, "/") {
			website += "/"
		}
		websites = append(websites, website)
		color.Cyan("Looking for robots.txt on: %s", scanner.Text())
		getRobot(website)
		color.Cyan("Looking for sitemap.xml on: %s ", scanner.Text())
		getSitemap(website)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("%v \n", err)
	}

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

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func main() {

	var cmdScan = &cobra.Command{
		Use:   "scan",
		Short: "It will run Nuclei templates, sslscan, dirsearch and more.",
		Long:  `We also make screenshot using gowitness and grap robots.txt, sitemaps.xml and gowitness.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Start scan: " + strings.Join(args, " "))
			os.Setenv("HTTP_PROXY", proxy)
			os.Setenv("HTTPS_PROXY", proxy)
			if proxy != "" {
				fmt.Println("Proxy configuration: ", proxy)
			} else {
				fmt.Println("No proxy has been set")
			}

			fmt.Println("Loading file: ")
			readFile()
		},
	}

	var createDirectories = &cobra.Command{
		Use:   "create -c [client name]",
		Short: "It will create all directories to work",
		Long:  "Obtain a clean-cut architecture at the launch of a mission",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Setup mission for: ", client)

			if shared == "" {
				shared = "."
			}
			baseDirectory = fmt.Sprintf("%s/%s", shared, client)
			_ = os.MkdirAll(baseDirectory, 0775)
			createDirectory(baseDirectory, []folder{
				{
					name: "Infrastructure-Penetration-Test",
				},
				{
					name:     "Web-Penetration-Test",
					children: folderNameFactory("nmap", "nessus", "report", "screenshot", "ssl"),
				},
			})
			copyCherryTreeAndTargets()
			out, _ := exec.Command("tree", baseDirectory).Output()
			fmt.Println(string(out))
		},
	}

	var rootCmd = createDirectories

	rootCmd.AddCommand(cmdScan)
	rootCmd.Flags().StringVarP(&client, "client", "c", "", "Client name")
	cmdScan.Flags().StringVarP(&targetPath, "target", "t", "", "Target file")
	cmdScan.Flags().StringVarP(&proxy, "proxy", "p", "", "Add HTTP proxy")
	cmdScan.Flags().BoolVarP(&insecure, "insecure", "k", false, "Allow insecure certificate")
	rootCmd.Flags().StringVarP(&shared, "shared", "s", "", "path to shared folder")
	rootCmd.Flags().StringVarP(&excludedType, "excludedType", "e", "", "excluded type")

	if err := rootCmd.MarkFlagRequired("client"); err != nil {
		panic(err)
	}

	if err := cmdScan.MarkFlagRequired("target"); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
