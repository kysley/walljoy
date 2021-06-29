package main

import (
	"encoding/json"
	"fmt"
	"io"
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

type Kv struct {
	Value int
}

var (
	c        *cron.Cron
	client   = resty.New()
	store, _ = Open("./walljoy.db")
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

	// sessionId := sessionData.SessionId

	// openBrowser("http://localhost:8080/register?sId=" + sessionId)
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

	// mReroll := systray.AddMenuItem("Shuffle Wallpaper", "Gets a new random wallpaper from Unsplash")
	systray.AddMenuItem("Current: ", "").Disable()
	systray.AddSeparator()
	mChan1 := systray.AddMenuItem("Earth", "")
	mChan2 := systray.AddMenuItem("Structure", "")
	mChan3 := systray.AddMenuItem("Random", "")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit Walljoy", "Goodbye")

	c = cron.New()

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mChan3.ClickedCh:
				onCollectionChange(3)
			case <-mChan2.ClickedCh:
				onCollectionChange(2)
			case <-mChan1.ClickedCh:
				onCollectionChange(1)
			}
		}
	}()

	getSessionId()

	err := store.Get("collectionId", Kv{})

	if err != nil {
		collectionId := Kv{1}
		store.Set("collectionId", collectionId)
		fmt.Print(collectionId.Value)
	}

	setWallpaper()

	c.AddFunc("@midnight", setWallpaper)
	c.Start()
}

func onExit() {
	systray.Quit()
	c.Stop()
}

func setWallpaper() {
	collectionId := Kv{}
	if err := store.Get("collectionId", &collectionId); err != nil {
		print(err)
	}
	res, err := http.Get(fmt.Sprintf("http://localhost:8081/c/%d", collectionId.Value))

	if err != nil {
		fmt.Print("error getting image from api")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		print("malformed body")
	}
	wallpaper.SetFromURL(string(body))

	fmt.Println("setwallpaper finished")
}

func onCollectionChange(cId int) {
	store.Set("collectionId", Kv{cId})
	setWallpaper()
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
