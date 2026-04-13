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
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/getlantern/systray"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/assets/images/32x32.ico
var trayIcon []byte

const appTitle = "dadclass"

func alreadyRunningMutex(name string) (bool, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	createMutex := kernel32.NewProc("CreateMutexW")

	ptr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return false, err
	}

	handle, _, lastErr := createMutex.Call(
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(ptr)),
	)
	if handle == 0 {
		return false, fmt.Errorf("CreateMutexW failed: %v", lastErr)
	}
	if syscall.Errno(lastErr.(syscall.Errno)) == syscall.ERROR_ALREADY_EXISTS {
		return true, nil
	}
	return false, nil
}

func main() {
	already, err := alreadyRunningMutex("Local\\DadClassNotifierMutex")
	if err != nil {
		fmt.Println("Mutex error:", err)
		os.Exit(1)
	}
	if already {
		return
	}

	app := NewApp()

	// ── Setup System Tray (Windows safe) ────────────
	go systray.Run(func() {
		systray.SetIcon(trayIcon)
		systray.SetTitle("dadclass")
		systray.SetTooltip("dadclass Notifier")

		mShow := systray.AddMenuItem("Show App", "Show the main window")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Exit the application completely")

		for {
			select {
			case <-mShow.ClickedCh:
				ShowInTaskbar(appTitle)
				if app.ctx != nil {
					runtime.WindowShow(app.ctx)
					runtime.WindowUnminimise(app.ctx)
					runtime.WindowSetAlwaysOnTop(app.ctx, true)
					runtime.WindowSetAlwaysOnTop(app.ctx, false)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				if app.ctx != nil {
					app.Quit()
				} else {
					os.Exit(0)
				}
				return
			}
		}
	}, func() {
		// onExit
	})

	// ── Window Menu (since v2 has no system tray) ────────────
	appMenu := menu.NewMenu()
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.Append(menu.Text("Show App", nil, func(_ *menu.CallbackData) {
		ShowInTaskbar(appTitle)
		runtime.WindowShow(app.ctx)
		runtime.WindowUnminimise(app.ctx)
		runtime.WindowSetAlwaysOnTop(app.ctx, true)
		runtime.WindowSetAlwaysOnTop(app.ctx, false)
	}))
	fileMenu.Append(menu.Separator())
	fileMenu.Append(menu.Text("Stop Notifications", nil, func(_ *menu.CallbackData) {
		app.StopNotifications()
	}))
	fileMenu.Append(menu.Separator())
	fileMenu.Append(menu.Text("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		app.Quit()
	}))

	err = wails.Run(&options.App{
		Title:       appTitle,
		Width:       360,
		Height:      480,
		StartHidden: false,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 18, A: 255},

		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: false,
			BackdropType:                      windows.Mica,
		},

		Menu: appMenu,

		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},

		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			selection, _ := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
				Type:          runtime.QuestionDialog,
				Title:         "Keep Running?",
				Message:       "Do you want to leave Dadclass running in the background to receive jokes?",
				Buttons:       []string{"Keep Running", "Quit"},
				DefaultButton: "Keep Running",
				CancelButton:  "Quit",
			})

			if selection == "Keep Running" {
				HideFromTaskbar(appTitle)
				runtime.WindowHide(ctx)
				return true
			}

			app.cleanup()
			return false
		},

		Bind: []interface{}{app},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
