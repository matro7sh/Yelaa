package tool

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	assetfinder "github.com/spiral-sec/assetfinder/scanner"
)

type fetchFn func(string) ([]string, error)

type AssetfinderConfig struct {
    scanPath  string
	outfile   string
	functions []fetchFn
}

type Assetfinder struct {
	opts      *AssetfinderConfig
	outfile   string
    functions []fetchFn
    scanPath  string
}

func (a *Assetfinder) Info(url string) {
	color.Cyan("Running Assetfinder on %s", url)
}

func (a *Assetfinder) Configure(c interface{}) {}

func (a *Assetfinder) Run(url string) {
	var wg sync.WaitGroup

	rl := assetfinder.NewRateLimiter(time.Second)
	out := make(chan string)

	for _, f := range a.opts.functions {
		wg.Add(1)
		fn := f

		go func() {
			defer wg.Done()

			rl.Block(fmt.Sprintf("%#v", fn))
			names, err := fn(url)

			if err != nil {
				return
			}

			for _, n := range names {
				n = assetfinder.CleanDomain(n)
				out <- n
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	printed := make(map[string]bool)
	file, err := os.OpenFile(a.opts.outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	for n := range out {
		if _, ok := printed[n]; ok {
			continue
		}
		printed[n] = true
		fmt.Println(n)
		file.WriteString(n + string('\n'))
	}
}

var _ ToolInterface = (*Assetfinder)(nil)
