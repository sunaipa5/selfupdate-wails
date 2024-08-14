package main

import (
	"context"
	"fmt"
	"selfupdate-wails/updater"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CheckUpdate() {
	updaterOptions := updater.Options{
		Author:         "sunaipa5",
		Repo:           "selfupdate-wails",
		CurrentVersion: "0.0.1",
		TagEnd:         "linux_amd64.tar.gz",
	}

	isUpdateAvailable, release := updaterOptions.CheckUpdate()

	if isUpdateAvailable {
		updaterOptions.ApplyUpdate(release)
	}
}
