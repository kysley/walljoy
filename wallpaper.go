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
	"strconv"

	"github.com/denisbrodbeck/machineid"
	"github.com/getlantern/systray"
	"github.com/go-resty/resty/v2"
	"github.com/reujab/wallpaper"
	"github.com/robfig/cron/v3"
)

type Kv struct {
	Value string
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
		fmt.Print(resp.Body())
	}
	var sessionData ackResponse

	if err := json.Unmarshal(resp.Body(), &sessionData); err != nil {
		fmt.Print("Client unmarshal failed: " + err.Error())
	}

	fmt.Print("hello")

	store.Set("sId", Kv{sessionData.SessionId})
}

func openSession() {
	var sessionId Kv

	store.Get("sId", &sessionId)
	openBrowser(fmt.Sprintf("http://localhost:8080/register?sId=%s", sessionId.Value))
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

	mDash := systray.AddMenuItem("Dashboard", "")
	systray.AddSeparator()
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
			case <-mDash.ClickedCh:
				openSession()
			}
		}
	}()

	getSessionId()

	err := store.Get("collectionId", Kv{})

	if err != nil {
		collectionId := Kv{"1"}
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
	res, err := http.Get(fmt.Sprintf("http://localhost:8081/c/%s", collectionId.Value))

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
	stringCId := strconv.Itoa(cId)
	store.Set("collectionId", Kv{stringCId})
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
