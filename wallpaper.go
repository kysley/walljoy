package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/reujab/wallpaper"
	"github.com/robfig/cron/v3"
)

type Kv struct {
	Value string
}

var (
	c *cron.Cron
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		if runtime.GOOS == "windows" {
			watchWindows()
		} else if runtime.GOOS == "darwin" {
			watchOSX()
		}
	}()

	systray.SetIcon(getIcon("smile_light.ico"))
	systray.SetTitle("Walljoy")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit Walljoy", ":wave:")
	mRefresh := systray.AddMenuItem("Refresh", "")

	c = cron.New()

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mRefresh.ClickedCh:
				setWallpaper()
			}
		}
	}()

	setWallpaper()

	c.AddFunc("@midnight", setWallpaper)
	c.Start()
}

func onExit() {
	systray.Quit()
	c.Stop()
}

func setWallpaper() {
	res, err := http.Get("https://source.unsplash.com/random/1920x1080")

	if err != nil {
		fmt.Print("error getting image from unplash")
	}

	wallpaper.SetFromURL(res.Request.URL.String())
}
