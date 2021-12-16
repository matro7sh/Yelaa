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
	"time"

	"github.com/CMEPW/Yelaa/tool"
	"github.com/common-nighthawk/go-figure"
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
	domain        string
	insecure      bool
)

type folder struct {
	name     string
	children []folder
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func readFile() {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

	file, err := os.Open(targetPath)
	defer file.Close()

	if err != nil {
		fmt.Printf("%v, %+v", err, targetPath)
	}

	scanner := bufio.NewScanner(file)

	body, err := ioutil.ReadFile(targetPath)
	color.Magenta("Loaded target: \n%v", strings.Replace(string(body), "\n", ", ", -1))

	var websites []string

	for scanner.Scan() {
		// check if its ip/domain
		website := scanner.Text()
		if !strings.HasSuffix(website, "/") {
			website += "/"
		}
		websites = append(websites, website)

		color.Cyan("Looking for robots.txt on: %s", website)
		tool.GetRobot(website)

		color.Cyan("Looking for sitemap.xml on: %s ", website)
		tool.GetSitemap(website)

		color.Cyan("Running Dirsearch on %s", website)
		tool.Dirsearch(website)

		color.Cyan("Running Nuclei on %s", website)
		tool.Nuclei(website)

		color.Cyan("Looking for subdomains on %s", website)
		tool.GetSubdomains(website)
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

func checkProxy() {
	os.Setenv("HTTP_PROXY", proxy)
	os.Setenv("HTTPS_PROXY", proxy)

	if proxy != "" {
		color.Cyan("Proxy configuration: %s", proxy)
	} else {
		color.Cyan("No proxy has been set")
	}
}

func main() {

	version := figure.NewColorFigure("Yelaa 1.2.3", "", "cyan", true)
	version.Print()

	var cmdScan = &cobra.Command{
		Use:   "scan",
		Short: "It will run Nuclei templates, sslscan, dirsearch and more.",
		Long:  `We also make screenshot using gowitness and grap robots.txt, sitemaps.xml and gowitness.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			currentTime := time.Now()
			color.Cyan("Start scan: %v", currentTime.Format("2006-01-02 15:04:05"))
			checkProxy()
			readFile()
		},
	}

	var cmdOsint = &cobra.Command{
		Use:   "osint",
		Short: "Run subfinder, dnsx and httpx to find ips and subdomains of a specific domain",
		Long:  "First run subfinder on the domain to find all the subdomains, then pass the subdomains to dnsx to find all the ips and finally use httx against all the domains found",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			checkProxy()

			fmt.Printf("\nTarget domain: %s\n\n", domain)

			subdomainsFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
			if err != nil {
				fmt.Printf("%s", err)
			}

			ipsFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
			if err != nil {
				fmt.Printf("%s", err)
			}

			color.Cyan("Searching for subdomains with subfinder")
			tool.Subfinder(domain, subdomainsFile.Name())

			color.Cyan("Running dnsx on subdomains to find IP address")
			tool.Dnsx(subdomainsFile.Name(), ipsFile.Name())

			color.Cyan("Running httpx to find http servers")
			tool.Httpx(subdomainsFile.Name())

			subdomainsFile.Close()
			ipsFile.Close()
		},
	}

	var createDirectories = &cobra.Command{
		Use:   "create -c [client name]",
		Short: "It will create all directories to work",
		Long:  "Obtain a clean-cut architecture at the launch of a mission and make some tests",
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

	rootCmd.AddCommand(cmdScan, cmdOsint)
	rootCmd.Flags().StringVarP(&client, "client", "c", "", "Client name")
	rootCmd.Flags().StringVarP(&shared, "shared", "s", "", "path to shared folder")
	rootCmd.Flags().StringVarP(&excludedType, "excludedType", "e", "", "excluded type")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "Add HTTP proxy")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Allow insecure certificate")

	cmdScan.Flags().StringVarP(&targetPath, "target", "t", "", "Target file")

	cmdOsint.Flags().StringVarP(&domain, "domain", "d", "", "Target domain")

	if err := rootCmd.MarkFlagRequired("client"); err != nil {
		panic(err)
	}

	if err := cmdScan.MarkFlagRequired("target"); err != nil {
		panic(err)
	}

	if err := cmdOsint.MarkFlagRequired("domain"); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
