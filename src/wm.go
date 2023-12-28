package wm

import (
	"fmt"
	"log"

	"github.com/vbsw/xlib"
)

type WindowManager struct {
	Display *xlib.Display
	Root    xlib.Window
}

func (wm *WindowManager) Create() *WindowManager {
	wm.Display = xlib.XOpenDisplay(nil)
	if wm.Display == nil {
		log.Fatal("Cannot open display")
	}

	wm.Root = xlib.XRootWindowOfScreen(xlib.XDefaultScreenOfDisplay(wm.Display))
	if wm.Root == 0 {
		log.Fatal("Cannot get root window")
	}

	return wm
}

func (wm *WindowManager) Destroy() {
	xlib.XCloseDisplay(wm.Display)
}

func (WindowManager *WindowManager) Run() {
	fmt.Println("Hello, World!")
}

func (wm *WindowManager) checkOtherWM() {
	// Check if another window manager is running
	//
}

func main() {
	wm := new(WindowManager).Create()
	if wm == nil {
		log.Fatal("Cannot create window manager")
	}
	defer wm.Destroy()
	wm.Run()
}
