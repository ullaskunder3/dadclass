//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const (
	wsExAppWindow  uint32 = 0x00040000
	wsExToolWindow uint32 = 0x00000080
)

// gwlExStylePtr returns the uintptr representation of GWL_EXSTYLE (-20)
// using two's complement, required for Windows syscall parameters.
func gwlExStylePtr() uintptr { return ^uintptr(19) }

var (
	user32dll         = syscall.NewLazyDLL("user32.dll")
	procFindWindowW   = user32dll.NewProc("FindWindowW")
	procGetWindowLong = user32dll.NewProc("GetWindowLongPtrW")
	procSetWindowLong = user32dll.NewProc("SetWindowLongPtrW")
)

func findHwnd(windowTitle string) uintptr {
	titlePtr, err := syscall.UTF16PtrFromString(windowTitle)
	if err != nil {
		return 0
	}
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	return hwnd
}

// HideFromTaskbar removes the window from the Windows taskbar.
// Adds WS_EX_TOOLWINDOW and removes WS_EX_APPWINDOW.
func HideFromTaskbar(windowTitle string) {
	hwnd := findHwnd(windowTitle)
	if hwnd == 0 {
		return
	}
	exStyle, _, _ := procGetWindowLong.Call(hwnd, gwlExStylePtr())
	newStyle := (exStyle &^ uintptr(wsExAppWindow)) | uintptr(wsExToolWindow)
	procSetWindowLong.Call(hwnd, gwlExStylePtr(), newStyle)
}

// ShowInTaskbar restores the window to the Windows taskbar.
// Removes WS_EX_TOOLWINDOW and adds WS_EX_APPWINDOW.
func ShowInTaskbar(windowTitle string) {
	hwnd := findHwnd(windowTitle)
	if hwnd == 0 {
		return
	}
	exStyle, _, _ := procGetWindowLong.Call(hwnd, gwlExStylePtr())
	newStyle := (exStyle &^ uintptr(wsExToolWindow)) | uintptr(wsExAppWindow)
	procSetWindowLong.Call(hwnd, gwlExStylePtr(), newStyle)
}
