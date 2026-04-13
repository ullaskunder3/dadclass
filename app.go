package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	ticker  *time.Ticker
	stop    chan struct{}
	running bool
	mu      sync.Mutex
}

func NewApp() *App {
	return &App{
		stop: make(chan struct{}),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Quit() {
	a.cleanup()
	runtime.Quit(a.ctx)
}

func (a *App) cleanup() {
	a.mu.Lock()
	if a.running {
		close(a.stop)
		a.running = false
	}
	a.mu.Unlock()
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

func (a *App) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.running
}

func (a *App) StartNotifications(seconds int) {
	a.mu.Lock()
	if a.running {
		close(a.stop)
		a.running = false
	}
	a.stop = make(chan struct{})
	a.running = true
	a.mu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Duration(seconds) * time.Second)
		defer ticker.Stop()

		joke := a.fetchJoke()
		if joke != "" {
			runtime.EventsEmit(a.ctx, "joke", joke)
		}

		for {
			select {
			case <-ticker.C:
				joke := a.fetchJoke()
				if joke != "" {
					runtime.EventsEmit(a.ctx, "joke", joke)
				}
			case <-a.stop:
				return
			}
		}
	}()
}

func (a *App) StopNotifications() {
	a.cleanup()
}

func (a *App) fetchJoke() string {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 8 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	var data struct {
		Joke string `json:"joke"`
	}
	json.NewDecoder(res.Body).Decode(&data)
	return data.Joke
}
