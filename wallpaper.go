package main

import (
	"fmt"
	"net/http"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/reujab/wallpaper"
	"github.com/robfig/cron/v3"
)

var (
	c *cron.Cron
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("It can't be this easy")

	mReroll := systray.AddMenuItem("Retry", "I don't blame you")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Goodbye forever")
	mQuit.SetIcon(icon.Data)

	c = cron.New()

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mReroll.ClickedCh:
				SetWallpaper()
			}
		}
	}()

	// c.AddFunc("@midnight", SetWallpaper)

	SetWallpaper()
}

func onExit() {
	systray.Quit()
	c.Stop()
}

func SetWallpaper() {
	res, err := http.Get("https://source.unsplash.com/random/1920x1080")

	if err != nil {
		fmt.Print("error getting image from unplash")
	}

	wallpaper.SetFromURL(res.Request.URL.String())

	fmt.Println("setwallpaper finished")
}
