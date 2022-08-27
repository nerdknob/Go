package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(newGameTheme())
	show(a)
	a.Run()
}

func show(app fyne.App) {
	w := app.NewWindow("Blackjack")
	w.SetPadded(false)
	w.CenterOnScreen()
	game := NewGame()
	w.SetContent(NewTable(game))
	w.Resize(fyne.NewSize(minWidth, minHeight))
	w.Show()
}
