package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var App fyne.App

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	App = app.New()
	InitUI()
	App.Run()
}
