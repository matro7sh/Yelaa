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
	toolList = append(toolList, &tool.Robot{}, &tool.Sitemap{}, &tool.GoBuster{}, &tool.Nuclei{})

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
			t.Run(website)
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

func scanDomain(domain string) {
	fmt.Printf("\nTarget domain for this loop: %s\n\n", domain)

	subdomainsFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
	if err != nil {
		fmt.Printf("%s", err)
	}

	ipsFile, err := ioutil.TempFile(os.TempDir(), "yelaa-")
	if err != nil {
		fmt.Printf("%s", err)
	}

	getSubDomainCrt, err := ioutil.TempFile(os.TempDir(), "yelaa-")
	if err != nil {
		fmt.Printf("%s", err)
	}

	sf := tool.Subfinder{}
	configuration := make(map[string]interface{})
	configuration["filename"] = subdomainsFile.Name()
	sf.Info("")
	sf.Configure(configuration)
	sf.Run(domain)

	color.Cyan("Make request to crt.sh on domain")
	tool.Crt(domain, getSubDomainCrt.Name())

	color.Cyan("Running dnsx on subdomains to find IP address")
	tool.Dnsx(subdomainsFile.Name(), ipsFile.Name())

	domainsFiles := []string{subdomainsFile.Name(), getSubDomainCrt.Name(), ipsFile.Name()}
	var domainBuffer bytes.Buffer
	UserHomeDir, err := os.UserHomeDir()

	if _, err := os.Stat(UserHomeDir + "/.yelaa"); os.IsNotExist(err) {
		fmt.Println("take care folder already exist")
	}
	err = os.Mkdir(UserHomeDir+"/.yelaa", 0755)

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

		err = ioutil.WriteFile(UserHomeDir+"/.yelaa/domains.txt", domainBuffer.Bytes(), 0644)
		if err != nil {
			fmt.Printf("%s", err)
		}
	}

	color.Cyan("Running httpx to find http servers")
	tool.Httpx(UserHomeDir + "/.yelaa/domains.txt")

	color.Cyan("Running Gowitness")
	tool.Gowitness(UserHomeDir + "/.yelaa/domains.txt")

	subdomainsFile.Close()
	ipsFile.Close()
	getSubDomainCrt.Close()
}

/*
func NewGobuster(opts *Options, plugin GobusterPlugin) (*Gobuster, error) {
	var g Gobuster
	g.Opts = opts
	g.plugin = plugin
	g.RequestsCountMutex = new(sync.RWMutex)
	g.resultChan = make(chan Result)
	g.errorChan = make(chan error)
	g.LogInfo = log.New(os.Stdout, "", log.LstdFlags)
	g.LogError = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	return &g, nil
}
*/

func main() {
	version := figure.NewColorFigure("Yelaa 1.3.2", "", "cyan", true)
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

			if targetPath == "" {
				scanDomain(domain)
				return
			}

			scanner := loadTargetFile()
			defer scanner.Close()

			for scanner.Scan() {
				targetDomain := scanner.Text()
				scanDomain(targetDomain)
			}
		},
	}

	var checkAndScreen = &cobra.Command{
		Use:   "checkAndScreen -l list_of_ip.txt",
		Short: "Run httpx and gowitness",
		Long:  "Run httpx on each IP and take screenshots of each server that are up",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkProxy()

			if targetPath == "" {
				color.Red("Please provide a list of ips/domains")
				return
			}

			color.Cyan("Running httpx to find http servers")
			tool.Httpx(targetPath)

			UserHomeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
			}

			color.Cyan("Running gowitness on server found by httpx")
			tool.Gowitness(UserHomeDir + "/.yelaa/checkAndScreen.txt")
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

	rootCmd.AddCommand(cmdScan, cmdOsint, checkAndScreen)
	rootCmd.Flags().StringVarP(&client, "client", "c", "", "Client name")
	rootCmd.Flags().StringVarP(&shared, "shared", "s", "", "path to shared folder")
	rootCmd.Flags().StringVarP(&excludedType, "excludedType", "e", "", "excluded type")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "Add HTTP proxy")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "Allow insecure certificate")

	cmdScan.Flags().StringVarP(&targetPath, "target", "t", "", "Target file")

	cmdOsint.Flags().StringVarP(&domain, "domain", "d", "", "Target domain")
	cmdOsint.Flags().StringVarP(&targetPath, "target", "t", "", "Target domains file")

	checkAndScreen.Flags().StringVarP(&targetPath, "list", "l", "", "list of ips/domains")

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
