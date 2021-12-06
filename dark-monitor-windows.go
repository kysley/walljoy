// https://gist.githubusercontent.com/jerblack/1d05bbcebb50ad55c312e4d7cf1bc909/raw/214bbd10d65a0d8ba3be649568f9cce8467ff4e8/dark-monitor-windows.go
package main

import (
	"log"
	"syscall"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

const (
	regKey  = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize` // in HKCU
	regName = `SystemUsesLightTheme`                                         // <- For taskbar & tray. Use AppsUseLightTheme for apps
)

func watchWindows() {
	// fmt.Println("Dark Mode on:", isDark())
	reactWindows(isDark())
	monitor(reactWindows)
}

// react to the change
func reactWindows(isDark bool) {
	if isDark {
		// fmt.Println("Dark Mode ON")
		systray.SetIcon(getIcon("smile_light.ico"))
	} else {
		// fmt.Println("Dark Mode OFF")
		systray.SetIcon(getIcon("smile_dark.ico"))
	}
}

func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, regKey, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(regName)
	if err != nil {
		log.Fatal(err)
	}
	return val == 0
}

func monitor(fn func(bool)) {
	var regNotifyChangeKeyValue *syscall.Proc
	changed := make(chan bool)

	if advapi32, err := syscall.LoadDLL("Advapi32.dll"); err == nil {
		if p, err := advapi32.FindProc("RegNotifyChangeKeyValue"); err == nil {
			regNotifyChangeKeyValue = p
		} else {
			log.Fatal("Could not find function RegNotifyChangeKeyValue in Advapi32.dll")
		}
	}
	if regNotifyChangeKeyValue != nil {
		go func() {
			k, err := registry.OpenKey(registry.CURRENT_USER, regKey, syscall.KEY_NOTIFY|registry.QUERY_VALUE)
			if err != nil {
				log.Fatal(err)
			}
			var wasDark uint64
			for {
				regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
				val, _, err := k.GetIntegerValue(regName)
				if err != nil {
					log.Fatal(err)
				}
				if val != wasDark {
					wasDark = val
					changed <- val == 0
				}
			}
		}()
	}
	for {
		val := <-changed
		fn(val)
	}

}
