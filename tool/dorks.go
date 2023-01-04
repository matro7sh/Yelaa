package tool

import (
	dorks "github.com/bogdzn/gork/cmd"
	"github.com/fatih/color"
)

type Dorks struct {
	outfile string
	proxy   string
}

func (d *Dorks) Info(url string) {
	color.Cyan("Running dorks on %s", url)
}

func (d *Dorks) Configure(c interface{}) {
	d.outfile = c.(map[string]interface{})["outfile"].(string)
	d.proxy = c.(map[string]interface{})["proxy"].(string)
}

func (d *Dorks) Run(domain string) {
	opts := &dorks.Options{
		Proxy:         d.proxy,
		Outfile:       d.outfile,
		AppendResults: false,
		Extensions:    dorks.DefaultFileExtensions(),
		Exclusions:    dorks.DefaultExclusions(),
		UserAgent:     dorks.DefaultUserAgent(),
		Target:        domain,
	}

	dorks.Run(opts)
	color.Cyan("Dorks are stored in %s", d.outfile)
}
