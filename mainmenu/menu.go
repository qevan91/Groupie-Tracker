package mainmenu

import (
	"fmt"
	search "gpo/Search"
	"gpo/data"
	"gpo/structs"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

	title := canvas.NewText("Groupie Tracker", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
		searchTerm := searchEntry.Text
		search.PerformPostJsonRequest(w, searchTerm)
	})

	artistContainer := container.NewHBox()

	overlay := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil))

	searchEntry.OnChanged = func(text string) {
		artistContainer.RemoveAll()
		overlay.RemoveAll()
		art, err := structs.GetArtistByName(text)
		if err != nil {
			fmt.Println("Problème")
			return
		}
		img := data.FetchImage(art.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))
		if text == "" {
			artistContainer.RemoveAll()
			overlay.RemoveAll()
			return
		}
		overlay.Add(img)
		artistContainer.Add(widget.NewLabel(fmt.Sprintf("Name: %s", art.Name)))
	}

	grid := container.NewGridWithRows(0)
	count := 0
	listeArtist := container.NewVBox()

	for _, artist := range artists {
		img := data.FetchImage(artist.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))
		infoButton := widget.NewButton("", func() {
			data.ShowArtistDetails(a, artist)
		})
		infoButton.Importance = widget.LowImportance
		infoButton.SetIcon(nil)
		artistInfo := fmt.Sprintf(artist.Name)
		artistLabel := widget.NewLabel(artistInfo)
		overlay := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
			img,
			artistLabel,
			infoButton,
		)
		grid.Add(overlay)
		count++
		if count == 4 {
			count = 0
			listeArtist.Add(grid)
			grid = container.NewGridWithRows(0)
		}
	}
	if count > 0 {
		listeArtist.Add(grid)
	}

	content := container.NewVBox(
		titleContainer,
		buttonHome,
		searchEntry,
		searchButton,
		overlay,
		artistContainer,
		listeArtist,
	)

	scroll := container.NewScroll(content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
