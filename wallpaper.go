package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	"github.com/getlantern/systray"
	"github.com/go-resty/resty/v2"
	"github.com/reujab/wallpaper"
	"github.com/robfig/cron/v3"
)

var (
	c      *cron.Cron
	client = resty.New()
)

func main() {
	systray.Run(onReady, onExit)
}

// [id, name, code]
func getDeviceIdentity() string {
	deviceName, _ := os.Hostname()
	deviceId, err := machineid.ProtectedID("walljoy")

	if err != nil {
		fmt.Print("Error getting device identity")
	}

	return deviceId + "," + deviceName
}

type ackResponse struct {
	SessionId string `json:"sessionId"`
	Code      string `json:"newCode"`
}

func getSessionId() {
	resp, err := client.R().SetBody(map[string]interface{}{"identity": getDeviceIdentity()}).Post("http://localhost:8081/ack")

	if err != nil {
		// writeFile(string(resp.Body()))
		fmt.Print(resp.Body())
	}
	var sessionData ackResponse

	if err := json.Unmarshal(resp.Body(), &sessionData); err != nil {
		fmt.Print("Client unmarshal failed: " + err.Error())
	}

	fmt.Print("hello")

	sessionId := sessionData.SessionId

	openBrowser("http://localhost:8080/register?sId=" + sessionId)
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
	getSessionId()

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

	// browser.OpenURL()
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
