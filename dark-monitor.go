// https://gist.github.com/jerblack/869a303d1a604171bf8f00bbbefa59c2#file-2-dark-monitor-go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
)

const plistPath = `/Library/Preferences/.GlobalPreferences.plist`

var plist = filepath.Join(os.Getenv("HOME"), plistPath)
var wasDark bool

func watchOSX() {
	// get initial state
	wasDark = checkDarkMode()
	reactOSX(wasDark)

	// Start watcher and give it a function to call when the state changes
	startWatcher(reactOSX)
}

// react to the change
func reactOSX(isDark bool) {
	if isDark {
		fmt.Println("Dark Mode ON")
		// @todo osx icons
		systray.SetIcon(getIcon("smile_light.ico"))
	} else {
		systray.SetIcon(getIcon("smile_dark.ico"))
		fmt.Println("Dark Mode OFF")
	}
}

func checkDarkMode() bool {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false
		}
	}
	return true
}

func startWatcher(fn func(bool)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					isDark := checkDarkMode()
					if isDark && !wasDark {
						fn(isDark)
						wasDark = isDark
					}
					if !isDark && wasDark {
						fn(isDark)
						wasDark = isDark
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(plist)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
