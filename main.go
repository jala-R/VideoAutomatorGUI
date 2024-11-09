package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	mainscreen "jalar.me/VideoCreatorGUI/MainScreen"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(mainscreen.GetGUI(w))

	w.Resize(fyne.NewSize(1080, 1920))
	w.ShowAndRun()
}
