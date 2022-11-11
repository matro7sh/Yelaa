package tool

import (
	"github.com/fatih/color"
    dorks "github.com/bogdzn/gork/cmd"
)

type Dorks struct{
    outfile         string
    proxy           string
    userAgent       string
    extensions      []string
}

func (d *Dorks) Info(url string) {
	color.Cyan("Running dorks on %s", url)
}

func (d *Dorks) Configure(c interface{}) {

    /*
        gork will parse the DOM instead of making an API request, because it's easier for the end user
        (no API key to worry about etc), so we probably should **not** be changing the page's layout
        but, it's here in case something breaks
    */
    defaultUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36";
    defaultExtensions := []string{"doc", "docx", "csv", "pdf", "txt", "log", "bak", "json", "xlsx"}

    d.extensions = defaultExtensions
    d.userAgent = defaultUserAgent
    d.outfile = c.(map[string]interface{})["outfile"].(string)
    d.proxy = c.(map[string]interface{})["proxy"].(string)
}

func (d *Dorks) Run(url string) {
    opts := &dorks.Options{
        Outfile: d.outfile,
        AppendResults: true, /* we could be running this in a loop, should not erase former results */
        Proxy: d.proxy,
        Extensions: d.extensions,
        UserAgent: d.userAgent,
        Target: url,
    }

    dorks.Run(opts)
    color.Cyan("Dorks are stored in %s", d.outfile)
}
