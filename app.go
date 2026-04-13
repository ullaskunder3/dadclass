package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	ticker *time.Ticker
	stop   chan bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		stop: make(chan bool),
	}
}

// startup saves the context so we can call runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

}

func (a *App) Quit() {
	// cleanup if needed
	close(a.stop) // stop ticker, or do other cleanup
	runtime.Quit(a.ctx)
}

// StartNotifications starts sending dad jokes every X seconds
// StartNotifications starts a ticker to send dad jokes every X seconds
func (a *App) StartNotifications(seconds int) {
	go func() {
		ticker := time.NewTicker(time.Duration(seconds) * time.Second)
		a.ticker = ticker

		for {
			select {
			case <-ticker.C:
				joke := a.fetchJoke() // <---- changed here
				runtime.EventsEmit(a.ctx, "joke", joke)
			case <-a.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

// StopNotifications stops the ticker
func (a *App) StopNotifications() {
	if a.stop != nil {
		a.stop <- true
	}
}

// fetchJoke gets a random dad joke from the API
func (a *App) fetchJoke() string {
	req, _ := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
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
