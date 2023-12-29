package main

import (
	"log"
)

func main() {
	wm := new(WindowManager).Create()
	if wm == nil {
		log.Fatal("Cannot create window manager")
	}
	defer wm.Destroy()
	wm.Run()
}
