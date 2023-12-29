package main

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>

#define MAX(a, b) ((a) > (b) ? (a) : (b))
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"log"
	"unsafe"
)

var None = C.None

type WindowManager struct {
	Display *C.Display
	// Root             xlib.Window
	WindowAttributes C.XWindowAttributes
}

func (wm *WindowManager) Create() *WindowManager {
	wm.Display = C.XOpenDisplay(nil)
	if wm.Display == nil {
		log.Fatal("Cannot open display")
	}

	// wm.Root = C.
	// if wm.Root == 0 {
	// 	log.Fatal("Cannot get root window")
	// }

	return wm
}

func (wm *WindowManager) Destroy() {
	C.XCloseDisplay(wm.Display)
}

func (WindowManager *WindowManager) Run() {
	var start C.XButtonEvent

	C.XGrabKey(
		WindowManager.Display,
		C.int(C.XKeysymToKeycode(WindowManager.Display, C.XStringToKeysym(C.CString("F1")))),
		C.Mod1Mask,
		C.XDefaultRootWindow(WindowManager.Display),
		1,
		C.GrabModeAsync,
		C.GrabModeAsync)
	C.XGrabButton(
		WindowManager.Display,
		1,
		C.Mod1Mask,
		C.XDefaultRootWindow(WindowManager.Display),
		C.Bool(1),
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
		C.GrabModeAsync,
		C.GrabModeAsync,
		C.None,
		C.None)
	C.XGrabButton(
		WindowManager.Display,
		3,
		C.Mod1Mask,
		C.XDefaultRootWindow(WindowManager.Display),
		C.Bool(1),
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
		C.GrabModeAsync,
		C.GrabModeAsync,
		C.None,
		C.None)

	start.subwindow = C.None

	for {
		loop(WindowManager)
	}

}

func loop(WindowManager *WindowManager) {
	var ev C.XEvent
	var start C.XButtonEvent

	C.XNextEvent(WindowManager.Display, &ev)

	if unionToXKeyEvent(ev).subwindow != C.None {
		switch unionToInt(ev) {
		case C.KeyPress:
			C.XRaiseWindow(WindowManager.Display, unionToXKeyEvent(ev).subwindow)

		case C.ButtonPress:
			C.XGetWindowAttributes(WindowManager.Display, unionToXButtonEvent(ev).subwindow, &WindowManager.WindowAttributes)

		case C.MotionNotify:
			xdiff := unionToXButtonEvent(ev).x_root - start.x_root
			ydiff := unionToXButtonEvent(ev).y_root - start.y_root

			var toDiffX C.int
			var toDiffY C.int

			if start.button == 1 {
				toDiffX = xdiff
				toDiffY = ydiff
			}

			var toWidth C.int
			var toHeight C.int

			if start.button == 3 {
				toWidth = xdiff
				toHeight = ydiff
			}

			C.XMoveResizeWindow(
				WindowManager.Display,
				start.subwindow,
				WindowManager.WindowAttributes.x+toDiffX,
				WindowManager.WindowAttributes.y+toDiffY,
				max(1, WindowManager.WindowAttributes.width+toWidth),
				max(1, WindowManager.WindowAttributes.height+toHeight))
		case C.ButtonRelease:
			start.subwindow = C.None
		}

	}

}

func (wm *WindowManager) checkOtherWM() {
	// Check if another window manager is running
	//
}

func unionToInt(cbytes [192]byte) (result int) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return *(*int)(unsafe.Pointer(uptr))
	}
	return 0
}

func unionToXKeyEvent(cbytes [192]byte) (result *C.XKeyEvent) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return (*C.XKeyEvent)(unsafe.Pointer(uptr))
	}
	return nil
}

func unionToXButtonEvent(cbytes [192]byte) (result *C.XButtonEvent) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return (*C.XButtonEvent)(unsafe.Pointer(uptr))
	}
	return nil
}

func max(a, b C.int) C.uint {
	if a > b {
		return C.uint(a)
	}

	return C.uint(b)
}
