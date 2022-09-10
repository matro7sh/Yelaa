package tool

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

type Dorks struct{}

func (d *Dorks) Info(url string) {
	color.Cyan("Use dorks on %s", url)
}

func (d *Dorks) Configure(c interface{}) {}

func (d *Dorks) Run(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", "https://www.google.com/search?q=site:"+url+"+ext:doc+OR+ext:docx+OR+ext:csv+OR+ext:pdf+OR+ext:txt+OR+ext:log+OR+ext:bak").Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
