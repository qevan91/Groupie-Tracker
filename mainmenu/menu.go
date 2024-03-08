package mainmenu

import (
	"fmt"
	search "gpo/Search"
	"gpo/data"
	_ "gpo/data"
	"gpo/structs"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MainMenu() {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	artists, err := structs.FetchArtists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des artistes: %v\n", err)
		return
	}

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

	for _, artist := range artists {
		artistInfo := fmt.Sprintf("Name: %s\nFirst Album: %s", artist.Name, artist.FirstAlbum)
		artistLabel := widget.NewLabel(artistInfo)

		img := data.FetchImage(artist.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))

		infoButton := widget.NewButton("", func() {
			data.ShowArtistDetails(a, artist)
		})
		infoButton.Importance = widget.LowImportance // Rend le bouton transparent
		infoButton.SetIcon(nil)                      // Aucune icône
		overlay := container.NewMax(img, infoButton) // Superpose l'image et le bouton transparent

		hbox := container.NewHBox(overlay, artistLabel)

		// Ajoutez hbox à content dans la boucle
		content.Add(hbox)
	}

	scroll := container.NewScroll(content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
