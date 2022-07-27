package tool

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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
	g.optDir = gobusterdir.NewOptionsDir()
	g.optDir.StatusCodesBlacklistParsed.Add(404)
	g.optDir.NoTLSValidation = true
	g.optDir.Method = "GET"
	g.optDir.Timeout = time.Second * 10
	g.optDir.WildcardForced = true
	g.opts = &libgobuster.Options{Threads: 10, Wordlist: "yelaa.txt", OutputFilename: "scan_log_gobuster.txt"}
}

func (g *GoBuster) Run(website string) {
	g.optDir.URL = strings.TrimSuffix(website, "/")
	ctx := context.Background()
	data, err := os.Create("scan_data.csv")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer data.Close()
	data.WriteString("Url:;" + g.optDir.URL + "\nMethod:;" + g.optDir.Method +
		"\nThreads:;" + fmt.Sprintf("%d\n", g.opts.Threads) + "Wordlist:;" +
		g.opts.Wordlist + "\nOutput File name:;" + g.opts.OutputFilename +
		"\nTimeout:;" + g.optDir.Timeout.String() + "\n")
	d, _ := gobusterdir.NewGobusterDir(ctx, g.opts, g.optDir)
	e := cli.Gobuster(ctx, g.opts, d)
	if e != nil {
		fmt.Println(e)
	}
}

var _ ToolInterface = (*GoBuster)(nil)
