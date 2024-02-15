package main

import (
	"gpo/fonction"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	fonction.Saucisse()

	a := app.New()
	w := a.NewWindow("Hello")
	w.SetContent(widget.NewLabel("sa dit quoi l ekip"))
	w.ShowAndRun()

	/*
		// Crée l'interface utilisateur et configure la fenêtre
		w := fonction.CreateUI()

		// Affiche la fenêtre et démarre l'application
		w.ShowAndRun()*/
}
