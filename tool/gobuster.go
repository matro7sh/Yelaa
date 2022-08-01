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
	optDir *gobusterdir.OptionsDir
	opts   *libgobuster.Options
}

func (s *GoBuster) Info(website string) {
	color.Cyan("Running gobuster on %s", website)
}

func (g *GoBuster) Configure(c interface{}) {
	outputDir := helper.YelaaPath + "/gobuster"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.Mkdir(outputDir, 0750); err != nil {
			fmt.Println(err)
		}
	}
	g.optDir = gobusterdir.NewOptionsDir()
	g.optDir.StatusCodesBlacklistParsed.Add(404)
	g.optDir.NoTLSValidation = true
	g.optDir.Method = "GET"
	g.optDir.Timeout = time.Second * 10
	g.optDir.WildcardForced = true
	g.opts = &libgobuster.Options{Threads: 10, Wordlist: "yelaa.txt"}
}

func (g *GoBuster) Run(website string) {
	g.opts.OutputFilename = helper.YelaaPath + "/gobuster/scan_log_gobuster-" +
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
