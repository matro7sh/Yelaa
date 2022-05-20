package tool

import (
    "fmt"
    "sync"
    "time"

    "github.com/fatih/color"
    assetfinder "github.com/spiral-sec/assetfinder/scanner"

    "github.com/CMEPW/Yelaa/helper"
)


type fetchFn func(string) ([]string, error)

type AssetfinderConfig struct {
    subsOnly bool
    outfile string
    functions []fetchFn
}

type Assetfinder struct {
    opts *AssetfinderConfig
    outfile string
}

func (a *Assetfinder) Info(url string) {
    color.Cyan("Running Assetfinder on %s", url)
}

func (a *Assetfinder) Configure(c interface{}) {
    a.opts.subsOnly = false
    a.outfile = helper.YelaaPath + "/assetfinder.txt"
    a.opts.functions = []fetchFn{
        assetfinder.CertSpotter,
        assetfinder.HackerTarget,
        assetfinder.ThreatCrowd,
        assetfinder.CrtSh,
        assetfinder.Facebook,
        assetfinder.VirusTotal,
        assetfinder.FindSubDomains,
        assetfinder.Urlscan,
        assetfinder.BufferOverrun,
    }
}


func (a *Assetfinder) Run(url string) {
    var wg sync.WaitGroup

    rl := assetfinder.NewRateLimiter(time.Second)
    out := make(chan string)

    for _, f := range a.opts.functions {
        wg.Add(1)
        fn := f

        go func () {
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

    for n := range out {
        if _, ok := printed[n]; ok {
            continue
        }
        printed[n] = true

        fmt.Println(n)
    }
}

var _ ToolInterface = (*Assetfinder)(nil)
