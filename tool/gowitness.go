package tool

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/remeh/sizedwaitgroup"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sensepost/gowitness/chrome"
	"github.com/sensepost/gowitness/lib"
	"github.com/sensepost/gowitness/storage"
)

var (
	chrm = chrome.NewChrome()
	db   = storage.NewDb()
)

type Gowitness struct {
	scanPath       string
	screenshotPath string
	file           string
	swg            sizedwaitgroup.SizedWaitGroup
    proxy          string
}

func (g *Gowitness) Info(_ string) {
	color.Cyan("Running gowitness on subdomains")
}

func (g *Gowitness) Configure(config interface{}) {
	chrm.ResolutionX = 1440
	chrm.ResolutionY = 900
	chrm.Delay = 0
	chrm.FullPage = false
	chrm.Timeout = 10
    chrm.Proxy = config.(map[string]interface{})["proxy"].(string)

	g.file = config.(map[string]interface{})["file"].(string)
	g.scanPath = config.(map[string]interface{})["scanPath"].(string)

	g.screenshotPath = g.scanPath + "/screenshots-" + time.Now().Format("2006-01-02_15-04-05")

	if _, err := os.Stat(g.screenshotPath); os.IsNotExist(err) {
		if err = os.Mkdir(g.screenshotPath, 0750); err != nil {
			fmt.Println(err)
		}
	}
	g.swg = sizedwaitgroup.New(4)
}

func (g *Gowitness) Run(_ string) {
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	f, err := os.Open(g.file)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	defer f.Close()

	db.Path = g.scanPath + "/gowitness.sqlite3"
	db, _ := db.Get()

	for scanner.Scan() {
		candidate := scanner.Text()
		if candidate == "" {
			return
		}

		for _, u := range getUrls(candidate) {
			g.swg.Add()

			go func(url *url.URL) {
				defer func() {
					g.swg.Done()
				}()

				p := &lib.Processor{
					Logger:         &l,
					Db:             db,
					Chrome:         chrm,
					URL:            url,
					ScreenshotPath: g.screenshotPath,
				}

				if err := p.Gowitness(); err != nil {
					log.Error().Err(err).Str("url", url.String()).Msg("failed to gowitness url")
				}
			}(u)
		}
	}

	g.swg.Wait()
	log.Info().Msg("processing complete")
}

func getUrls(target string) (c []*url.URL) {
	if strings.HasPrefix(target, "http") {
		u, err := url.Parse(target)
		if err == nil {
			c = append(c, u)
		}

		return
	}

	if !strings.HasPrefix(target, "http://") {
		u, err := url.Parse("http://" + target)
		if err == nil {
			c = append(c, u)
		}
	}

	if !strings.HasPrefix(target, "https://") {
		u, err := url.Parse("https://" + target)
		if err == nil {
			c = append(c, u)
		}
	}

	return
}
