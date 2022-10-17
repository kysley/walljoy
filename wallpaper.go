package main

import (
	"fmt"
	"log"
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
func getDeviceIdentity() []string {
	deviceName, _ := os.Hostname()
	deviceId, err := machineid.ProtectedID("walljoy")

	if err != nil {
		fmt.Print("Error getting device identity")
	}

	return []string{deviceId, deviceName}
}

func openSession() {
	var code Kv
	identity := getDeviceIdentity()

	store.Get("code", &code)
	fmt.Print(code)
	openBrowser(fmt.Sprintf("http://127.0.0.1:5173/register?code=%s&deviceId=%s&name=%s", code.Value, identity[0], identity[1]))
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

	code := RegisterDevice(getDeviceIdentity()[0])
	fmt.Print(code)
	store.Set("code", Kv{code})

	setWallpaper(GetCollectionLatest("1"))

	// if err != nil {
	// 	collectionId := Kv{"1"}
	// 	store.Set("collectionId", collectionId)
	// 	fmt.Print(collectionId.Value)
	// }

	// setWallpaper()

	// c.AddFunc("@midnight", setWallpaper)
	c.Start()
}

func onExit() {
	systray.Quit()
	c.Stop()
}

func setWallpaper(url string) {
	// collectionId := Kv{}
	// if err := store.Get("collectionId", &collectionId); err != nil {
	// 	print(err)
	// }

	wallpaper.SetFromURL(url)

	fmt.Println("setwallpaper finished")
}

func onCollectionChange(cId int) {
	stringCId := strconv.Itoa(cId)
	store.Set("collectionId", Kv{stringCId})
	// setWallpaper()
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
