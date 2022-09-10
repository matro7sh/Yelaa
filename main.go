package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/CMEPW/Yelaa/helper"
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
	scanPath      string
	targetPath    string
	proxy         string
	domain        string
	insecure      bool
	dryRun        bool
)

type folder struct {
	name     string
	children []folder
}

type FileScanner struct {
	io.Closer
	*bufio.Scanner
}

func loadTargetFile() *FileScanner {
	file, err := os.Open(targetPath)
	if err != nil {
		fmt.Printf("%v, %+v", err, targetPath)
	}

	scanner := bufio.NewScanner(file)

	body, err := ioutil.ReadFile(targetPath)

	if err != nil {
		fmt.Printf("%v, %+v", err, targetPath)
	}

	color.Magenta("Loaded target: \n%v", strings.Replace(string(body), "\n", ", ", -1))

	return &FileScanner{file, scanner}
}

func readFile() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

	var toolList []tool.ToolInterface
	toolList = append(toolList, &tool.Robot{}, &tool.Sitemap{}, &tool.Nuclei{})

	gb := tool.GoBuster{}
	gbCfg := make(map[string]interface{})
	gbCfg["scanPath"] = scanPath

	gb.Configure(gbCfg)

	var i interface{}
	for _, t := range toolList {
		t.Configure(i)
	}

	scanner := loadTargetFile()
	for scanner.Scan() {
		// check if its ip/domain
		website := scanner.Text()
		defer scanner.Close()
		if !strings.HasSuffix(website, "/") {
			website += "/"
		}

		color.Cyan("Running tools on %s", website)
		for _, t := range toolList {
			if t == toolList[len(toolList)-1] {
				break
			}

			if !dryRun {
				t.Run(website)
			}
		}

		if !dryRun {
			gb.Run(website)
		}
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

func createOutDirectory() {
	if _, err := os.Stat(scanPath); os.IsNotExist(err) {
		color.Cyan("Creating " + scanPath + " folder")
		if err = os.MkdirAll(scanPath, 0755); err != nil {
			fmt.Println(err)
		}
	}
}

func scanDomain(domain string) {
	fmt.Printf("\nTarget domain for this loop: %s\n\n", domain)
	dorks := tool.Dorks{}
	dorksCfg := make(map[string]interface{})
	dorks.Configure(dorksCfg)
	dorks.Info(domain)
	dorks.Run("https://" + domain)

	subdomainsFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
	if err != nil {
		fmt.Printf("%s", err)
	}

	ipsFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
	if err != nil {
		fmt.Printf("%s", err)
	}

	sf := tool.Subfinder{}
	configuration := make(map[string]interface{})
	configuration["filename"] = subdomainsFile.Name()
	sf.Info("")
	sf.Configure(configuration)

	if !dryRun {
		sf.Run(domain)
	}

	asf := tool.Assetfinder{}
	asfCfg := make(map[string]interface{})

	asfOutfile := scanPath + "/assetfinder.txt"

	asfCfg["scanPath"] = scanPath
	asfCfg["outfile"] = asfOutfile

	asf.Configure(asfCfg)
	asf.Info(domain)

	if !dryRun {
		asf.Run(domain)
	}

	dnsx := tool.Dnsx{}
	dnsxConfig := make(map[string]interface{})
	dnsxConfig["subdomainsFilename"] = subdomainsFile.Name()
	dnsxConfig["ipsFilename"] = ipsFile.Name()
	dnsx.Info("")
	dnsx.Configure(dnsxConfig)

	if !dryRun {
		dnsx.Run("")
	}

	domainsFiles := []string{asfOutfile, subdomainsFile.Name(), ipsFile.Name()}
	var domainBuffer bytes.Buffer

	for _, file := range domainsFiles {
		newDomain, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("%s", err)
		}

		domainBuffer.Write(newDomain)
		// check if directory already exist + name of projet

		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(scanPath+"/domains.txt", domainBuffer.Bytes(), 0644)
		if err != nil {
			fmt.Printf("%s", err)
		}
	}

	filepath := scanPath + "/osint.domains.txt"
	httpx := tool.Httpx{}
	httpxConfig := make(map[string]interface{})
	httpxConfig["input"] = scanPath + "/domains.txt"
	httpxConfig["output"] = filepath
	httpx.Info("")
	httpx.Configure(httpxConfig)

	if !dryRun {
		httpx.Run("")
	}

	gw := tool.Gowitness{}
	gwConfig := make(map[string]interface{})
	gwConfig["file"] = filepath
	gwConfig["scanPath"] = scanPath

	gw.Info("")
	gw.Configure(gwConfig)

	if !dryRun {
		gw.Run("")
	}

	subdomainsFile.Close()
	ipsFile.Close()
}

