package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/getlantern/systray"
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
	go func() {
		if runtime.GOOS == "windows" {
			watchWindows()
		} else if runtime.GOOS == "darwin" {
			watchOSX()
		}
	}()

	systray.SetIcon(getIcon("smile_light.ico"))
	systray.SetTitle("It can't be this easy")

	mReroll := systray.AddMenuItem("Shuffle Wallpaper", "Gets a new random wallpaper from Unsplash")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit Walljoy", "Goodbye")

	c = cron.New()

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mReroll.ClickedCh:
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

	fmt.Println("setwallpaper finished")
}
