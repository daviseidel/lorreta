package wm

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>

#define MAX(a, b) ((a) > (b) ? (a) : (b))
*/
import "C"
import (
	"log"

	"github.com/daviseidel/xlib"
)

type WindowManager struct {
	Display          *xlib.Display
	Root             xlib.Window
	WindowAttributes xlib.WindowAttributes
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
	xlib.XGrabKey(
		WindowManager.Display,
		xlib.XKeysymToKeycode(WindowManager.Display, xlib.XStringToKeysym("F1")),
		int(xlib.Mod1Mask),
		WindowManager.Root,
		1,
		xlib.GrabModeAsync,
		xlib.GrabModeAsync)
	xlib.XGrabButton(
		WindowManager.Display,
		1,
		int(xlib.Mod1Mask),
		WindowManager.Root,
		xlib.Bool(1),
		xlib.ButtonPressMask|xlib.ButtonReleaseMask|xlib.PointerMotionMask,
		xlib.GrabModeAsync,
		xlib.GrabModeAsync,
		C.None,
		C.None)
	xlib.XGrabButton(
		WindowManager.Display,
		1,
		int(xlib.Mod1Mask),
		WindowManager.Root,
		xlib.Bool(1),
		xlib.ButtonPressMask|xlib.ButtonReleaseMask|xlib.PointerMotionMask,
		xlib.GrabModeAsync,
		xlib.GrabModeAsync,
		C.None,
		C.None)

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
