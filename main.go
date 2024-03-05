package main

import (
	"gpo/fonction" // Importer le package fonction

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
		searchTerm := searchEntry.Text
		fonction.PerformPostJsonRequest(w, searchTerm) // Utiliser la fonction PerformPostJsonRequest du package fonction
	})

	// Créer un conteneur avec un layout vertical pour organiser les éléments
	content := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(), // Utiliser layout.NewVBoxLayout() pour organiser verticalement
		buttonHome,
		searchEntry,
		searchButton,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
