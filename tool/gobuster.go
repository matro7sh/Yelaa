package tool

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CMEPW/Yelaa/helper"
	"github.com/OJ/gobuster/v3/cli"
	"github.com/OJ/gobuster/v3/gobusterdir"
	"github.com/OJ/gobuster/v3/libgobuster"
	"github.com/fatih/color"
)

type GoBuster struct {
	optDir   *gobusterdir.OptionsDir
	opts     *libgobuster.Options
	scanPath string
    proxy    string
}

func (s *GoBuster) Info(website string) {
	color.Cyan("Running gobuster on %s", website)
}

func (g *GoBuster) Configure(c interface{}) {

	wordlist := c.(map[string]interface{})["wordlist"].(string)
	g.scanPath = c.(map[string]interface{})["scanPath"].(string)
	outputDir := g.scanPath + "/gobuster"

	blacklist := libgobuster.NewIntSet()
	blacklist.Add(302)
	blacklist.Add(404)
	blacklist.Add(500)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0750); err != nil {
			fmt.Println(err)
		}
	}

    g.proxy = c.(map[string]interface{})["proxy"].(string)
    g.optDir.UserAgent = helper.GetUserAgent()
	g.optDir = gobusterdir.NewOptionsDir()
	g.optDir.StatusCodesBlacklistParsed.Add(404)
	g.optDir.NoTLSValidation = true
	g.optDir.Method = "GET"
	g.optDir.Timeout = time.Second * 10
	g.optDir.WildcardForced = true
	g.optDir.StatusCodesBlacklistParsed = blacklist
	g.opts = &libgobuster.Options{
        Threads: 10,
        Wordlist: wordlist,
    }
}

func (g *GoBuster) Run(website string) {
	g.opts.OutputFilename = g.scanPath + "/gobuster/scan_log_gobuster-" +
		time.Now().Format("2006-01-02_15-04-05") + ".txt"
	g.optDir.URL = strings.TrimSuffix(website, "/")
	ctx := context.Background()
	d, _ := gobusterdir.NewGobusterDir(ctx, g.opts, g.optDir)
	e := cli.Gobuster(ctx, g.opts, d)
	if e != nil {
		fmt.Println(e)
	}
	CsvWriterGobuster(g)
}

var _ ToolInterface = (*GoBuster)(nil)
