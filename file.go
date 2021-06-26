package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func writeFile(s string) {
	d := []byte(s)

	if err := os.Truncate("./wall.joy", 0); err != nil {
		fmt.Printf("Failed to truncate: %v", err)
	}

	if err := ioutil.WriteFile("./wall.joy", d, 0644); err != nil {
		fmt.Print("Failed to write")
	}

}
