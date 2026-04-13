package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func alreadyRunningMutex(name string) (bool, error) {
	// Windows API: CreateMutexW
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	createMutex := kernel32.NewProc("CreateMutexW")

	// UTF16 pointer for the name
	ptr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return false, err
	}

	// Call CreateMutexW(nil, FALSE, name)
	handle, _, lastErr := createMutex.Call(
		uintptr(0),                   // lpMutexAttributes = NULL
		uintptr(0),                   // bInitialOwner   = FALSE
		uintptr(unsafe.Pointer(ptr)), // lpName
	)
	if handle == 0 {
		// Some real error occurred
		return false, fmt.Errorf("CreateMutexW failed: %v", lastErr)
	}

	// If ERROR_ALREADY_EXISTS, return true
	if syscall.Errno(lastErr.(syscall.Errno)) == syscall.ERROR_ALREADY_EXISTS {
		return true, nil
	}

	// Otherwise, we hold the mutex and will release it on exit.
	return false, nil
}

func main() {
	// 1) Check single instance
	already, err := alreadyRunningMutex("Local\\DadJokeNotifierMutex")
	if err != nil {
		fmt.Println("Mutex error:", err)
		// If you can’t create the mutex, we bail out to avoid chaos.
		os.Exit(1)
	}
	if already {
		// Another instance is already running. Just exit silently.
		return
	}

	// 2) If we reach here, no other instance is running. Proceed:
	app := NewApp()

	// 3) Build a system tray menu
	trayMenu := menu.NewMenu()
	trayMenu.Append(menu.Text("Show App", nil, func(_ *menu.CallbackData) {
		runtime.Show(app.ctx)
	}))
	trayMenu.Append(menu.Text("Stop Notifications", nil, func(_ *menu.CallbackData) {
		app.StopNotifications()
	}))
	trayMenu.Append(menu.Text("Quit", nil, func(_ *menu.CallbackData) {
		app.Quit()
	}))

	// 4) Run Wails, but override OnBeforeClose to hide instead of quit.
	err = wails.Run(&options.App{
		Title:  "Dad Joke Notifier",
		Width:  320,
		Height: 500,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 20, B: 20, A: 150},

		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
			BackdropType:                      windows.Mica,
		},
		Menu: trayMenu,

		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			// Instead of quitting, just hide the window.
			runtime.Hide(ctx)
			return true // “prevent” the window from really closing
		},
		Bind: []interface{}{app},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
