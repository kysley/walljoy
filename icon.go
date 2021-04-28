package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
)

func getIcon(s string) []byte {
	var b []byte
	if runtime.GOOS == "windows" {
		ba, err := ioutil.ReadFile(s)
		if err != nil {
			fmt.Print(err)
		}
		b = ba
	} else if runtime.GOOS == "darwin" {
		ba, err := ioutil.ReadFile(s)
		if err != nil {
			fmt.Print(err)
		}
		b = ba
	}
	return b
}
