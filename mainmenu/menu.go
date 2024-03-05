package mainmenu

import (
	search "gpo/Search"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MainMenu() {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
		searchTerm := searchEntry.Text
		search.PerformPostJsonRequest(w, searchTerm)
	})

	content := container.NewVBox(
		buttonHome,
		searchEntry,
		searchButton,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
