package tool

import (
	"context"
	"fmt"

	"github.com/OJ/gobuster/v3/cli"
	"github.com/OJ/gobuster/v3/gobusterdir"
	"github.com/OJ/gobuster/v3/libgobuster"
)

type GoBuster struct {
	currentURL string
}

func (*GoBuster) Configure(c interface{}) {}

func (g *GoBuster) Run(website string) {
	optDir := gobusterdir.NewOptionsDir()
	optDir.URL = g.currentURL
	optDir.StatusCodesBlacklistParsed.Add(404)

	options := &libgobuster.Options{Threads: 2, Wordlist: "yelaa.txt"}
	ctx := context.Background()

	d, _ := gobusterdir.NewGobusterDir(ctx, options, optDir)
	e := cli.Gobuster(ctx, options, d)
	if e != nil {
		fmt.Println(e)
	}
}

var _ ToolInterface = (*GoBuster)(nil)