func main() {
	version := figure.NewColorFigure("Yelaa 1.5.5", "", "cyan", true)
	version.Print()

	var cmdScan = &cobra.Command{
		Use:   "scan",
		Short: "It will run Nuclei templates, dirsearch and more.",
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
			createOutDirectory()
			checkProxy()

			if targetPath == "" {
				scanDomain(domain)
				return
			}

			scanner := loadTargetFile()
			defer scanner.Close()

			for scanner.Scan() {
				targetDomain := scanner.Text()
				dorks := tool.Dorks{}
				dorksCfg := make(map[string]interface{})
				dorks.Configure(dorksCfg)
				dorks.Info(targetDomain)
				dorks.Run(targetDomain)
			}

		},
	}

	var checkAndScreen = &cobra.Command{
		Use:   "checkAndScreen -t list_of_ip.txt",
		Short: "Run httpx and gowitness",
		Long:  "Run httpx on each IP and take screenshots of each server that are up",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createOutDirectory()
			checkProxy()

			if targetPath == "" {
				color.Red("Please provide a list of ips/domains")
				return
			}

			filepath := scanPath + "/checkAndScreen.txt"
			httpx := tool.Httpx{}
			httpxConfig := make(map[string]interface{})
			httpxConfig["input"] = targetPath
			httpxConfig["output"] = filepath
			httpx.Info("")
			httpx.Configure(httpxConfig)

			if !dryRun {
				httpx.Run("")
			}

			gw := tool.Gowitness{}
			gwConfig := make(map[string]interface{})
			gwConfig["file"] = filepath
			gw.Info("")
			gw.Configure(gwConfig)

			if !dryRun {
				gw.Run("")
			}
		},
	}

	var createDirectories = &cobra.Command{
		Use:   "create -c [client name]",
		Short: "It will create all directories to work",
		Long:  "Obtain a clean-cut architecture at the launch of a mission and make some tests",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createOutDirectory()
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

	rootCmd.AddCommand(cmdScan, cmdOsint, checkAndScreen)

	rootCmd.Flags().StringVarP(&client, "client", "c", "", "Client name")
	rootCmd.Flags().StringVarP(&shared, "shared", "s", "", "path to shared folder")
	rootCmd.Flags().StringVarP(&excludedType, "excludedType", "e", "", "excluded type")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "Add HTTP proxy")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Allow insecure certificate")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Run in dry-run mode")
	rootCmd.PersistentFlags().StringVar(&scanPath, "path", helper.YelaaPath, "Output path")

	cmdScan.Flags().StringVarP(&targetPath, "target", "t", "", "Target file")

	cmdOsint.Flags().StringVarP(&domain, "domain", "d", "", "Target domain")
	cmdOsint.Flags().StringVarP(&targetPath, "target", "t", "", "Target domains file")

	checkAndScreen.Flags().StringVarP(&targetPath, "target", "t", "", "list of ips/domains")

	if err := rootCmd.MarkFlagRequired("client"); err != nil {
		panic(err)
	}

	if err := cmdScan.MarkFlagRequired("target"); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
	tool.TmpRemover()
}
